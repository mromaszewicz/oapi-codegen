package exp

import (
	_ "embed"

	"gopkg.in/yaml.v2"
)

// OpenAPIToGoTypeMapping is a map of maps which converts OpenAPI
// type and format specifiers into a Go typename. See the default
// type mapping below for syntax.
type OpenAPIToGoTypeMapping map[string]GoTypesForFormats

// GoTypeWithImport represents a go type, eg, time.Time, with the corresponding
// import, if any.
type GoTypeWithImport struct {
	GoType   string `yaml:"type"`
	Nullable bool   `yaml:"nullable,omitempty"`
	Import   string `yaml:"import,omitempty"`
}

type GoTypesForFormats struct {
	Formats map[string]GoTypeWithImport `yaml:"formats"`
	Default string                      `yaml:"default"`
}

//go:embed defaults/typemapping.yaml
var defaultTypeMapping []byte

// GetDefaultTypeMapping loads the type mapping defaults, or panics if it fails.
// We unit test this function, so it should never panic given that we also
// ship the defaults.
func GetDefaultTypeMapping() OpenAPIToGoTypeMapping {
	var result OpenAPIToGoTypeMapping
	err := yaml.Unmarshal(defaultTypeMapping, &result)
	if err != nil {
		// This should never happen, as we ensure via unit test that the default
		// mapping loads correctly.
		panic(err)
	}
	return result
}

// OverrideTypeMapping returns a new type mapping, where values in 'base' are
// replaced or supplemented with values from 'override'
func OverrideTypeMapping(base, override OpenAPIToGoTypeMapping) OpenAPIToGoTypeMapping {
	result := make(OpenAPIToGoTypeMapping)

	mappings := []OpenAPIToGoTypeMapping{base, override}

	for _, mapping := range mappings {
		for typeName, formats := range mapping {
			newFormats := result[typeName]
			if formats.Default != "" {
				newFormats.Default = formats.Default
			}
			for formatName, formatSpec := range formats.Formats {
				if newFormats.Formats == nil {
					newFormats.Formats = make(map[string]GoTypeWithImport)
				}
				newFormats.Formats[formatName] = formatSpec
			}
			result[typeName] = newFormats
		}
	}
	return result
}
