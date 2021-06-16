package web

import "net/http"

// auth authenticates and authorizes a user for accessing a requested resource
func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			panic("IMPLEMENT ME")

			// executes the next function in the chain. Do not remove this.
			next.ServeHTTP(w, r)
		})
}
