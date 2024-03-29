package xmldb

//DBConfig of the XML database
type DBConfig struct {
	//RootDir - here all database files reside
	DBDir string `yaml:"db_dir"`

	//AuthRelDir - relative path (to root dir) where authentication files are stored.
	AuthRelDir string

	//AuthDir - absolute path of where authentication files are stored.
	AuthDir string

	//UserRelDir - relative path (to root dir) where user files are stored.
	UserRelDir string

	//UserDir
	UserDir string

	//CalendarRelDir - relative path (to root dir) where calendar files are stored.
	CalendarRelDir string

	//CalendarDir
	CalendarDir string

	// CacheSize - how many bytes (e.g. for users) will be cached simultaneously;
	//			   cache can be used to prevent RAM getting flooded with elements.
	CacheSize int
}
