package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestToolsService_MakeToolDefinitions(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().AllSpecs().Return(nil)

	svc := newToolsService(idx, func() bool { return false })
	defs, err := svc.MakeToolDefinitions()
	require.NoError(t, err)
	require.NotEmpty(t, defs.Instruction)
	require.NotEmpty(t, defs.Tools)
}

func TestToolsService_MakeToolDefinitions_authDisabled(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().AllSpecs().Return(nil)

	svc := newToolsService(idx, func() bool { return true })
	defs, err := svc.MakeToolDefinitions()
	require.NoError(t, err)

	for _, tool := range defs.Tools {
		require.NotEqual(t, Auth, tool.Name)
	}
}

func TestLoadToolFromEmbed(t *testing.T) {
	t.Parallel()

	tool, err := loadToolFromEmbed("spec_list.md")
	require.NoError(t, err)
	require.Equal(t, SpecList, tool.Name)
	require.NotEmpty(t, tool.Description)
}

func TestLoadInstructionFromEmbed(t *testing.T) {
	t.Parallel()

	instruction, err := loadInstructionFromEmbed()
	require.NoError(t, err)
	require.NotEmpty(t, instruction)
}
