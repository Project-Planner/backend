package web

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

// createToken takes the userid and generates a JWT authentication token for it. The token is returned as model.JWT
// and signed string ready for transmission to the client's browser, and an error, if one occurs
func createToken(userid string) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	tokenID := id.String()
	exp := time.Now().Add(time.Hour * 365 * 24).Unix()

	c := jwt.MapClaims{}
	c["authorized"] = true
	c["token_id"] = tokenID
	c["user_id"] = userid
	c["expiry"] = exp

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString(conf.JWTSecret)
}
