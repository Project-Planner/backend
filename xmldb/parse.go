package xmldb

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

//parse takes in the content of the file behind @source and
//transforms it into the given struct @target.
func parse(source string, target interface{}) {
	file, _ := os.Open(source)
	content, _ := ioutil.ReadAll(file)
	xml.Unmarshal(content, target)
}
