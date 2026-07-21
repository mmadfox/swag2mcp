package index

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"
	"github.com/blugelabs/bluge/analysis/tokenizer"
	"github.com/blugelabs/bluge/search"
	querystring "github.com/blugelabs/query_string"
	"github.com/mmadfox/swag2mcp/internal/model"
)

const initialSpecsCapacity = 8

// newAnalyzer creates a default text analyzer with unicode tokenizer and lower-case filter.
func newAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Tokenizer: tokenizer.NewUnicodeTokenizer(),
		TokenFilters: []analysis.TokenFilter{
			token.NewLowerCaseFilter(),
		},
	}
}

// EndpointCursor represents a position in the index.
type EndpointCursor struct {
	Spec       *model.Spec
	Tag        *model.Tag
	Endpoint   *model.Endpoint
	Collection *model.Collection
}

// TagCursor represents a position in the index.
type TagCursor struct {
	Spec       *model.Spec
	Tag        *model.Tag
	Collection *model.Collection
}

// CollectionCursor represents a position in the index.
type CollectionCursor struct {
	Spec       *model.Spec
	Collection *model.Collection
}

// Index represents the in-memory index structure.
type Index struct {
	mu                    sync.RWMutex
	specs                 map[string]*model.Spec         // specID -> Spec
	allSpecs              []*model.Spec                  // all specs
	collectionsBySpec     map[string][]*model.Collection // specID -> []Collection
	collectionByID        map[string]*model.Collection   // collectionID -> Collection
	tagsByCollection      map[string][]*model.Tag        // collectionID -> []Tag
	tagByID               map[string]*model.Tag          // tagID -> Tag
	tagBySpec             map[string][]*model.Tag        // specID -> []Tag
	endpointsByTag        map[string][]*model.Endpoint   // tagID -> []Endpoint
	endpointsBySpec       map[string][]*model.Endpoint   // specID -> []Endpoint
	endpointsByCollection map[string][]*model.Endpoint   // collectionID -> []Endpoint
	endpointByID          map[string]*model.Endpoint     // endpointID -> Endpoint
	blugeWriter           *bluge.Writer
	blugeReader           atomic.Pointer[bluge.Reader]
	analyzer              *analysis.Analyzer
	noFullText            bool
}

// NewOption is a functional option for configuring an Index.
type NewOption func(*Index)

// WithNoFullText disables full-text search indexing.
// Use this for CLI commands that only need in-memory lookups.
func WithNoFullText() NewOption {
	return func(idx *Index) {
		idx.noFullText = true
	}
}

// New creates an empty in-memory index with type-based structures.
func New(opts ...NewOption) (*Index, error) {
	writer, err := bluge.OpenWriter(bluge.InMemoryOnlyConfig())
	if err != nil {
		return nil, fmt.Errorf("bluge open: %w", err)
	}
	idx := &Index{
		blugeWriter:           writer,
		specs:                 make(map[string]*model.Spec),
		collectionsBySpec:     make(map[string][]*model.Collection),
		collectionByID:        make(map[string]*model.Collection),
		tagsByCollection:      make(map[string][]*model.Tag),
		tagByID:               make(map[string]*model.Tag),
		tagBySpec:             make(map[string][]*model.Tag),
		endpointsByTag:        make(map[string][]*model.Endpoint),
		endpointsBySpec:       make(map[string][]*model.Endpoint),
		endpointsByCollection: make(map[string][]*model.Endpoint),
		endpointByID:          make(map[string]*model.Endpoint),
		allSpecs:              make([]*model.Spec, 0, initialSpecsCapacity),
		analyzer:              newAnalyzer(),
	}
	for _, opt := range opts {
		opt(idx)
	}
	return idx, nil
}

