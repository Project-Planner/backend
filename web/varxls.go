package web

import "fmt"

type varXLS struct {
	name  string
	value string
}

func (v varXLS) String() string {
	return fmt.Sprintf(`<xsl:variable name="%s" select="'%s'"/>`, v.name, v.value) + "\n"
}
