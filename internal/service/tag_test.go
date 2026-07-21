package service

import (
	"context"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTagService_TagByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tagID   string
		tag     *model.Tag
		wantErr bool
	}{
		{name: "found", tagID: "t1", tag: &model.Tag{ID: "t1", Name: "mytag"}},
		{name: "not found", tagID: "missing", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			idx := NewMockIndexReader(ctrl)

			if tt.tag != nil {
				idx.EXPECT().TagByID(tt.tagID).Return(tt.tag, nil)
			} else {
				idx.EXPECT().TagByID(tt.tagID).Return(nil, errNotFound("tag", tt.tagID))
			}

			svc := newTagService(idx, fakeValidator{})
			resp, err := svc.TagByID(context.Background(), TagByIDRequest{ID: tt.tagID})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.tag.Name, resp.Tag.Title)
		})
	}
}

func TestTagService_TagsBySpec(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		specID  string
		tags    []*model.Tag
		wantErr bool
	}{
		{name: "found", specID: "s1", tags: []*model.Tag{{ID: "t1", Name: "t1"}}},
		{name: "not found", specID: "missing", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			idx := NewMockIndexReader(ctrl)

			if tt.tags != nil {
				idx.EXPECT().TagsBySpec(tt.specID).Return(tt.tags, nil)
			} else {
				idx.EXPECT().TagsBySpec(tt.specID).Return(nil, errNotFound("spec", tt.specID))
			}

			svc := newTagService(idx, fakeValidator{})
			resp, err := svc.TagsBySpec(context.Background(), TagsBySpecRequest{SpecID: tt.specID})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, resp.Tags, len(tt.tags))
		})
	}
}

func TestTagService_TagsByCollection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		collectionID string
		coll         *model.Collection
		spec         *model.Spec
		tags         []*model.Tag
		wantErr      bool
	}{
		{name: "found", collectionID: "c1", coll: &model.Collection{ID: "c1", SpecID: "s1", Title: "C1"}, spec: &model.Spec{ID: "s1", Domain: "d1"}, tags: []*model.Tag{{ID: "t1", Name: "t1"}}},
		{name: "not found", collectionID: "missing", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			idx := NewMockIndexReader(ctrl)

			if tt.coll != nil {
				idx.EXPECT().CollectionByID(tt.collectionID).Return(tt.coll, nil)
				idx.EXPECT().SpecByID(tt.coll.SpecID).Return(tt.spec, nil)
				idx.EXPECT().TagsByCollection(tt.collectionID).Return(tt.tags, nil)
			} else {
				idx.EXPECT().CollectionByID(tt.collectionID).Return(nil, errNotFound("collection", tt.collectionID))
			}

			svc := newTagService(idx, fakeValidator{})
			resp, err := svc.TagsByCollection(context.Background(), TagsByCollectionRequest{CollectionID: tt.collectionID})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, resp.Tags, len(tt.tags))
		})
	}
}
