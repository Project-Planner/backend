package model

import (
	"encoding/json"
	"encoding/xml"
	"testing"
)

const docXML = `
<auth>
    <login>
        <name val="Peter" />
        <hash val="trololol" />
    </login>
    <login>
        <name val="Günther" />
        <hash val="trololol" />
    </login>
</auth>`

func TestParseAuth(t *testing.T) {
	var a Auth

	err := xml.Unmarshal([]byte(docXML), &a)
	if err != nil {
		t.Fatal(err)
	}

	if len(a.Logins) != 2 || a.Logins[1].Hash.Val != "trololol" {
		s, _ := json.MarshalIndent(a, "", "\t")
		t.Fatal("want: " + docXML + "\ngot: " + string(s))
	}
}
