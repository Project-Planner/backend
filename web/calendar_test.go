package web

import (
	"context"
	"fmt"
	"github.com/Project-Planner/backend/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testOwner = "lambda"
var userView = "userView"
var userEdit = "userEdit"
var userNone = "userNone"

var defCalendar = model.Calendar{
	Name: struct {
		Text string `xml:",chardata"`
		Val  string `xml:"val,attr"`
	}{Val: testOwner},
	Owner: struct {
		Text string `xml:",chardata"`
		Val  string `xml:"val,attr"`
	}{Val: testOwner},
	ID: struct {
		Text string `xml:",chardata"`
		Val  string `xml:"val,attr"`
	}{Val: testOwner + "/" + testOwner},
	Permissions: struct {
		Text string `xml:",chardata"`
		View struct {
			Text string            `xml:",chardata"`
			User []model.Attribute `xml:"user"`
		} `xml:"view"`
		Edit struct {
			Text string            `xml:",chardata"`
			User []model.Attribute `xml:"user"`
		} `xml:"edit"`
	}(struct {
		Text string
		View struct {
			Text string            `xml:",chardata"`
			User []model.Attribute `xml:"user"`
		}
		Edit struct {
			Text string            `xml:",chardata"`
			User []model.Attribute `xml:"user"`
		}
	}{
		View: struct {
			Text string            `xml:",chardata"`
			User []model.Attribute `xml:"user"`
		}{User: []model.Attribute{
			{
				Val: userView,
			},
		}},
		Edit: struct {
			Text string            `xml:",chardata"`
			User []model.Attribute `xml:"user"`
		}{User: []model.Attribute{
			{
				Val: userEdit,
			},
		}}}),
}

func TestGetCalendarHandler(t *testing.T) {

	defaultCalendar := dbMock{data: map[string]struct {
		d interface{}
		e error
	}{
		"GetCalendar": {
			d: defCalendar,
		},
	}}

	tt := []struct {
		db        dbMock
		path      string
		urlParams string
		authed    string // user authed by middleware
		code      int
	}{
		// Kosher Case direct link
		{
			path:      "/c/" + testOwner + "/" + testOwner,
			urlParams: "?date=1.1.1970&time=15:14",
			authed:    testOwner,
			code:      http.StatusOK,
			db:        defaultCalendar,
		},
		// Kosher Case auto link
		{
			path:      "/c",
			urlParams: "?date=1.1.1970&time=15:14",
			authed:    testOwner,
			code:      http.StatusOK,
			db:        defaultCalendar,
		},
		// Unauthorized
		{
			path:      "/c",
			urlParams: "?date=1.1.1970&time=15:14",
			authed:    "",
			code:      http.StatusUnauthorized,
			db:        defaultCalendar,
		},
		// Forbidden
		{
			path:      "/c/" + testOwner + "/" + testOwner,
			urlParams: "?date=1.1.1970&time=15:14",
			authed:    userNone,
			code:      http.StatusForbidden,
			db:        defaultCalendar,
		},
		// user view
		{
			path:      "/c/" + testOwner + "/" + testOwner,
			urlParams: "?date=1.1.1970&time=15:14",
			authed:    userView,
			code:      http.StatusOK,
			db:        defaultCalendar,
		},
	}

	for _, tc := range tt {
		db = tc.db

		r, err := http.NewRequest("GET", tc.path+tc.urlParams, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(getCalendarHandler)

		p := strings.Split(tc.path, "/")
		muxVars := make(map[string]string)
		if len(p) == 3 {
			muxVars[calendarIDStr] = p[2]
		} else if len(p) == 4 {
			muxVars[userIDStr] = p[2]
			muxVars[calendarIDStr] = p[3]
		}

		ctx := context.WithValue(r.Context(), userIDStr, tc.authed)
		ctx = context.WithValue(ctx, 0, muxVars)

		if tc.authed == "" {
			handler.ServeHTTP(rr, r)
		} else {
			handler.ServeHTTP(rr, r.WithContext(ctx))
		}

		if rr.Code != tc.code {
			t.Fatalf("wrong status code: got: %d want: %d \n%s\n%v", rr.Code, tc.code, rr.Body.String(), tc)
		}

		if tc.code == http.StatusOK &&
			!strings.Contains(rr.Body.String(),
				fmt.Sprintf(`href="%s"`, conf.AuthedPathName+"/calendar.xsl"+tc.urlParams)) {
			t.Fatalf("url params not added successfully. Got: %s \nIn test case %v", rr.Body.String(), tc)
		}
	}
}

func TestLegalName(t *testing.T) {
	tt := []struct {
		name string
		want bool
	}{
		// Kosher case
		{
			name: "abcDEF09-_",
			want: true,
		},
		// Critical illegal character
		{
			name: "a/bc/",
			want: false,
		},
		// Empty name
		{
			name: "",
			want: false,
		},
		// space only
		{
			name: " ",
			want: false,
		},
	}

	for _, tc := range tt {
		if got := legalName(tc.name); tc.want != got {
			t.Errorf("got: %v\ntc: %v", got, tc)
		}
	}
}

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
