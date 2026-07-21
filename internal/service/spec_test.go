package service

import (
	"context"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSpecService_Specs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		specs   []*model.Spec
		wantLen int
	}{
		{name: "empty", specs: nil, wantLen: 0},
		{name: "multiple", specs: []*model.Spec{{ID: "a", Domain: "x"}, {ID: "b", Domain: "y"}}, wantLen: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			idx := NewMockIndexReader(ctrl)
			idx.EXPECT().AllSpecs().Return(tt.specs)

			svc := newSpecService(idx, fakeValidator{})
			resp, err := svc.Specs(context.Background())
			require.NoError(t, err)
			require.Len(t, resp.Specs, tt.wantLen)
		})
	}
}

func TestSpecService_SpecByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		specID      string
		spec        *model.Spec
		collections []*model.Collection
		wantErr     bool
	}{
		{
			name:        "found",
			specID:      "s1",
			spec:        &model.Spec{ID: "s1", Domain: "d1"},
			collections: []*model.Collection{{ID: "c1", Title: "C1"}},
		},
		{name: "not found", specID: "missing", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			idx := NewMockIndexReader(ctrl)

			if tt.spec != nil {
				idx.EXPECT().SpecByID(tt.specID).Return(tt.spec, nil)
				idx.EXPECT().CollectionsBySpec(tt.specID).Return(tt.collections, nil)
			} else {
				idx.EXPECT().SpecByID(tt.specID).Return(nil, errNotFound("spec", tt.specID))
			}

			svc := newSpecService(idx, fakeValidator{})
			resp, err := svc.SpecByID(context.Background(), SpecByIDRequest{ID: tt.specID})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.spec.ID, resp.Spec.ID)
		})
	}
}
