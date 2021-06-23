package web

// ServerConfig of the web web
type ServerConfig struct {
	// Here go the config fields that are relevant to the webserver
	// Port for listening to http requests
	Port int `yaml:"port"`
	// StaticDir where the static web pages reside, that may be viewed by anyone
	StaticDir string `yaml:"static_dir"`
	// HTMLDir for other html files
	HTMLDir string `yaml:"html_dir"`
	// AuthedPathName - Path prefix of routes requiring authentication
	AuthedPathName string `yaml:"authed_path_name"`
	// JWTSecret contains the private, server-sided secret to sign JWTs
	JWTSecret string `yaml:"jwt_secret"`
}
