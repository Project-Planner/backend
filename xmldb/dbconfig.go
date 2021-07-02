package xmldb

// DBConfig of the XML database
type DBConfig struct {
	// DBRootDir - where the database resides on the system
	RootDir string `yaml:"root_dir"`
	// AuthDir - where the authentication file is stored
	AuthDir string `yaml:"auth_dir"`
	// UserDir - where the user files are stored
	// AuthSubpath - where the authentication file actually is stored relative
	// to AuthDir
	AuthSubpath string `yaml:"auth_subpath"`
	UserDir     string `yaml:"user_dir"`
	// CalendarDir - where the calendar files are stored
	CalendarDir string `yaml:"calendar_dir"`
	// CacheSize - how many elements (e.g. users) will be cached simultaneously;
	// amount will be limited, so that RAM doesn't get flooded with elements.
	CacheSize int `yaml:"cache_size"`
}
