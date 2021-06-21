package config

import (
	"github.com/Project-Planner/backend/web"
	"github.com/Project-Planner/backend/xmldb"
)

// Config of the project, embedding all sub-configs
type Config struct {
	// here go the config fields and embedded configs of other packages
	web.ServerConfig `yaml:",inline"`
	xmldb.DBConfig   `yaml:",inline"`
}