// EnsureIndex indexes all provided data: (spec, collections, tags, endpoints).
func (idx *Index) EnsureIndex(
	spec *model.Spec,
	colls []*model.Collection,
	tags []*model.Tag,
	endpoints []*model.Endpoint,
) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if err := idx.indexSpec(spec); err != nil {
		return fmt.Errorf("indexing spec: %w", err)
	}

	if err := idx.indexCollections(colls); err != nil {
		return fmt.Errorf("indexing collections: %w", err)
	}

	if err := idx.indexTags(tags); err != nil {
		return fmt.Errorf("indexing tags: %w", err)
	}

	if err := idx.indexEndpoints(endpoints); err != nil {
		return fmt.Errorf("indexing endpoints: %w", err)
	}

	if !idx.noFullText {
		return idx.index(endpoints)
	}
	return nil
}

// indexSpec indexes a spec.
func (idx *Index) indexSpec(spec *model.Spec) error {
	if _, exists := idx.specs[spec.ID]; exists {
		return fmt.Errorf("spec-%q(%s) already exists", spec.ID, spec.Domain)
	}
	idx.specs[spec.ID] = spec
	idx.allSpecs = append(idx.allSpecs, spec)
	return nil
}

// indexCollections indexes collections.
func (idx *Index) indexCollections(colls []*model.Collection) error {
	for _, coll := range colls {
		if _, exists := idx.specs[coll.SpecID]; !exists {
			return fmt.Errorf("spec-%q(%s) not found", coll.SpecID, coll.LLMTitle)
		}

		if _, exists := idx.collectionByID[coll.ID]; exists {
			return fmt.Errorf("collection-%q(%s) already exists", coll.ID, coll.LLMTitle)
		}

		idx.collectionByID[coll.ID] = coll
		idx.collectionsBySpec[coll.SpecID] = append(idx.collectionsBySpec[coll.SpecID], coll)
	}
	return nil
}

// indexTags indexes tags.
func (idx *Index) indexTags(tags []*model.Tag) error {
	for _, tag := range tags {
		if _, exists := idx.specs[tag.SpecID]; !exists {
			return fmt.Errorf("spec-%q not found", tag.SpecID)
		}

		if _, exists := idx.collectionByID[tag.CollectionID]; !exists {
			return fmt.Errorf("collection-%q not found", tag.ID)
		}

		if _, exists := idx.tagByID[tag.ID]; exists {
			return fmt.Errorf("tag-%q(%s) already exists", tag.ID, tag.Name)
		}

		idx.tagByID[tag.ID] = tag
		idx.tagsByCollection[tag.CollectionID] = append(idx.tagsByCollection[tag.CollectionID], tag)
		idx.tagBySpec[tag.SpecID] = append(idx.tagBySpec[tag.SpecID], tag)
	}
	return nil
}

// indexEndpoints indexes endpoints.
func (idx *Index) indexEndpoints(endpoints []*model.Endpoint) error {
	for _, ep := range endpoints {
		if _, exists := idx.specs[ep.SpecID]; !exists {
			return fmt.Errorf("spec-%q not found", ep.SpecID)
		}

		if _, exists := idx.collectionByID[ep.CollectionID]; !exists {
			return fmt.Errorf("collection-%q not found", ep.ID)
		}

		if _, exists := idx.tagByID[ep.TagID]; !exists {
			return fmt.Errorf("tag-%q not found", ep.ID)
		}

		if _, exists := idx.endpointByID[ep.ID]; exists {
			return fmt.Errorf("endpoint-%q already exists", ep.ID)
		}

		idx.endpointByID[ep.ID] = ep
		idx.endpointsByTag[ep.TagID] = append(idx.endpointsByTag[ep.TagID], ep)
		idx.endpointsBySpec[ep.SpecID] = append(idx.endpointsBySpec[ep.SpecID], ep)
		idx.endpointsByCollection[ep.CollectionID] = append(idx.endpointsByCollection[ep.CollectionID], ep)
	}
	return nil
}

