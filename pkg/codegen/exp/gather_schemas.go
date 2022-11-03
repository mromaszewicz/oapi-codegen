package exp

import (
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
)

// GatherSchemasFromComponents gathers all the explicitly named schemas in the
// components section of an OpenAPI specification.
func GatherSchemasFromComponents(components *openapi3.Components) []*PathTreeNode {
	// #/components/schemas
	var result []*PathTreeNode
	if len(components.Schemas) > 0 {
		schemasNode := &PathTreeNode{PathElement: "schemas", ElementName: "Schema"}
		for _, schemaName := range SortedMapKeys(components.Schemas) {
			schemaRef := components.Schemas[schemaName]
			childNode := &PathTreeNode{
				PathElement: schemaName,
				Schema:      schemaRef,
				Ref:         schemaRef.Ref,
			}
			gatherNestedSchemas(childNode, schemaRef)
			schemasNode.AddChild(childNode)
		}
		result = append(result, schemasNode)
	}

	// #/components/parameters
	if len(components.Parameters) > 0 {
		parametersNode := &PathTreeNode{PathElement: "parameters", ElementName: "Parameter"}
		for _, paramName := range SortedMapKeys(components.Parameters) {
			paramRef := components.Parameters[paramName]
			childNode := &PathTreeNode{
				PathElement: paramName,
				Schema:      paramRef,
				Ref:         paramRef.Ref,
			}
			gatherNestedSchemas(childNode, paramRef.Value.Schema)
			parametersNode.AddChild(childNode)
		}
		result = append(result, parametersNode)
	}

	// #/components/headers
	if len(components.Headers) > 0 {
		headersNode := &PathTreeNode{PathElement: "headers", ElementName: "Header"}
		for _, headerName := range SortedMapKeys(components.Headers) {
			headerRef := components.Headers[headerName]
			childNode := &PathTreeNode{
				PathElement: headerName,
				Schema:      headerRef,
				Ref:         headerRef.Ref,
			}
			gatherNestedSchemas(childNode, headerRef.Value.Schema)
			headersNode.AddChild(childNode)
		}
		result = append(result, headersNode)
	}

	// #/components/requestBodies. This one is more complex, as we have multiple
	// schemas possible within the map of content types. We will have a content
	// type as a path element, which can contain a slash, but the tree structure
	// can express this unambiguously. It will be tricky to search by string versus
	// by path elements.
	if len(components.RequestBodies) > 0 {
		requestBodiesNode := &PathTreeNode{PathElement: "requestBodies", ElementName: "RequestBody"}
		for _, rbName := range SortedMapKeys(components.RequestBodies) {
			rbRef := components.RequestBodies[rbName]
			childNode := &PathTreeNode{
				PathElement: rbName,
				Schema:      rbRef,
				Ref:         rbRef.Ref,
			}

			for _, contentTypeName := range SortedMapKeys(rbRef.Value.Content) {
				mediaType := rbRef.Value.Content[contentTypeName]
				contentNode := &PathTreeNode{
					PathElement: contentTypeName,
					Schema:      mediaType.Schema,
					Ref:         mediaType.Schema.Ref,
				}
				gatherNestedSchemas(contentNode, mediaType.Schema)
				childNode.AddChild(contentNode)
			}
			requestBodiesNode.AddChild(childNode)
		}
		result = append(result, requestBodiesNode)
	}

	// #/components/responses
	if len(components.Responses) > 0 {
		node := &PathTreeNode{PathElement: "responses", ElementName: "Response"}
		for _, responseName := range SortedMapKeys(components.Responses) {
			responseRef := components.Responses[responseName]
			childNode := &PathTreeNode{
				PathElement: responseName,
				Schema:      responseRef,
				Ref:         responseRef.Ref,
			}

			for _, contentTypeName := range SortedMapKeys(responseRef.Value.Content) {
				mediaType := responseRef.Value.Content[contentTypeName]
				contentNode := &PathTreeNode{
					PathElement: contentTypeName,
					Schema:      mediaType.Schema,
					Ref:         mediaType.Schema.Ref,
				}
				gatherNestedSchemas(contentNode, mediaType.Schema)
				childNode.AddChild(contentNode)
			}
			node.AddChild(childNode)
		}
		result = append(result, node)
	}

	// #/components/securitySchemes
	if len(components.SecuritySchemes) > 0 {
		node := &PathTreeNode{PathElement: "securitySchemes", ElementName: "SecurityScheme"}
		for ssName, ssRef := range components.SecuritySchemes {
			childNode := &PathTreeNode{
				PathElement: ssName,
				Schema:      ssRef,
				Ref:         ssRef.Ref,
			}
			node.AddChild(childNode)
		}
		result = append(result, node)
	}

	return result
}

