package config

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

// Load returns the config and a possible error.
func Load() (Config, error) {
	path := "./config.yaml"
	return load(path)
}

func load(path string) (Config, error) {
	var c Config

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(b, &c)
	return c, err
}
