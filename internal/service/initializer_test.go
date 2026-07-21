package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitializer_InitWorkspace_sameDir(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	init := newInitializer(svc)
	err = init.initWorkspace(svc.ws.Root())
	require.NoError(t, err)
}

func TestInitializer_InitWorkspace_newDir(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	init := newInitializer(svc)
	tmpDir := t.TempDir()
	err = init.initWorkspace(tmpDir)
	require.NoError(t, err)
	require.Equal(t, tmpDir, svc.ws.Root())
}
