package model

import (
	"encoding/json"
	"encoding/xml"
	"testing"
)

const docUserXML = `
<user>
    <name val="Günther Pascal" /> <!-- Has to be limited to (ASCII-)letters, numbers and +, -, _ (FE and BE)-->
    <items>
        <calendar href="file1" />
        <calendar href="file2" />
        <calendar href="file3" />
    </items>
</user>`

func TestParseUser(t *testing.T) {
	var u User

	err := xml.Unmarshal([]byte(docUserXML), &u)
	if err != nil {
		t.Fatal(err)
	}

	if u.Name.Val != "Günther Pascal" || len(u.Items.Calendars) != 3 || u.Items.Calendars[1].Link != "file2" {
		s, _ := json.MarshalIndent(u, "", "\t")
		t.Fatal("want: " + docUserXML + "\ngot: " + string(s))
	}
}
