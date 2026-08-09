/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

import (
	"log/slog"
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

	svc := newToolsService(idx, func() bool { return false }, slog.Default())
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

	svc := newToolsService(idx, func() bool { return true }, slog.Default())
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

func TestInstructionMentionsAuthInjectedParameters(t *testing.T) {
	t.Parallel()

	instruction, err := loadInstructionFromEmbed()
	require.NoError(t, err)
	require.Contains(t, instruction, "Auth-injected parameters are automatic")
	require.Contains(t, instruction, "api_key")
	require.Contains(t, instruction, "timestamp")
	require.Contains(t, instruction, "signature")
}

func TestInvokeToolMentionsAuthInjectedParameters(t *testing.T) {
	t.Parallel()

	tool, err := loadToolFromEmbed("invoke.md")
	require.NoError(t, err)
	require.Contains(t, tool.Description, "Auth-injected parameters are automatic")
	require.Contains(t, tool.Description, "api_key")
}

func TestInspectToolMentionsAuthInjectedParameters(t *testing.T) {
	t.Parallel()

	tool, err := loadToolFromEmbed("inspect.md")
	require.NoError(t, err)
	require.Contains(t, tool.Description, "auth-injected parameters")
	require.Contains(t, tool.Description, "api_key")
}

func TestInstructionMentionsInvokePacingAndRetry(t *testing.T) {
	t.Parallel()

	instruction, err := loadInstructionFromEmbed()
	require.NoError(t, err)
	require.Contains(t, instruction, "Invoke one at a time, with bounded retries")
	require.Contains(t, instruction, "one outstanding invoke")
	require.Contains(t, instruction, "rate limited")
	require.Contains(t, instruction, "at most")
}

func TestInvokeToolMentionsPacingAndRetry(t *testing.T) {
	t.Parallel()

	tool, err := loadToolFromEmbed("invoke.md")
	require.NoError(t, err)
	require.Contains(t, tool.Description, "One at a time, bounded retries")
	require.Contains(t, tool.Description, "one outstanding invoke")
	require.Contains(t, tool.Description, "rate limited")
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

	svc := newToolsService(idx, func() bool { return false }, slog.Default())
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

	svc := newToolsService(idx, func() bool { return false }, slog.Default())
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

	svc := newToolsService(idx, func() bool { return false }, slog.Default())
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

	svc := newToolsService(idx, func() bool { return false }, slog.Default())
	result := svc.makeAvailableSpecs()
	require.Contains(t, result, "more than 10 collections")
}
