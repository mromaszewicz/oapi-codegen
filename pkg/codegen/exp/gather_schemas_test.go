package exp

import (
	"embed"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed test_files
var testFiles embed.FS

func MustRead(t *testing.T, filename string) []byte {
	buf, err := testFiles.ReadFile("test_files/" + filename)
	require.NoErrorf(t, err, "loading test file")
	return buf
}

func MustLoad(t *testing.T, filename string) *openapi3.T {
	buf := MustRead(t, filename)
	loader := openapi3.NewLoader()
	spec, err := loader.LoadFromData(buf)
	require.NoErrorf(t, err, "parsing spec from buffer")
	return spec
}

func TestGatherComponents(t *testing.T) {
	spec := MustLoad(t, "all_components.yaml")
	root := BuildSchemaTree(spec)

	// Ensure that we have some properly traversed the data. Generally, once
	// we unwrap a top level container, like Header or Parameter, and have a SchemaRef,
	// the nested property recursion will work the same for all types, so we only
	// need to test the recursion on one kind of object, but ensure that the rest
	// unwrap correctly.

	// Make sure that objects and their inner properties turn into paths
	assert.NotNil(t, root.GetNodeByPath("components/schemas/SimpleObject"))
	assert.NotNil(t, root.GetNodeByPath("components/schemas/SimpleObject/Name"))
	assert.NotNil(t, root.GetNodeByPath("components/schemas/SimpleObject/Color"))
	assert.NotNil(t, root.GetNodeByPath("components/schemas/ObjectWithAnonymousType"))
	assert.NotNil(t, root.GetNodeByPath("components/schemas/ObjectWithAnonymousType/Name"))
	assert.NotNil(t, root.GetNodeByPath("components/schemas/ObjectWithAnonymousType/CustomProperty"))
	assert.NotNil(t, root.GetNodeByPath("components/schemas/ObjectWithAnonymousType/CustomProperty/Field1"))
	assert.NotNil(t, root.GetNodeByPath("components/schemas/ObjectWithAnonymousType/CustomProperty/Field2"))
	assert.NotNil(t, root.GetNodeByPath("components/schemas/ObjectWithAnonymousType/CustomProperty/Field2/Prop1"))
	assert.NotNil(t, root.GetNodeByPath("components/schemas/ObjectWithAnonymousType/CustomProperty/Field2/Prop2"))

	// Make sure the types are correct deep down.
	assert.Equal(t, "object", root.GetNodeByPath("components/schemas/ObjectWithAnonymousType/CustomProperty/Field2").Schema.(*openapi3.SchemaRef).Value.Type)
	assert.Equal(t, "number", root.GetNodeByPath("components/schemas/ObjectWithAnonymousType/CustomProperty/Field2/Prop2").Schema.(*openapi3.SchemaRef).Value.Type)

	// Requests and responses need to have the content type in their path elements, so let's
	// ensure that string based search works.
	assert.Equal(t,
		root.GetNodeByPathElements([]string{"components", "requestBodies", "MultiContent", "application/json"}),
		root.GetNodeByPath("components/requestBodies/MultiContent/application/json"))
	assert.Equal(t,
		root.GetNodeByPathElements([]string{"components", "responses", "MultiContent", "application/json"}),
		root.GetNodeByPath("components/responses/MultiContent/application/json"))

}
