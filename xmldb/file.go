package xmldb

import (
	"io/ioutil"
	"log"
	"os"
)

//exists checks whether a file or directory at
//the given path does exist or not.
func exists(path string) bool {
	var _, err = os.Stat(path)
	return err == nil
}

//write overwrites the file behind @path with @content
func write(path, content string) {
	if err := ioutil.WriteFile(path, []byte(content), 0666); err != nil {
		log.Fatal(err)
	}
}

//ensureDir makes sure that a directory at the given
//path exists by creating it if necessary.
func ensureDir(path string) {
	if !exists(path) {
		if err := os.Mkdir(path, 0755); err != nil {
			log.Fatal(err)
		}
	}
}
