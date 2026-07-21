package service

import (
	"testing"

	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/stretchr/testify/require"
)

func TestIndexer_EnsureIndex(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	sp := &model.Spec{ID: "s1", Domain: "test"}
	colls := map[string]*model.Collection{"c1": {ID: "c1", SpecID: "s1"}}
	tags := map[string]*model.Tag{"t1": {ID: "t1", SpecID: "s1", CollectionID: "c1"}}
	eps := map[string]*model.Endpoint{
		"ep1": {ID: "ep1", SpecID: "s1", CollectionID: "c1", TagID: "t1", Operation: &spec.Operation{Summary: "test"}},
	}

	err = svc.indexSpec(sp, colls, tags, eps)
	require.NoError(t, err)

	specs := svc.index.AllSpecs()
	require.Len(t, specs, 1)
	require.Equal(t, "s1", specs[0].ID)
}
