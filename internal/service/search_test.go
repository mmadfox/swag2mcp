package service

import (
	"context"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSearchService_Search(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		query   string
		limit   int
		results []*model.Endpoint
		spec    *model.Spec
		coll    *model.Collection
		tag     *model.Tag
		wantErr bool
	}{
		{
			name:  "found",
			query: "pet", limit: 10,
			results: []*model.Endpoint{{ID: "ep1", SpecID: "s1", CollectionID: "c1", TagID: "t1", Name: "GET", Path: "/pet", Operation: &spec.Operation{Summary: "get pet"}}},
			spec:    &model.Spec{ID: "s1", Domain: "pets"},
			coll:    &model.Collection{ID: "c1", Title: "Pet Store"},
			tag:     &model.Tag{ID: "t1", Name: "pets"},
		},
		{name: "no results", query: "zzz", limit: 5, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			idx := NewMockIndexReader(ctrl)

			if tt.results != nil {
				idx.EXPECT().Search(gomock.Any(), gomock.Any(), tt.limit).Return(tt.results, nil)
				idx.EXPECT().SpecByID(tt.results[0].SpecID).Return(tt.spec, nil)
				idx.EXPECT().CollectionByID(tt.results[0].CollectionID).Return(tt.coll, nil)
				idx.EXPECT().TagByID(tt.results[0].TagID).Return(tt.tag, nil)
			} else {
				idx.EXPECT().Search(gomock.Any(), gomock.Any(), tt.limit).Return(nil, errNotFound("search", ""))
			}

			svc := newSearchService(idx, fakeValidator{})
			resp, err := svc.Search(context.Background(), SearchRequest{Query: tt.query, Limit: tt.limit})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, resp.Endpoints, len(tt.results))
		})
	}
}

func TestSearchService_Search_validationError(t *testing.T) {
	t.Parallel()

	svc := newSearchService(NewMockIndexReader(gomock.NewController(t)), strictValidator{})
	_, err := svc.Search(context.Background(), SearchRequest{})
	require.Error(t, err)
}

func TestSearchService_Search_specNotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().Search(gomock.Any(), gomock.Any(), 10).Return(
		[]*model.Endpoint{{ID: "ep1", SpecID: "s1", CollectionID: "c1", TagID: "t1", Name: "GET", Path: "/pet", Operation: &spec.Operation{}}}, nil)
	idx.EXPECT().SpecByID("s1").Return(nil, errNotFound("spec", "s1"))

	svc := newSearchService(idx, fakeValidator{})
	_, err := svc.Search(context.Background(), SearchRequest{Query: "pet", Limit: 10})
	require.Error(t, err)
}

func TestSearchService_Search_collectionNotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().Search(gomock.Any(), gomock.Any(), 10).Return(
		[]*model.Endpoint{{ID: "ep1", SpecID: "s1", CollectionID: "c1", TagID: "t1", Name: "GET", Path: "/pet", Operation: &spec.Operation{}}}, nil)
	idx.EXPECT().SpecByID("s1").Return(&model.Spec{ID: "s1", Domain: "pets"}, nil)
	idx.EXPECT().CollectionByID("c1").Return(nil, errNotFound("collection", "c1"))

	svc := newSearchService(idx, fakeValidator{})
	_, err := svc.Search(context.Background(), SearchRequest{Query: "pet", Limit: 10})
	require.Error(t, err)
}

func TestSearchService_Search_tagNotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().Search(gomock.Any(), gomock.Any(), 10).Return(
		[]*model.Endpoint{{ID: "ep1", SpecID: "s1", CollectionID: "c1", TagID: "t1", Name: "GET", Path: "/pet", Operation: &spec.Operation{}}}, nil)
	idx.EXPECT().SpecByID("s1").Return(&model.Spec{ID: "s1", Domain: "pets"}, nil)
	idx.EXPECT().CollectionByID("c1").Return(&model.Collection{ID: "c1", Title: "Pet Store"}, nil)
	idx.EXPECT().TagByID("t1").Return(nil, errNotFound("tag", "t1"))

	svc := newSearchService(idx, fakeValidator{})
	_, err := svc.Search(context.Background(), SearchRequest{Query: "pet", Limit: 10})
	require.Error(t, err)
}