// gatherNestedSchemas traverses a SchemaRef, and adds new nodes into
// the given parentNode based on the properties of the schema. The function
// is recursive and bottoms out when it hits a Ref, or hits the end of
// the schema definition structure.
func gatherNestedSchemas(parentNode *PathTreeNode, schemaRef *openapi3.SchemaRef) {
	if len(schemaRef.Value.Properties) > 0 {
		for propName, propRef := range schemaRef.Value.Properties {
			node := &PathTreeNode{
				PathElement: propName,
				Schema:      propRef,
				Ref:         propRef.Ref,
			}
			if propRef.Ref == "" {
				gatherNestedSchemas(node, propRef)
			}
			parentNode.AddChild(node)
		}
		gatherAdditionalProperties(parentNode, schemaRef)
	} else if len(schemaRef.Value.AllOf) > 0 {
		// For allOf, we're going to create a node named "allOf", and children
		// will have names based on their index, eg, [0], so in the path, we
		// will have .../allOf/0/..., etc.
		allOfNode := &PathTreeNode{PathElement: "allOf"}
		parentNode.AddChild(allOfNode)

		for i := range schemaRef.Value.AllOf {
			sref := schemaRef.Value.AllOf[i]
			node := &PathTreeNode{
				PathElement: strconv.Itoa(i),
				Schema:      sref,
				Ref:         sref.Ref,
			}
			gatherNestedSchemas(node, sref)
			gatherAdditionalProperties(node, sref)
			allOfNode.AddChild(node)
		}
	} else if len(schemaRef.Value.OneOf) > 0 {
		oneOfNode := &PathTreeNode{PathElement: "oneOf"}
		parentNode.AddChild(oneOfNode)

		for i := range schemaRef.Value.OneOf {
			sref := schemaRef.Value.OneOf[i]
			node := &PathTreeNode{
				PathElement: strconv.Itoa(i),
				Schema:      sref,
				Ref:         sref.Ref,
			}
			gatherNestedSchemas(node, sref)
			gatherAdditionalProperties(node, sref)
			oneOfNode.AddChild(node)
		}
	} else if len(schemaRef.Value.AnyOf) > 0 {
		anyOfNode := &PathTreeNode{PathElement: "anyOf"}
		parentNode.AddChild(anyOfNode)

		for i := range schemaRef.Value.AnyOf {
			sref := schemaRef.Value.AnyOf[i]
			node := &PathTreeNode{
				PathElement: strconv.Itoa(i),
				Schema:      sref,
				Ref:         sref.Ref,
			}
			gatherNestedSchemas(node, sref)
			gatherAdditionalProperties(node, sref)
			anyOfNode.AddChild(node)
		}
	}
}

// gatherParameters creates nodes from the specified parameters and adds them
// as children to the specified parent node.
func gatherParameters(p openapi3.Parameters, parent *PathTreeNode) {
	for _, parameterSpec := range p {
		// We skip Refs to parameters, since we only care about new definitions.
		if parameterSpec.Ref != "" {
			continue
		}
		parameterNode := &PathTreeNode{
			PathElement: parameterSpec.Value.Name,
			ElementName: parameterSpec.Value.Name,
			Schema:      parameterSpec,
		}
		gatherNestedSchemas(parameterNode, parameterSpec.Value.Schema)
		parent.AddChild(parameterNode)
	}
}

