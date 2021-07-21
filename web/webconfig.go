package web

// ServerConfig of the web web
type ServerConfig struct {
	// Here go the config fields that are relevant to the webserver
	// Port for listening to http requests
	Port int `yaml:"port"`
	// FrontendDir where the static web pages reside, that may be viewed by anyone
	FrontendDir string `yaml:"frontend_dir"`
	// AuthedPathName - Path prefix of routes requiring authentication
	AuthedPathName string `yaml:"authed_path_name"`
	// JWTSecret contains the private, server-sided secret to sign JWTs
	JWTSecret string `yaml:"jwt_secret"`
}
