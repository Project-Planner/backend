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
        <jwts>
            <jwt>
                <authorized val="true" />
                <token_id val="uuid" />
                <user_id val="userid" />
                <expiry val="val" />
            </jwt>
            <jwt>
                <authorized val="true" />
                <token_id val="uuid" />
                <user_id val="userid" />
                <expiry val="val" />
            </jwt>
        </jwts>
    </login>
    <login>
        <name val="Günther" />
        <hash val="trololol" />
        <jwts>
            <jwt>
                <authorized val="true" />
                <token_id val="uuid" />
                <user_id val="userid" />
                <expiry val="val" />
            </jwt>
            <jwt>
                <authorized val="true" />
                <token_id val="uuid" />
                <user_id val="userid" />
                <expiry val="val" />
            </jwt>
        </jwts>
    </login>
</auth>`

func TestParseAuth(t *testing.T) {
	var a Auth

	err := xml.Unmarshal([]byte(docXML), &a)
	if err != nil {
		t.Fatal(err)
	}

	if len(a.Logins) != 2 || len(a.Logins[0].JWTs.JWTs) != 2 || a.Logins[0].JWTs.JWTs[0].Authorized.Val != "true" {
		s, _ := json.MarshalIndent(a, "", "\t")
		t.Fatal("want: " + docXML + "\ngot: " + string(s))
	}
}