// index indexes the full-text search endpoints.
func (idx *Index) index(endpoints []*model.Endpoint) error {
	batch := bluge.NewBatch()
	for _, ep := range endpoints {
		summary := strings.ToLower(ep.SummaryOrFallback())
		doc := bluge.NewDocument(ep.ID).
			AddField(bluge.NewKeywordField("method", strings.ToLower(ep.Name)).StoreValue()).
			AddField(bluge.NewKeywordField("tag", strings.ToLower(ep.Tag)).StoreValue()).
			AddField(bluge.NewKeywordField("path", strings.ToLower(ep.Path)).StoreValue()).
			AddField(bluge.NewTextField("summary", strings.ToLower(summary)).WithAnalyzer(idx.analyzer).StoreValue().SearchTermPositions()).
			AddField(bluge.NewTextField("_all", fmt.Sprintf("%s %s %s %s", strings.ToLower(ep.Name), strings.ToLower(ep.Path), strings.ToLower(ep.Tag), strings.ToLower(summary))).WithAnalyzer(idx.analyzer).SearchTermPositions())

		batch.Update(bluge.Identifier(ep.ID), doc)
	}

	if err := idx.blugeWriter.Batch(batch); err != nil {
		return fmt.Errorf("indexing endpoints: %w", err)
	}

	return nil
}

// AllSpecs returns all indexed specs.
func (idx *Index) AllSpecs() []*model.Spec {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	return append([]*model.Spec(nil), idx.allSpecs...)
}

// SpecByID returns a spec by its ID.
func (idx *Index) SpecByID(specID string) (*model.Spec, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	spec, ok := idx.specs[specID]
	if !ok {
		return nil, fmt.Errorf("spec by id %q not found", specID)
	}
	return spec, nil
}

// CollectionByID returns a collection by its ID.
func (idx *Index) CollectionByID(collectionID string) (*model.Collection, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	col, ok := idx.collectionByID[collectionID]
	if !ok {
		return nil, fmt.Errorf("collection by id %q not found", collectionID)
	}
	return col, nil
}

// CollectionsBySpec returns all collections for a given spec ID.
func (idx *Index) CollectionsBySpec(specID string) ([]*model.Collection, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	cols, ok := idx.collectionsBySpec[specID]
	if !ok {
		return nil, fmt.Errorf("collection by spec %q not found", specID)
	}
	return append([]*model.Collection(nil), cols...), nil
}

// TagByID returns a tag by its ID.
func (idx *Index) TagByID(tagID string) (*model.Tag, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	tag, ok := idx.tagByID[tagID]
	if !ok {
		return nil, fmt.Errorf("tag by id %q not found", tagID)
	}
	return tag, nil
}

// TagsByCollection returns all tags for a given collection ID.
func (idx *Index) TagsByCollection(collectionID string) ([]*model.Tag, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	tags, ok := idx.tagsByCollection[collectionID]
	if !ok {
		return nil, fmt.Errorf("tags by collection %q not found", collectionID)
	}
	return append([]*model.Tag(nil), tags...), nil
}

// TagsBySpec returns all tags for a given spec ID.
func (idx *Index) TagsBySpec(specID string) ([]*model.Tag, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	tags, ok := idx.tagBySpec[specID]
	if !ok {
		return nil, fmt.Errorf("tags by spec %q not found", specID)
	}
	return append([]*model.Tag(nil), tags...), nil
}

// EndpointsByTag returns all endpoints for a given tag ID.
func (idx *Index) EndpointsByTag(tagID string) ([]*model.Endpoint, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	endpoints, ok := idx.endpointsByTag[tagID]
	if !ok {
		return nil, fmt.Errorf("endpoints by tag %q not found", tagID)
	}
	return append([]*model.Endpoint(nil), endpoints...), nil
}

// EndpointsBySpec returns all endpoints for a given spec ID.
func (idx *Index) EndpointsBySpec(specID string) ([]*model.Endpoint, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	endpoints, ok := idx.endpointsBySpec[specID]
	if !ok {
		return nil, fmt.Errorf("spec %q not found", specID)
	}
	return append([]*model.Endpoint(nil), endpoints...), nil
}

