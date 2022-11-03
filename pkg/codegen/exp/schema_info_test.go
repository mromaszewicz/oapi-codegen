package exp

import "testing"

func TestGenerateSchemaInfosFromTree(t *testing.T) {
	spec := MustLoad(t, "all_components.yaml")
	root := BuildSchemaTree(spec)
	si := GenerateSchemaInfosFromTree(root)

	for _, pathName := range SortedMapKeys(si) {
		t.Log(pathName)
	}
}
