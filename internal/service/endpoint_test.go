package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestEndpointService_EndpointByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		epID    string
		ep      *model.Endpoint
		spec    *model.Spec
		coll    *model.Collection
		tag     *model.Tag
		wantErr bool
	}{
		{
			name: "found",
			epID: "ep1",
			ep:   &model.Endpoint{ID: "ep1", SpecID: "s1", CollectionID: "c1", TagID: "t1", Name: "GET", Path: "/test", Operation: &spec.Operation{Summary: "test ep"}},
			spec: &model.Spec{ID: "s1", Domain: "d1"},
			coll: &model.Collection{ID: "c1", Title: "C1"},
			tag:  &model.Tag{ID: "t1", Name: "tag1"},
		},
		{name: "not found", epID: "missing", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			idx := NewMockIndexReader(ctrl)

			if tt.ep != nil {
				idx.EXPECT().EndpointByID(tt.epID).Return(tt.ep, nil)
				idx.EXPECT().SpecByID(tt.ep.SpecID).Return(tt.spec, nil)
				idx.EXPECT().CollectionByID(tt.ep.CollectionID).Return(tt.coll, nil)
				idx.EXPECT().TagByID(tt.ep.TagID).Return(tt.tag, nil)
			} else {
				idx.EXPECT().EndpointByID(tt.epID).Return(nil, errNotFound("endpoint", tt.epID))
			}

			svc := newEndpointService(idx, fakeValidator{})
			resp, err := svc.EndpointByID(context.Background(), EndpointByIDRequest{ID: tt.epID})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.ep.Name, resp.Endpoint.Method)
		})
	}
}

func TestEndpointService_EndpointsByTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tagID   string
		tag     *model.Tag
		coll    *model.Collection
		spec    *model.Spec
		eps     []*model.Endpoint
		wantErr bool
	}{
		{
			name:  "found",
			tagID: "t1",
			tag:   &model.Tag{ID: "t1", CollectionID: "c1"},
			coll:  &model.Collection{ID: "c1", SpecID: "s1", Title: "C1"},
			spec:  &model.Spec{ID: "s1", Domain: "d1"},
			eps:   []*model.Endpoint{{ID: "ep1", Name: "GET", Path: "/test", Operation: &spec.Operation{Summary: "test"}}},
		},
		{name: "not found", tagID: "missing", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			idx := NewMockIndexReader(ctrl)

			if tt.tag != nil {
				idx.EXPECT().TagByID(tt.tagID).Return(tt.tag, nil)
				idx.EXPECT().CollectionByID(tt.tag.CollectionID).Return(tt.coll, nil)
				idx.EXPECT().SpecByID(tt.coll.SpecID).Return(tt.spec, nil)
				idx.EXPECT().EndpointsByTag(tt.tagID).Return(tt.eps, nil)
			} else {
				idx.EXPECT().TagByID(tt.tagID).Return(nil, errNotFound("tag", tt.tagID))
			}

			svc := newEndpointService(idx, fakeValidator{})
			resp, err := svc.EndpointsByTag(context.Background(), EndpointsByTagRequest{TagID: tt.tagID})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, resp.Endpoints, len(tt.eps))
		})
	}
}

func TestEndpointService_EndpointsByCollection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		collectionID string
		coll         *model.Collection
		spec         *model.Spec
		eps          []*model.Endpoint
		tag          *model.Tag
		wantErr      bool
	}{
		{
			name:         "found",
			collectionID: "c1",
			coll:         &model.Collection{ID: "c1", SpecID: "s1", Title: "C1"},
			spec:         &model.Spec{ID: "s1", Domain: "d1"},
			eps:          []*model.Endpoint{{ID: "ep1", TagID: "t1", Name: "GET", Path: "/test", Operation: &spec.Operation{Summary: "test"}}},
			tag:          &model.Tag{ID: "t1", Name: "tag1"},
		},
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
				idx.EXPECT().EndpointByCollection(tt.collectionID).Return(tt.eps, nil)
				idx.EXPECT().TagByID(tt.eps[0].TagID).Return(tt.tag, nil)
			} else {
				idx.EXPECT().CollectionByID(tt.collectionID).Return(nil, errNotFound("collection", tt.collectionID))
			}

			svc := newEndpointService(idx, fakeValidator{})
			resp, err := svc.EndpointsByCollection(context.Background(), EndpointsByCollectionRequest{CollectionID: tt.collectionID})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, resp.Endpoints, len(tt.eps))
		})
	}
}

func TestEndpointService_EndpointsBySpec(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		specID  string
		eps     []*model.Endpoint
		spec    *model.Spec
		coll    *model.Collection
		tag     *model.Tag
		wantErr bool
	}{
		{
			name:   "found",
			specID: "s1",
			eps:    []*model.Endpoint{{ID: "ep1", SpecID: "s1", CollectionID: "c1", TagID: "t1", Name: "GET", Path: "/test", Operation: &spec.Operation{Summary: "test"}}},
			spec:   &model.Spec{ID: "s1", Domain: "d1"},
			coll:   &model.Collection{ID: "c1", Title: "C1"},
			tag:    &model.Tag{ID: "t1", Name: "tag1"},
		},
		{name: "not found", specID: "missing", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			idx := NewMockIndexReader(ctrl)

			if tt.eps != nil {
				idx.EXPECT().EndpointsBySpec(tt.specID).Return(tt.eps, nil)
				idx.EXPECT().SpecByID(tt.eps[0].SpecID).Return(tt.spec, nil)
				idx.EXPECT().CollectionByID(tt.eps[0].CollectionID).Return(tt.coll, nil)
				idx.EXPECT().TagByID(tt.eps[0].TagID).Return(tt.tag, nil)
			} else {
				idx.EXPECT().EndpointsBySpec(tt.specID).Return(nil, errNotFound("spec", tt.specID))
			}

			svc := newEndpointService(idx, fakeValidator{})
			resp, err := svc.EndpointsBySpec(context.Background(), EndpointsBySpecRequest{SpecID: tt.specID})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, resp.Endpoints, len(tt.eps))
		})
	}
}