// EndpointByCollection returns all endpoints for a given collection ID.
func (idx *Index) EndpointByCollection(collectionID string) ([]*model.Endpoint, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	endpoints, ok := idx.endpointsByCollection[collectionID]
	if !ok {
		return nil, fmt.Errorf("endpoints by collection %q not found", collectionID)
	}
	return append([]*model.Endpoint(nil), endpoints...), nil
}

// EndpointByID returns an endpoint by its ID.
func (idx *Index) EndpointByID(id string) (*model.Endpoint, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	if ep, ok := idx.endpointByID[id]; ok {
		return ep, nil
	}
	return nil, fmt.Errorf("endpoint by %q not found", id)
}

// IterateByEndpoints returns an iterator over all endpoints.
func (idx *Index) IterateByEndpoints() iter.Seq[*EndpointCursor] {
	return func(yield func(*EndpointCursor) bool) {
		idx.mu.RLock()
		defer idx.mu.RUnlock()

		for _, ep := range idx.endpointByID {
			spec := idx.specs[ep.SpecID]
			coll := idx.collectionByID[ep.CollectionID]
			tag := idx.tagByID[ep.TagID]
			cursor := &EndpointCursor{
				Endpoint:   ep,
				Spec:       spec,
				Collection: coll,
				Tag:        tag,
			}
			if !yield(cursor) {
				return
			}
		}
	}
}

// IterateByTags returns an iterator over all tags.
func (idx *Index) IterateByTags() iter.Seq[*TagCursor] {
	return func(yield func(*TagCursor) bool) {
		idx.mu.RLock()
		defer idx.mu.RUnlock()

		for _, tag := range idx.tagByID {
			spec := idx.specs[tag.SpecID]
			coll := idx.collectionByID[tag.CollectionID]
			cursor := &TagCursor{
				Spec:       spec,
				Collection: coll,
				Tag:        tag,
			}
			if !yield(cursor) {
				return
			}
		}
	}
}

// IterateByCollections returns an iterator over all collections.
func (idx *Index) IterateByCollections() iter.Seq[*CollectionCursor] {
	return func(yield func(*CollectionCursor) bool) {
		idx.mu.RLock()
		defer idx.mu.RUnlock()

		for _, col := range idx.collectionByID {
			spec := idx.specs[col.SpecID]
			coll := idx.collectionByID[col.ID]
			cursor := &CollectionCursor{
				Spec:       spec,
				Collection: coll,
			}
			if !yield(cursor) {
				return
			}
		}
	}
}

// Search returns endpoints matching the query.
func (idx *Index) Search(ctx context.Context, q string, limit int) ([]*model.Endpoint, error) {
	if len(q) == 0 {
		return nil, errors.New("query string is required")
	}

	if limit <= 0 || limit > 50 {
		limit = 50
	}

	query := idx.buildQuery(q)
	r, err := idx.reader()
	if err != nil {
		return nil, err
	}

	return idx.collectResults(ctx, r, query, limit)
}

// buildQuery parses a query string into a bluge query, falling back to a match-all or match query.
func (idx *Index) buildQuery(q string) bluge.Query {
	if q == "*" {
		return bluge.NewMatchAllQuery()
	}

	qsOpts := querystring.DefaultOptions().
		WithDefaultAnalyzer(idx.analyzer).
		WithAnalyzerForField("summary", idx.analyzer).
		WithAnalyzerForField("_all", idx.analyzer)

	if parsedQuery, err := querystring.ParseQueryString(q, qsOpts); err == nil {
		return parsedQuery
	}

	return bluge.NewMatchQuery(q).SetField("_all").SetAnalyzer(idx.analyzer)
}

