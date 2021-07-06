package web

import (
	"strings"
	"testing"
)

func TestVarXLS_String(t *testing.T) {
	want := `<xsl:variable name="weekDate" select="'1.1.1970'"/>` + "\n"
	got := varXLS{name: "weekDate", value: "1.1.1970"}.String()
	if want != got {
		t.Error("want: " + want + "\ngot: " + got)
	}
}

func TestVarsIntoXLS(t *testing.T) {
	want := `<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet version="1.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
<xsl:variable name="weekDate" select="'1.1.1970'"/>
<xsl:variable name="displayMode" select="'calendar'"/>
<xsl:variable name="calendarMode" select="'week'"/>
  
  <xsl:template match="/">`

	got := varsIntoXSL(xlsTruncated,
		varXLS{"weekDate", "1.1.1970"},
		varXLS{"displayMode", "calendar"},
		varXLS{"calendarMode", "week"},
	)

	w := strings.ReplaceAll(want, "\r", "")
	w = strings.ReplaceAll(w, " ", "")
	w = strings.ReplaceAll(w, "\t", "")
	w = strings.ReplaceAll(w, "\n", "")

	g := strings.ReplaceAll(got, "\r", "")
	g = strings.ReplaceAll(g, " ", "")
	g = strings.ReplaceAll(g, "\t", "")
	g = strings.ReplaceAll(g, "\n", "")

	if w != g {
		t.Error("want: " + w + "\ngot: " + g)
	}
}

const xlsTruncated = `<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet version="1.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  
  <xsl:template match="/">`
