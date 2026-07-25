package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"testing"

	"github.com/mmadfox/swag2mcp/internal/model"
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

func TestMakeAvailableSpecs_withSpecs(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().AllSpecs().Return([]*model.Spec{
		{ID: "s1", Domain: "api.example.com", LLMTitle: "Example API"},
	})
	idx.EXPECT().CollectionsBySpec("s1").Return([]*model.Collection{
		{ID: "c1", LLMTitle: "Users"},
		{ID: "c2", LLMTitle: "Orders"},
	}, nil)

	svc := newToolsService(idx, func() bool { return false })
	result := svc.makeAvailableSpecs()
	require.Contains(t, result, "api.example.com")
	require.Contains(t, result, "Users")
	require.Contains(t, result, "Orders")
}

func TestMakeAvailableSpecs_withInstruction(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().AllSpecs().Return([]*model.Spec{
		{ID: "s1", Domain: "api.example.com", LLMTitle: "API", LLMInstruction: "Use with care"},
	})
	idx.EXPECT().CollectionsBySpec("s1").Return(nil, errNotFound("collections", "s1"))

	svc := newToolsService(idx, func() bool { return false })
	result := svc.makeAvailableSpecs()
	require.Contains(t, result, "Use with care")
	require.Contains(t, result, "No available collections")
}

func TestMakeAvailableSpecs_withCollectionInstruction(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().AllSpecs().Return([]*model.Spec{
		{ID: "s1", Domain: "api.example.com", LLMTitle: "API"},
	})
	idx.EXPECT().CollectionsBySpec("s1").Return([]*model.Collection{
		{ID: "c1", LLMTitle: "Users", LLMInstruction: "User management"},
	}, nil)

	svc := newToolsService(idx, func() bool { return false })
	result := svc.makeAvailableSpecs()
	require.Contains(t, result, "User management")
}

func TestMakeAvailableSpecs_moreThan10Collections(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().AllSpecs().Return([]*model.Spec{
		{ID: "s1", Domain: "api.example.com", LLMTitle: "API"},
	})
	colls := make([]*model.Collection, 15)
	for i := range colls {
		colls[i] = &model.Collection{ID: "c", LLMTitle: "C"}
	}
	idx.EXPECT().CollectionsBySpec("s1").Return(colls, nil)

	svc := newToolsService(idx, func() bool { return false })
	result := svc.makeAvailableSpecs()
	require.Contains(t, result, "more than 10 collections")
}
