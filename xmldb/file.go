package xmldb

import (
	"log"
	"os"
)

//exists checks whether a file or directory at
//the given path does exist or not.
func exists(path string) bool {
	var _, err = os.Stat(path)
	return err == nil
}

//set overwrites the file behind @path with @content
func setFile(path, content string) {
	var f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err == nil {
		f.Truncate(0)
		f.Seek(0, 0)
		f.WriteString(content)
		f.Close()
	}
}

//create creates a file at the given path and
//initially fills it with the given content.
func create(path, content string) {
	var f, err = os.Create(path)
	if err != nil {
		log.Fatal(err)
	} else {
		f.WriteString(content)
		f.Close()
	}
}

//ensureDir makes sure that a directory at the given
//path exists by creating it if necessary.
func ensureDir(path string) {
	if !exists(path) {
		var err = os.Mkdir(path, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}