// reader returns a bluge reader, lazily initializing it via CAS on first call.
func (idx *Index) reader() (*bluge.Reader, error) {
	r := idx.blugeReader.Load()
	if r != nil {
		return r, nil
	}

	newReader, err := idx.blugeWriter.Reader()
	if err != nil {
		return nil, fmt.Errorf("bluge reader: %w", err)
	}

	if !idx.blugeReader.CompareAndSwap(nil, newReader) {
		if closeErr := newReader.Close(); closeErr != nil {
			slog.Default().Debug("closing bluge reader after CAS", "error", closeErr)
		}
	}

	return idx.blugeReader.Load(), nil
}

// collectResults iterates search results and maps document IDs to indexed endpoints.
func (idx *Index) collectResults(
	ctx context.Context,
	r *bluge.Reader,
	query bluge.Query,
	limit int,
) ([]*model.Endpoint, error) {
	req := bluge.NewTopNSearch(limit, query)
	itr, err := r.Search(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("bluge search: %w", err)
	}

	idx.mu.RLock()
	defer idx.mu.RUnlock()

	out := make([]*model.Endpoint, 0, limit)
	match, iterErr := itr.Next()
	for iterErr == nil && match != nil {
		docID := extractDocID(match)
		if docID != "" {
			if ep, ok := idx.endpointByID[docID]; ok {
				out = append(out, ep)
			}
		}
		match, iterErr = itr.Next()
	}
	if iterErr != nil {
		return out, fmt.Errorf("bluge iterate: %w", iterErr)
	}

	return out, nil
}

// extractDocID reads the _id stored field from a bluge document match.
func extractDocID(match *search.DocumentMatch) string {
	var docID string
	if err := match.VisitStoredFields(func(field string, value []byte) bool {
		if field == "_id" {
			docID = string(value)
			return false
		}
		return true
	}); err != nil {
		slog.Default().Debug("visiting stored fields", "error", err)
	}
	return docID
}

// RefreshSearchReader forces the bluge reader to be re-created from the writer,
// making newly indexed documents visible to Search.
func (idx *Index) RefreshSearchReader() {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if r := idx.blugeReader.Load(); r != nil {
		if err := r.Close(); err != nil {
			slog.Default().Debug("closing bluge reader on refresh", "error", err)
		}
	}

	newReader, err := idx.blugeWriter.Reader()
	if err != nil {
		return
	}
	idx.blugeReader.Store(newReader)
}

// Close releases all index resources.
func (idx *Index) Close() error {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	if r := idx.blugeReader.Load(); r != nil {
		if err := r.Close(); err != nil {
			slog.Default().Debug("closing bluge reader on close", "error", err)
		}
	}
	return idx.blugeWriter.Close()
}

// Size returns the total number of indexed endpoints.
func (idx *Index) Size() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return len(idx.endpointByID)
}

// RemoveSpec removes a spec from the index by its ID.
func (idx *Index) RemoveSpec(specID string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	delete(idx.specs, specID)
}

// RemoveCollection removes a collection from the index by its ID.
func (idx *Index) RemoveCollection(collectionID string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	delete(idx.collectionByID, collectionID)
}

// RemoveTag removes a tag from the index by its ID.
func (idx *Index) RemoveTag(tagID string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	delete(idx.tagByID, tagID)
}

// RemoveAllTags removes all tags from the index.
func (idx *Index) RemoveAllTags() {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.tagByID = make(map[string]*model.Tag)
	idx.tagsByCollection = make(map[string][]*model.Tag)
	idx.tagBySpec = make(map[string][]*model.Tag)
}

// RemoveCollectionsBySpec removes all collections for a given spec ID.
func (idx *Index) RemoveCollectionsBySpec(specID string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	delete(idx.collectionsBySpec, specID)
}

// AddCollection adds a collection to the index.
func (idx *Index) AddCollection(coll *model.Collection) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.collectionByID[coll.ID] = coll
	idx.collectionsBySpec[coll.SpecID] = append(idx.collectionsBySpec[coll.SpecID], coll)
}
