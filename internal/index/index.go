package index

import (
	"context"
	"fmt"
	"iter"
	"sync"
	"sync/atomic"

	"github.com/blugelabs/bluge"
	querystring "github.com/blugelabs/query_string"
	"github.com/mmadfox/swag2mcp/internal/types"
)

// Cursor represents a position in the index.
type EndpointCursor struct {
	Spec       *types.Spec
	Tag        *types.Tag
	Endpoint   *types.Endpoint
	Collection *types.Collection
}

// TagCursor represents a position in the index.
type TagCursor struct {
	Spec       *types.Spec
	Tag        *types.Tag
	Collection *types.Collection
}

// CollectionCursor represents a position in the index.
type CollectionCursor struct {
	Spec       *types.Spec
	Collection *types.Collection
}

// Index represents the in-memory index structure.
type Index struct {
	mu                    sync.RWMutex
	specs                 map[string]*types.Spec         // specID -> Spec
	allSpecs              []*types.Spec                  // all specs
	collectionsBySpec     map[string][]*types.Collection // specID -> []Collection
	collectionByID        map[string]*types.Collection   // collectionID -> Collection
	tagsByCollection      map[string][]*types.Tag        // collectionID -> []Tag
	tagByID               map[string]*types.Tag          // tagID -> Tag
	tagBySpec             map[string][]*types.Tag        // specID -> []Tag
	endpointsByTag        map[string][]*types.Endpoint   // tagID -> []Endpoint
	endpointsBySpec       map[string][]*types.Endpoint   // specID -> []Endpoint
	endpointsByCollection map[string][]*types.Endpoint   // collectionID -> []Endpoint
	endpointByID          map[string]*types.Endpoint     // endpointID -> Endpoint
	blugeWriter           *bluge.Writer
	blugeReader           atomic.Pointer[bluge.Reader]
}

// New creates an empty in-memory index with type-based structures.
func New() (*Index, error) {
	writer, err := bluge.OpenWriter(bluge.InMemoryOnlyConfig())
	if err != nil {
		return nil, fmt.Errorf("bluge open: %w", err)
	}
	return &Index{
		blugeWriter:           writer,
		specs:                 make(map[string]*types.Spec),
		collectionsBySpec:     make(map[string][]*types.Collection),
		collectionByID:        make(map[string]*types.Collection),
		tagsByCollection:      make(map[string][]*types.Tag),
		tagByID:               make(map[string]*types.Tag),
		tagBySpec:             make(map[string][]*types.Tag),
		endpointsByTag:        make(map[string][]*types.Endpoint),
		endpointsBySpec:       make(map[string][]*types.Endpoint),
		endpointsByCollection: make(map[string][]*types.Endpoint),
		endpointByID:          make(map[string]*types.Endpoint),
		allSpecs:              make([]*types.Spec, 0, 8),
	}, nil
}

// EnsureIndex indexes all provided data: (spec, collections, tags, endpoints).
func (idx *Index) EnsureIndex(
	spec *types.Spec,
	colls []*types.Collection,
	tags []*types.Tag,
	endpoints []*types.Endpoint,
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

	return idx.index(endpoints)
}

// indexSpec indexes a spec.
func (idx *Index) indexSpec(spec *types.Spec) error {
	if _, exists := idx.specs[spec.ID]; exists {
		return fmt.Errorf("spec-%q(%s) already exists", spec.ID, spec.Domain)
	}
	idx.specs[spec.ID] = spec
	idx.allSpecs = append(idx.allSpecs, spec)
	return nil
}

