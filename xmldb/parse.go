package xmldb

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

//parse takes in the content of the file behind @source and
//transforms it into the given struct @target.
func parse(source string, target interface{}) error {
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return xml.Unmarshal(content, target)
}
