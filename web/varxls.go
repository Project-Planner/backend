package web

import (
	"fmt"
	"net/url"
)

type varXLS struct {
	name  string
	value string
}

func (v varXLS) String() string {
	return fmt.Sprintf(`<xsl:variable name="%s" select="'%s'"/>`, v.name, v.value) + "\n"
}

func allFromURL(params url.Values) []varXLS {
	vars := make([]varXLS, len(params))
	i := 0
	for k, v := range params {
		vars[i] = varXLS{
			name:  k,
			value: v[0],
		}
		i++
	}

	return vars
}