// indexCollections indexes collections.
func (idx *Index) indexCollections(colls []*types.Collection) error {
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
func (idx *Index) indexTags(tags []*types.Tag) error {
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
func (idx *Index) indexEndpoints(endpoints []*types.Endpoint) error {
	for _, ep := range endpoints {
		if _, exists := idx.specs[ep.SpecID]; !exists {
			return fmt.Errorf("spec-%q not found", ep.SpecID)
		}

		if _, exists := idx.collectionByID[ep.CollectionID]; !exists {
			return fmt.Errorf("collection-%q not foudn", ep.ID)
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
func (idx *Index) index(endpoints []*types.Endpoint) error {
	batch := bluge.NewBatch()
	for _, ep := range endpoints {
		summary := ep.SummaryOrFallback()
		doc := bluge.NewDocument(ep.ID).
			AddField(bluge.NewKeywordField("method", ep.Name).StoreValue()).
			AddField(bluge.NewKeywordField("tag", ep.Tag).StoreValue()).
			AddField(bluge.NewTextField("path", ep.Path).StoreValue()).
			AddField(bluge.NewTextField("summary", summary).StoreValue()).
			AddField(bluge.NewTextField("_all", fmt.Sprintf("%s %s %s %s", ep.Name, ep.Path, ep.Tag, ep.SummaryOrFallback())))

		batch.Update(bluge.Identifier(ep.ID), doc)
	}

	if err := idx.blugeWriter.Batch(batch); err != nil {
		return fmt.Errorf("indexing endpoints: %w", err)
	}

	return nil
}

// AllSpecs returns all indexed specs.
func (idx *Index) AllSpecs() []*types.Spec {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	return idx.allSpecs
}

// GetSpec returns a spec by its ID.
func (idx *Index) SpecByID(specID string) (*types.Spec, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	spec, ok := idx.specs[specID]
	if !ok {
		return nil, fmt.Errorf("spec by id %q not found", specID)
	}
	return spec, nil
}

// CollectionByID returns a collection by its ID.
func (idx *Index) CollectionByID(collectionID string) (*types.Collection, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	col, ok := idx.collectionByID[collectionID]
	if !ok {
		return nil, fmt.Errorf("collection by id %q not found", collectionID)
	}
	return col, nil
}

// CollectionsBySpec returns all collections for a given spec ID.
func (idx *Index) CollectionsBySpec(specID string) ([]*types.Collection, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	cols, ok := idx.collectionsBySpec[specID]
	if !ok {
		return nil, fmt.Errorf("collection by spec %q not found", specID)
	}
	return cols, nil
}

// TagByID returns a tag by its ID.
func (idx *Index) TagByID(tagID string) (*types.Tag, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	tag, ok := idx.tagByID[tagID]
	if !ok {
		return nil, fmt.Errorf("tag by id %q not found", tagID)
	}
	return tag, nil
}

// TagsByCollection returns all tags for a given collection ID.
func (idx *Index) TagsByCollection(collectionID string) ([]*types.Tag, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	tags, ok := idx.tagsByCollection[collectionID]
	if !ok {
		return nil, fmt.Errorf("tags by collection %q not found", collectionID)
	}
	return tags, nil
}

// TagsBySpec returns all tags for a given spec ID.
func (idx *Index) TagsBySpec(specID string) ([]*types.Tag, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	tags, ok := idx.tagBySpec[specID]
	if !ok {
		return nil, fmt.Errorf("tags by spec %q not found", specID)
	}
	return tags, nil
}

// EndpointsByTag returns all endpoints for a given tag ID.
func (idx *Index) EndpointsByTag(tagID string) ([]*types.Endpoint, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	endpoints, ok := idx.endpointsByTag[tagID]
	if !ok {
		return nil, fmt.Errorf("endpoints by tag %q not found", tagID)
	}
	return endpoints, nil
}

// EndpointsBySpec returns all endpoints for a given spec ID.
func (idx *Index) EndpointsBySpec(specID string) ([]*types.Endpoint, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	endpoints, ok := idx.endpointsBySpec[specID]
	if !ok {
		return nil, fmt.Errorf("spec %q not found", specID)
	}
	return endpoints, nil
}

// EndpointsByCollection returns all endpoints for a given collection ID.
func (idx *Index) EndpointByCollection(collectionID string) ([]*types.Endpoint, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	endpoints, ok := idx.endpointsByCollection[collectionID]
	if !ok {
		return nil, fmt.Errorf("endpoints by collection %q not found", collectionID)
	}
	return endpoints, nil
}

// EndpointByID returns an endpoint by its ID.
func (idx *Index) EndpointByID(id string) (*types.Endpoint, error) {
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
func (idx *Index) Search(ctx context.Context, q string, limit int) ([]*types.Endpoint, error) {
	if len(q) == 0 {
		return nil, fmt.Errorf("query string is required")
	}

	if limit <= 0 || limit > 20 {
		limit = 20
	}

	var query bluge.Query
	if q == "*" {
		query = bluge.NewMatchAllQuery()
	} else {
		if parsedQuery, err := querystring.ParseQueryString(q, querystring.DefaultOptions()); err == nil {
			query = parsedQuery
		}
		if query == nil {
			query = bluge.NewMatchQuery(q).SetField("_all")
		}
	}

	reader := idx.blugeReader.Load()
	if reader == nil {
		r, err := idx.blugeWriter.Reader()
		if err != nil {
			return nil, fmt.Errorf("bluge reader: %w", err)
		}
		if !idx.blugeReader.CompareAndSwap(nil, r) {
			r.Close()
		}
		reader = idx.blugeReader.Load()
	}

	req := bluge.NewTopNSearch(limit, query)
	itr, err := reader.Search(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("bluge search: %w", err)
	}

	idx.mu.RLock()
	defer idx.mu.RUnlock()

	out := make([]*types.Endpoint, 0, limit)
	match, err := itr.Next()
	for err == nil && match != nil {
		var docID string
		_ = match.VisitStoredFields(func(field string, value []byte) bool {
			if field == "_id" {
				docID = string(value)
				return false
			}
			return true
		})
		if docID != "" {
			if ep, ok := idx.endpointByID[docID]; ok {
				out = append(out, ep)
			}
		}
		match, err = itr.Next()
	}
	if err != nil {
		return out, fmt.Errorf("bluge iterate: %w", err)
	}

	return out, nil
}

// Close releases all index resources.
func (idx *Index) Close() error {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	if r := idx.blugeReader.Load(); r != nil {
		_ = r.Close()
	}
	return idx.blugeWriter.Close()
}

// Size returns the total number of indexed endpoints.
func (idx *Index) Size() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return len(idx.endpointByID)
}
