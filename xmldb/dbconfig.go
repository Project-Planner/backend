package xmldb

// DBConfig of the XML database
type DBConfig struct {
	// DBRootDir - where the database resides on the system
	DBRootDir string `yaml:"db_root_dir"`
}