func gatherAdditionalProperties(parentNode *PathTreeNode, schemaRef *openapi3.SchemaRef) {
	s := schemaRef.Value

	hasAdditionalProperties := false

	if s.AdditionalProperties.Has != nil && *s.AdditionalProperties.Has {
		hasAdditionalProperties = true
	}
	if s.AdditionalProperties.Schema != nil {
		hasAdditionalProperties = true
	}

	if hasAdditionalProperties {
		addPropNode := &PathTreeNode{
			PathElement: "additionalProperties",
			Schema:      s.AdditionalProperties.Schema,
			Ref:         s.AdditionalProperties.Schema.Ref,
		}
		parentNode.AddChild(addPropNode)
		if s.AdditionalProperties.Schema != nil {
			gatherNestedSchemas(addPropNode, s.AdditionalProperties.Schema)
		}
	}
}

// BuildSchemaTree walks all the schemas in the spec, and builds a tree,
// based on path in the spec, that maps the path to the specific schema description
func BuildSchemaTree(spec *openapi3.T) *PathTreeNode {
	root := &PathTreeNode{}

	if spec.Components != nil {
		componentsNode := &PathTreeNode{PathElement: "components"}
		components := GatherSchemasFromComponents(spec.Components)
		for _, c := range components {
			componentsNode.AddChild(c)
		}
		root.AddChild(componentsNode)
	}

	if spec.Paths.Len() > 0 {
		allPathsNode := &PathTreeNode{PathElement: "paths"}
		for _, path := range SortedMapKeys(spec.Paths.Map()) {
			pathItem := spec.Paths.Map()[path]
			// The paths node above will contain a child identified by
			// the path URI. Since we use / as a delimiter, we will url escape
			// the actual path when embedding it in our schema path.
			pathNode := &PathTreeNode{
				PathElement: path,
			}
			allPathsNode.AddChild(pathNode)

			// Parameters may be defined outside http schemes, so we'll
			// put them under a child named "parameters", eg, parameters
			// for path /foo/{x} will be /foo/{x}/parameters/x
			if len(pathItem.Parameters) > 0 {
				pathParametersNode := &PathTreeNode{
					PathElement: "parameters",
					ElementName: "Parameter",
				}
				pathNode.AddChild(pathParametersNode)

				gatherParameters(pathItem.Parameters, pathParametersNode)
			}

			// Now, we go through all the operations under the path.
			for _, opScheme := range SortedMapKeys(pathItem.Operations()) {
				operation := pathItem.Operations()[opScheme]
				operationNode := &PathTreeNode{
					PathElement: operation.OperationID,
					// PathElement: opScheme,
					// ElementName: operation.OperationID,
				}
				pathNode.AddChild(operationNode)

				// Operations have many nested types, so we'll also convert them
				// into paths.
				if len(operation.Parameters) > 0 {
					parametersNode := &PathTreeNode{
						PathElement: "parameters",
						ElementName: "Parameter",
					}
					gatherParameters(operation.Parameters, parametersNode)
					operationNode.AddChild(parametersNode)
				}

				if operation.Responses.Len() > 0 {
					responsesNode := &PathTreeNode{PathElement: "responses"}
					operationNode.AddChild(responsesNode)

					for _, responseCode := range SortedMapKeys(operation.Responses.Map()) {
						responseSpec := operation.Responses.Map()[responseCode]
						responseNode := &PathTreeNode{
							PathElement: responseCode,
							Schema:      responseSpec,
						}
						responsesNode.AddChild(responseNode)
					}

				}

				if operation.RequestBody != nil {

				}
			}

		}
		root.AddChild(allPathsNode)
	}

	return root
}
