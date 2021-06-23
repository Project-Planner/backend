package web

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	theConfig := ServerConfig{JWTSecret: "some%secret"}

	conf = theConfig

	createExpiredToken := func(username string) string {
		exp := time.Now().Add(-jwtDuration).Unix()

		c := jwt.MapClaims{}
		c[authorizedStr] = true
		c[tokenIDStr] = uuid.New()
		c[userIDStr] = username
		c[expiryStr] = exp

		t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
		tStr, err := t.SignedString([]byte(conf.JWTSecret))
		if err != nil {
			panic(err)
		}
		return tStr
	}

	un := "someusername"
	okCookie, err := createCookie(un)
	if err != nil {
		t.Fatal(err)
	}

	tt := []struct {
		un               string
		c                http.Cookie
		jwtMissing       bool
		invalidSignature bool
		code             int
	}{
		// kosher case
		{
			un:   un,
			c:    okCookie,
			code: http.StatusOK,
		},
		// cookie missing
		{
			jwtMissing: true,
			c:          okCookie,
			code:       http.StatusUnauthorized,
		},
		// untrusted signature
		{
			invalidSignature: true,
			c:                okCookie,
			code:             http.StatusUnauthorized,
		},
		// invalid (empty) cookie
		{
			c: http.Cookie{
				Name:    okCookie.Name,
				Value:   "",
				Expires: okCookie.Expires,
			},
			code: http.StatusUnauthorized,
		},
		// expired token
		{
			un: un,
			c: http.Cookie{
				Name:  authStr,
				Value: createExpiredToken(un),
			},
			code: http.StatusUnauthorized,
		},
	}

	for _, tc := range tt {
		conf = theConfig

		r, err := http.NewRequest("POST", "/authorize", nil)
		if err != nil {
			t.Fatal(err)
		}

		// special test cases
		if !tc.jwtMissing {
			r.AddCookie(&tc.c)
		}
		if tc.invalidSignature {
			conf = ServerConfig{
				JWTSecret: "otherSecret",
			}
		}

		// tests whether the context
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			n, ok := r.Context().Value(userIDStr).(string)
			if !ok {
				t.Errorf("%s not in request context: got %q", userIDStr, n)
			}
			if n != tc.un {
				t.Errorf("username incorrect, want: %s got: %s", tc.un, n)
			}

			w.WriteHeader(http.StatusOK)
		})

		rr := httptest.NewRecorder()
		handler := auth(testHandler)

		handler.ServeHTTP(rr, r)

		if rr.Code != tc.code {
			t.Fatalf("wrong status code: got: %d want: %d \n%s\n%v", rr.Code, tc.code, rr.Body.String(), tc)
		}
	}

}

func createCookie(username string) (http.Cookie, error) {
	t, err := createToken(username)
	if err != nil {
		return http.Cookie{}, err
	}

	c := http.Cookie{
		Name:     authStr,
		Value:    t,
		Expires:  time.Now().Add(jwtDuration),
		HttpOnly: true,
	}

	return c, nil
}
