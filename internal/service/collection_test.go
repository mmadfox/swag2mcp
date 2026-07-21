package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCollectionService_CollectionsBySpec(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		specID  string
		spec    *model.Spec
		colls   []*model.Collection
		wantErr bool
	}{
		{name: "found", specID: "s1", spec: &model.Spec{ID: "s1", Domain: "d1"}, colls: []*model.Collection{{ID: "c1", Title: "C1"}}},
		{name: "not found", specID: "missing", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			idx := NewMockIndexReader(ctrl)

			if tt.spec != nil {
				idx.EXPECT().SpecByID(tt.specID).Return(tt.spec, nil)
				idx.EXPECT().CollectionsBySpec(tt.specID).Return(tt.colls, nil)
			} else {
				idx.EXPECT().SpecByID(tt.specID).Return(nil, errNotFound("spec", tt.specID))
			}

			svc := newCollectionService(idx, fakeValidator{})
			resp, err := svc.CollectionsBySpec(context.Background(), CollectionsRequest{SpecID: tt.specID})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, resp.Collections, len(tt.colls))
		})
	}
}

func TestCollectionService_CollectionByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		collID  string
		coll    *model.Collection
		spec    *model.Spec
		tags    []*model.Tag
		wantErr bool
	}{
		{name: "found", collID: "c1", coll: &model.Collection{ID: "c1", SpecID: "s1", Title: "C1"}, spec: &model.Spec{ID: "s1", Domain: "d1"}, tags: []*model.Tag{{ID: "t1", Name: "tag1"}}},
		{name: "not found", collID: "missing", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			idx := NewMockIndexReader(ctrl)

			if tt.coll != nil {
				idx.EXPECT().CollectionByID(tt.collID).Return(tt.coll, nil)
				idx.EXPECT().SpecByID(tt.coll.SpecID).Return(tt.spec, nil)
				idx.EXPECT().TagsByCollection(tt.collID).Return(tt.tags, nil)
			} else {
				idx.EXPECT().CollectionByID(tt.collID).Return(nil, errNotFound("collection", tt.collID))
			}

			svc := newCollectionService(idx, fakeValidator{})
			resp, err := svc.CollectionByID(context.Background(), CollectionByIDRequest{ID: tt.collID})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.coll.ID, resp.Collection.ID)
			require.Len(t, resp.Tags, len(tt.tags))
		})
	}
}
