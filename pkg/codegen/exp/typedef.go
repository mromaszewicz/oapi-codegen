package exp

import (
	"fmt"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
)

type TypeDefinition struct {
	// SchemaPath is the path of the traversal down the OpenAPI 3 DOM to reach
	// the given schema. It is like the $ref path, but goes deeper based on
	// our own traversal of composite types.
	SchemaPath string `yaml:"schema_path"`
	// The RawSchema will point to any top level schema in the OpenAPI 3 spec, such
	// as regular Schemas, Parameters, etc.
	RawSchema any `yaml:"raw_schema"`
}

func GenerateTypeDefinitions(nodes PathTreeNodesByPath, cfg Configuration) ([]TypeDefinition, error) {
	var typeDefinitions []TypeDefinition
	for _, nodePath := range SortedMapKeys(nodes) {
		node := nodes[nodePath]

		// If the node refers to a reference, we do not traverse deeper, as
		// we will generate types for the reference definition, not for its
		// use.
		if node.Ref != "" {
			continue
		}

		// The RawSchema field will be a top level schema of any kind, SchemaRef,
		// ParameterRef, etc.
		switch concreteSchema := node.Schema.(type) {
		case *openapi3.SchemaRef:
			tds, err := schemaRefToTypeDefinitions(nodePath, concreteSchema)
			if err != nil {
				return nil, fmt.Errorf("generating type definitions for openapi3.SchemaRef at path (%s): %w", nodePath, err)
			}
			typeDefinitions = append(typeDefinitions, tds...)
		case *openapi3.ParameterRef:
		case *openapi3.RequestBodyRef:
		case *openapi3.ResponseRef:

		default:
			return nil, fmt.Errorf("unhandled schema type in spec: %v", reflect.TypeOf(concreteSchema))
		}
	}

	return typeDefinitions, nil
}

// schemaRefToTypeDefinitions creates a list of type definitions for Go based
// on a SchemaRef input. The mapping is 1:N not 1:1 because we may need to generate
// additional type definitions for inner types, such as when creating an additionalProperties
// helper object, which can not be inlined anonymously, or when the user has called for
// specific instantiation via annotation or config file override.
func schemaRefToTypeDefinitions(schemaPath string, sRef *openapi3.SchemaRef) ([]TypeDefinition, error) {
	return []TypeDefinition{{
		SchemaPath: schemaPath,
	}}, nil
}
