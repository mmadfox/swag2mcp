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

func TestInspectService_Inspect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		epID    string
		ep      *model.Endpoint
		spec    *model.Spec
		coll    *model.Collection
		wantErr bool
	}{
		{
			name: "found",
			epID: "ep1",
			ep:   &model.Endpoint{ID: "ep1", SpecID: "s1", CollectionID: "c1", TagID: "t1", Name: "GET", Path: "/pet", Operation: &spec.Operation{ID: "op1"}},
			spec: &model.Spec{ID: "s1", Domain: "pets", BaseURL: "https://api.example.com"},
			coll: &model.Collection{ID: "c1", Title: "Pet Store"},
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
			} else {
				idx.EXPECT().EndpointByID(tt.epID).Return(nil, errNotFound("endpoint", tt.epID))
			}

			svc := newInspectService(idx, fakeValidator{})
			resp, err := svc.Inspect(context.Background(), InspectRequest{EndpointID: tt.epID})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.ep.Name, resp.Method)
			require.Equal(t, tt.spec.BaseURL, resp.BaseURL)
		})
	}
}
