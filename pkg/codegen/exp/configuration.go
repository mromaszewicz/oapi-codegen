package exp

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	ImportMapping map[string]ImportSpec  `yaml:"import-mapping,omitempty"`
	TypeOverrides OpenAPIToGoTypeMapping `yaml:"type-mapping,omitempty"`

	// typeMapping is the composite mapping of the default and the overrides
	typeMapping OpenAPIToGoTypeMapping `yaml:"-"`
}

type ImportSpec struct {
	Alias   string `yaml:"alias,omitempty"`
	Package string `yaml:"package"`
}

func LoadConfiguration(buf []byte) (Configuration, error) {
	var cfg Configuration
	err := yaml.Unmarshal(buf, &cfg)
	if err != nil {
		return Configuration{}, fmt.Errorf("loading configuration: %w", err)
	}

	// Now, merge any overrides on top of default values.
	cfg.typeMapping = OverrideTypeMapping(cfg.typeMapping, GetDefaultTypeMapping())

	return cfg, nil
}
