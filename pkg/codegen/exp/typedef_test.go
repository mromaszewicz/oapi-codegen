package exp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateTypeDefinitions(t *testing.T) {
	spec := MustLoad(t, "all_components.yaml")
	root := BuildSchemaTree(spec)
	nodesByPath := GenerateSchemaInfosFromTree(root)

	// Loading nothing should result in the default configuration.
	config, err := LoadConfiguration(nil)
	require.NoError(t, err)

	typeDefs, err := GenerateTypeDefinitions(nodesByPath, config)
	require.NoError(t, err)

	for _, td := range typeDefs {
		t.Log(td)
	}
}
