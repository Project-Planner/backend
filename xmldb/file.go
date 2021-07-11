package xmldb

import (
	"io/ioutil"
	"os"
)

//exists checks whether a file or directory at
//the given path does exist or not.
func exists(path string) bool {
	var _, err = os.Stat(path)
	return err == nil
}

//write overwrites the file behind @path with @content
func write(path, content string) error {
	return ioutil.WriteFile(path, []byte(content), 0666)
}

//ensureDir makes sure that a directory at the given
//path exists by creating it if necessary.
func ensureDir(path string) error {
	if !exists(path) {
		return os.Mkdir(path, 0755)
	}
	return nil
}
