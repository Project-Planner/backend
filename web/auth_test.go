package web

import (
	"errors"
	"github.com/Project-Planner/backend/model"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	conf = ServerConfig{
		JWTSecret: "some%secret",
	}

	hash := func(pw string) string {
		hashedPw, err := bcrypt.GenerateFromPassword([]byte(pw), 12)
		if err != nil {
			t.Fatal(err)
		}
		return string(hashedPw)
	}

	userFoundDB := dbMock{data: map[string]struct {
		d interface{}
		e error
	}{
		"GetLogin": {
			d: model.Login{
				Name: struct {
					Text string `xml:",chardata"`
					Val  string `xml:"val,attr"`
				}{Val: "someusername"},
				Hash: struct {
					Text string `xml:",chardata"`
					Val  string `xml:"val,attr"`
				}{Val: hash("supersafepassword%&$")},
			},
		},
	},
	}

	type form struct {
		pwField string
		unField string
	}

	okForm := form{
		pwField: "password",
		unField: "username",
	}

	tt := []struct {
		pw     string
		un     string
		db     dbMock
		badReq bool
		form   form
		code   int
	}{
		// Everything kosher test case
		{
			pw:   "supersafepassword%&$",
			un:   "someusername",
			code: http.StatusOK,
			db:   userFoundDB,
			form: okForm,
		},
		// Wrong password test case
		{
			pw:   "supersafeWRONGpassword%&$",
			un:   "someusername",
			code: http.StatusUnauthorized,
			db:   userFoundDB,
			form: okForm,
		},
		// Wrong username test case
		{
			pw:   "supersafepassword%&$",
			un:   "someWRONGusername",
			code: http.StatusUnauthorized,
			db:   userFoundDB,
			form: okForm,
		},
		// Wrong password DB side
		{
			pw:   "supersafepassword%&$",
			un:   "someusername",
			code: http.StatusUnauthorized,
			form: okForm,
			db: dbMock{data: map[string]struct {
				d interface{}
				e error
			}{
				"GetLogin": {
					d: model.Login{
						Name: struct {
							Text string `xml:",chardata"`
							Val  string `xml:"val,attr"`
						}{Val: "someusername"},
						Hash: struct {
							Text string `xml:",chardata"`
							Val  string `xml:"val,attr"`
						}{Val: hash("supersafeWRONGpassword%&$")},
					},
				},
			},
			},
		},
		// User not found
		{
			pw:   "supersafepassword%&$",
			un:   "someusername",
			code: http.StatusUnauthorized,
			form: okForm,
			db: dbMock{data: map[string]struct {
				d interface{}
				e error
			}{
				"GetLogin": {
					d: model.Login{},
					e: model.ErrNotFound,
				},
			}},
		},
		// DB error other than user not found
		{
			pw:   "supersafepassword%&$",
			un:   "someusername",
			code: http.StatusInternalServerError,
			form: okForm,
			db: dbMock{data: map[string]struct {
				d interface{}
				e error
			}{
				"GetLogin": {
					d: model.Login{},
					e: errors.New("whoops, sth went wrong in the database"),
				},
			}},
		},
		// Not sending form data
		{
			badReq: true,
			form:   okForm,
			code:   http.StatusBadRequest,
		},
		// wrong spelled password form field
		{
			code: http.StatusUnprocessableEntity,
			db:   userFoundDB,
			form: form{
				pwField: "paZZword",
				unField: "username",
			},
		},
		// wrong spelled username form field
		{
			code: http.StatusUnprocessableEntity,
			db:   userFoundDB,
			form: form{
				pwField: "password",
				unField: "jusername",
			},
		},
	}

	for _, tc := range tt {
		db = tc.db

		data := url.Values{}
		data.Set(tc.form.pwField, tc.pw)
		data.Set(tc.form.unField, tc.un)

		var rdr io.Reader
		rdr = strings.NewReader(data.Encode())
		if tc.badReq {
			rdr = nil
		}

		r, err := http.NewRequest("POST", "/authorize", rdr)
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(loginHandler)

		handler.ServeHTTP(rr, r)

		if rr.Code != tc.code {
			t.Fatalf("wrong status code: got: %d want: %d \n%s\n%v", rr.Code, tc.code, rr.Body.String(), tc)
		}

		// do not execute the rest of the test, if it's not a test case looking for success
		if tc.code != http.StatusOK {
			continue
		}

		cRaw := rr.Header().Get("Set-Cookie")
		h := http.Header{}
		h.Add("Cookie", cRaw)
		c, _ := (&http.Request{Header: h}).Cookie(authStr)

		token, err := parseTokenAndVerifySignature(c.Value)
		if err != nil {
			t.Fatal(err)
		}

		if tc.un != token.Claims.(jwt.MapClaims)[userIDStr].(string) {
			t.Fatal("username wrong cookie not parsed correctly")
		}
	}
}

func TestRegisterHandler(t *testing.T) {
	un := "nickname"
	pw := "mypw"

	okDB := dbMock{data: map[string]struct {
		d interface{}
		e error
	}{
		"GetLogin": {
			d: model.Login{
				Name: struct {
					Text string `xml:",chardata"`
					Val  string `xml:"val,attr"`
				}{Val: un},
			},
		},
	}}

	tt := []struct {
		un   string
		pw   string
		db   dbMock
		code int
	}{
		// Kosher case
		{
			un:   "otherUsername",
			pw:   pw,
			db:   okDB,
			code: http.StatusCreated,
		},
		// Username already exists
		{
			un:   un,
			pw:   pw,
			db:   okDB,
			code: http.StatusConflict,
		},
		// Database issue while getting user
		{
			un: un,
			pw: pw,
			db: dbMock{data: map[string]struct {
				d interface{}
				e error
			}{
				"GetLogin": {
					d: model.Login{},
					e: errors.New("whupsi, some error"),
				},
			}},
			code: http.StatusInternalServerError,
		},
		// Database issue while adding user
		{
			un: un,
			pw: pw,
			db: dbMock{data: map[string]struct {
				d interface{}
				e error
			}{
				"GetLogin": {
					e: model.ErrNotFound,
				},
				"AddUser": {
					e: errors.New("whupsi, some error"),
				},
			}},
			code: http.StatusInternalServerError,
		},
	}

	for _, tc := range tt {
		db = tc.db

		data := url.Values{}
		data.Set("password", tc.pw)
		data.Set("username", tc.un)

		r, err := http.NewRequest("POST", "/authorize", strings.NewReader(data.Encode()))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(registerHandler)

		handler.ServeHTTP(rr, r)

		if rr.Code != tc.code {
			t.Fatalf("wrong status code: got: %d want: %d \n%s\n%v", rr.Code, tc.code, rr.Body.String(), tc)
		}
	}
}

type dbMock struct {
	data map[string]struct {
		d interface{}
		e error
	}
}

func (d dbMock) GetCalendar(calendarid string) (model.Calendar, error) {
	e := d.data["GetCalendar"].e
	if e != nil {
		return model.Calendar{}, e
	}

	un := d.data["GetCalendar"].d.(model.Calendar)
	return un, e
}

func (d dbMock) AddUser(userid, hashedPW string) error {
	return d.data["AddUser"].e
}

func (d dbMock) GetLogin(userid string) (model.Login, error) {
	e := d.data["GetLogin"].e
	if e != nil {
		return model.Login{}, e
	}

	un := d.data["GetLogin"].d.(model.Login)
	if userid != un.Name.Val {
		return model.Login{}, model.ErrNotFound
	}
	return un, e
}
