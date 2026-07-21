package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestResponseService_ResponseOutline_validationError(t *testing.T) {
	t.Parallel()

	svc := newResponseService(newServiceContext(), NewMockWorkspaceOps(gomock.NewController(t)), strictValidator{})
	_, err := svc.ResponseOutline(context.Background(), ResponseOutlineRequest{})
	require.Error(t, err)
}

func TestResponseService_ResponseCompress_validationError(t *testing.T) {
	t.Parallel()

	svc := newResponseService(newServiceContext(), NewMockWorkspaceOps(gomock.NewController(t)), strictValidator{})
	_, err := svc.ResponseCompress(context.Background(), ResponseCompressRequest{})
	require.Error(t, err)
}

func TestResponseService_ResponseSlice_validationError(t *testing.T) {
	t.Parallel()

	svc := newResponseService(newServiceContext(), NewMockWorkspaceOps(gomock.NewController(t)), strictValidator{})
	_, err := svc.ResponseSlice(context.Background(), ResponseSliceRequest{})
	require.Error(t, err)
}
