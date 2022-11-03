package exp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadDefaultTypeMapping(t *testing.T) {
	dm := GetDefaultTypeMapping()
	assert.NotEmpty(t, dm)
}

func TestOverrideTypeMapping(t *testing.T) {
	defaultMapping := GetDefaultTypeMapping()

	// override the int32 default to be "banana"
	overrideMapping := OpenAPIToGoTypeMapping{
		"integer": GoTypesForFormats{
			Formats: map[string]GoTypeWithImport{
				"int32": GoTypeWithImport{
					GoType: "banana",
					Import: "encoding/banana",
				},
			},
		},
		// We'll add new string format named bob
		"string": GoTypesForFormats{
			Formats: map[string]GoTypeWithImport{
				"bob": GoTypeWithImport{
					GoType: "bob",
				},
			},
		},
		// Replace the default type for numbers to something new.
		"number": GoTypesForFormats{
			Default: "float128",
		},
	}

	// First, a nil override makes no changes.
	result := OverrideTypeMapping(defaultMapping, nil)
	assert.EqualValues(t, defaultMapping, result)

	// Applying a populated mapping over an empty one should produce equality
	// as well.
	result = OverrideTypeMapping(make(OpenAPIToGoTypeMapping), defaultMapping)
	assert.EqualValues(t, defaultMapping, result)

	// Now, do our overrides and make sure everything is as expected.
	result = OverrideTypeMapping(defaultMapping, overrideMapping)
	assert.Equal(t, "banana", result["integer"].Formats["int32"].GoType)
	assert.Equal(t, "encoding/banana", result["integer"].Formats["int32"].Import)
	assert.Equal(t, "bob", result["string"].Formats["bob"].GoType)
	assert.Equal(t, "float128", result["number"].Default)
}
