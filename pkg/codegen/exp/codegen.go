package exp

import (
	"io"
	"log"

	"github.com/getkin/kin-openapi/openapi3"
)

func Generate(spec *openapi3.T, cfg Configuration, output io.Writer, logger *log.Logger) error {
	tree := BuildSchemaTree(spec)
	schemasByPath := GenerateSchemaInfosFromTree(tree)

	sortedSchemaKeys := SortedMapKeys(schemasByPath)
	logger.Println("Found paths:")
	for _, k := range sortedSchemaKeys {
		logger.Println("  ", k)
	}

	return nil
}
