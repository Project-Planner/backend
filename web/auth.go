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
	username, pw, err := parseForm(w, r)
	if err != nil {
		return
	}

	errIncorrect := errors.New("username or password incorrect")

	l, err := db.GetLogin(username)
	if err == model.ErrNotFound {
		writeError(w, errIncorrect.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		writeError(w, "", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(l.Hash.Val), []byte(pw)); err != nil {
		writeError(w, errIncorrect.Error(), http.StatusUnauthorized)
		return
	}

	t, err := createToken(username)
	if err != nil {
		writeError(w, "", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	c := http.Cookie{
		Name:     authStr,
		Value:    t,
		Expires:  time.Now().Add(jwtDuration),
		HttpOnly: true,
		Path: conf.AuthedPathName,
	}
	http.SetCookie(w, &c)

	http.Redirect(w, r, "/html/mainPage.html", http.StatusSeeOther)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(authStr)
	if err != nil {
		writeError(w, "no authentication token (jwt) provided, please log in.\n"+err.Error(),
			http.StatusUnauthorized)
		return
	}

	deleteCookie(w, c)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userid, ok := r.Context().Value(userIDStr).(string)
	if !ok {
		writeError(w, "", http.StatusUnauthorized)
		return
	}

	err := db.DeleteUser(userid)
	if err != nil {
		log.Println(err)
		writeError(w, "", http.StatusInternalServerError)
		return
	}

	c, err := r.Cookie(authStr)
	if err != nil {
		writeError(w, "no authentication token (jwt) provided, please log in.\n"+err.Error(),
			http.StatusUnauthorized)
		return
	}

	deleteCookie(w, c)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	username, pw, err := parseForm(w, r)
	if err != nil {
		return
	}

	_, err = db.GetLogin(username)
	if err != nil && err != model.ErrNotFound {
		writeError(w, "", http.StatusInternalServerError)
		log.Println(err)
		return
	} else if err == nil {
		writeError(w, "username already exists", http.StatusConflict)
		return
	}

	if !legalName(username) {
		writeError(w, "illegal name", http.StatusUnprocessableEntity)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(pw), 14)
	if err != nil {
		writeError(w, "", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if err := db.AddUser(username, string(hashed)); err != nil {
		writeError(w, "", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	t, err := createToken(username)
	if err != nil {
		writeError(w, "", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	c := http.Cookie{
		Name:     authStr,
		Value:    t,
		Expires:  time.Now().Add(jwtDuration),
		HttpOnly: true,
		Path: conf.AuthedPathName,
	}
	http.SetCookie(w, &c)

	http.Redirect(w, r, "/html/mainPage.html", http.StatusSeeOther)
}

func parseForm(w http.ResponseWriter, r *http.Request) (username string, pw string, error error) {
	error = errors.New("parsing failed")

	// Parse HTML form from body
	if err := r.ParseForm(); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	uns, ok := r.Form["username"]
	if !ok || len(uns) != 1 {
		writeError(w, "username missing, html input must have name 'username'", http.StatusUnprocessableEntity)
		return
	}
	pws, ok := r.Form["password"]
	if !ok || len(pws) != 1 {
		writeError(w, "password missing, html input must have name 'password'", http.StatusUnprocessableEntity)
		return
	}

	return uns[0], pws[0], nil
}

// deleteCookie deletes the given cookie WITHOUT modifying the provided cookie, even though it is a pointer.
func deleteCookie(w http.ResponseWriter, c *http.Cookie) {
	delC := http.Cookie{
		Name:     c.Name,
		Value:    "",
		Path:     c.Path,
		Domain:   c.Domain,
		Expires:  time.Now().Add(-7 * 24 * time.Hour), // THIS DELETES THE COOKIE
		MaxAge:   -1,                                  // Tells browser to delete cookie NOW, but doesn't work with IE, hence 'Expires'
		HttpOnly: c.HttpOnly,
		SameSite: c.SameSite,
	}

	http.SetCookie(w, &delC)
}

