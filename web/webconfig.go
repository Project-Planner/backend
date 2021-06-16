package web

// ServerConfig of the web web
type ServerConfig struct {
	// Here go the config fields that are relevant to the webserver
	Port           int
	StaticDir      string
	AuthedPathName string
}
