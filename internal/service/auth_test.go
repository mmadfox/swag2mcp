package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type errAuthClient struct{}

func (errAuthClient) New() error { return nil }

func (errAuthClient) Type() auth.Type { return auth.NoAuth }

func (errAuthClient) Apply(_ *http.Request, _ *auth.Info) error {
	return errors.New("auth apply failed")
}

func (errAuthClient) Validate() error { return nil }

func TestAuthService_Auth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		specID    string
		spec      *model.Spec
		disabled  bool
		wantToken bool
		wantErr   bool
	}{
		{name: "disabled", specID: "s1", disabled: true, wantToken: false},
		{name: "not found", specID: "missing", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			idx := NewMockIndexReader(ctrl)

			if tt.wantErr {
				idx.EXPECT().SpecByID(tt.specID).Return(nil, errNotFound("spec", tt.specID))
			}

			svc := newAuthService(idx, func() bool { return tt.disabled })
			resp, err := svc.Auth(context.Background(), AuthRequest{SpecID: tt.specID})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.disabled {
				require.Empty(t, resp.Token)
			}
		})
	}
}

func TestAuthService_Auth_nilAuth(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().SpecByID("s1").Return(&model.Spec{ID: "s1"}, nil)

	svc := newAuthService(idx, func() bool { return false })
	resp, err := svc.Auth(context.Background(), AuthRequest{SpecID: "s1"})
	require.NoError(t, err)
	require.Empty(t, resp.Token)
}

func TestAuthService_Auth_applyError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().SpecByID("s1").Return(&model.Spec{ID: "s1", Auth: errAuthClient{}}, nil)

	svc := newAuthService(idx, func() bool { return false })
	_, err := svc.Auth(context.Background(), AuthRequest{SpecID: "s1"})
	require.Error(t, err)
}
