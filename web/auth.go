package web

import (
	"errors"
	"github.com/Project-Planner/backend/model"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse HTML form from body
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uns, ok := r.Form["username"]
	if !ok || len(uns) != 1{
		http.Error(w, "username missing, html input must have name 'username'", http.StatusUnprocessableEntity)
	}
	pws, ok := r.Form["password"]
	if !ok || len(pws) != 1{
		http.Error(w, "password missing, html input must have name 'password'", http.StatusUnprocessableEntity)
	}

	username := uns[0]
	pw := pws[0]

	errIncorrect := errors.New("username or password incorrect")

	l, err := db.GetLogin(username)
	if err == model.ErrNotFound {
		http.Error(w, errIncorrect.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(l.Hash.Val), []byte(pw)); err != nil {
		http.Error(w, errIncorrect.Error(), http.StatusUnauthorized)
		return
	}

	t, err := createToken(username)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	c := http.Cookie{
		Name:       authStr,
		Value:      t,
		Expires:    time.Now().Add(time.Hour * 365 * 24),
		HttpOnly:   true,
	}
	http.SetCookie(w, &c)

	w.WriteHeader(http.StatusOK) // possibly redirect to another page later
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(authStr)
	if err != nil {
		http.Error(w, "no authentication token (jwt) provided, please log in.\n" + err.Error(),
			http.StatusUnauthorized)
		return
	}

	deleteCookie(w, c)

	w.WriteHeader(http.StatusOK) // possibly redirect to another page later
}

// deleteCookie deletes the given cookie WITHOUT modifying the provided cookie, even though it is a pointer.
func deleteCookie(w http.ResponseWriter, c *http.Cookie) {
	delC := http.Cookie{
		Name:       c.Name,
		Value:      "",
		Path:       c.Path,
		Domain:     c.Domain,
		Expires:    time.Now().Add(-7 * 24 * time.Hour), // THIS DELETES THE COOKIE
		MaxAge:     -1, // Tells browser to delete cookie NOW, but doesn't work with IE, hence 'Expires'
		HttpOnly:   c.HttpOnly,
		SameSite:   c.SameSite,
	}

	http.SetCookie(w, &delC)
}