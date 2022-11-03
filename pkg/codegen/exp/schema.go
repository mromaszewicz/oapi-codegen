package exp

import "github.com/getkin/kin-openapi/openapi3"

// SchemaAdapter wraps the openapi3 Schema and provides a bunch of useful
// utilities around it.
type SchemaAdapter struct {
	Schema *openapi3.Schema
}
