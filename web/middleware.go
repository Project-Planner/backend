package web

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strconv"
	"time"
)

// auth authenticates and authorizes a user for accessing a requested resource
func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie(authStr)
			if err != nil {
				http.Error(w, "no authentication token (jwt) provided, please log in.\n"+err.Error(),
					http.StatusUnauthorized)
				return
			}

			t, err := parseTokenAndVerifySignature(c.Value)
			if err != nil {
				log.Println(err)
				http.Error(w, "untrusted signature, please log in again.", http.StatusUnauthorized)
				return
			}

			claims, ok := t.Claims.(jwt.MapClaims)
			if !ok || !t.Valid {
				http.Error(w, "token invalid, please log in again", http.StatusUnauthorized)
				return
			}

			uid, ok := claims[userIDStr]
			userID, cast := uid.(string)
			if !ok || !cast {
				http.Error(w, "user_id missing", http.StatusUnauthorized)
				return
			}
			exp, ok := claims[expiryStr]
			expiry, err := strconv.ParseInt(fmt.Sprintf("%.f", exp), 10, 64)
			if !ok || err != nil {
				http.Error(w, "expiry date missing", http.StatusUnauthorized)
				return
			}
			if time.Now().Unix() > expiry {
				deleteCookie(w, c)
				http.Error(w, "your session has expired, please log in again", http.StatusUnauthorized)
				return
			}

			// Sets the verified user context, this user is authenticated
			ctx := context.WithValue(r.Context(), userIDStr, userID)

			// executes the next function in the chain. Do not remove this.
			next.ServeHTTP(w, r.WithContext(ctx))
		})
}
