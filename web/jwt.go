package web

import (
	"errors"
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
	tkID := id.String()
	exp := time.Now().Add(time.Hour * 365 * 24).Unix()

	c := jwt.MapClaims{}
	c[authorizedStr] = true
	c[tokenIDStr] = tkID
	c[userIDStr] = userid
	c[expiryStr] = exp

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString(conf.JWTSecret)
}

func parseTokenAndVerifySignature(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("jwt: signing method not correct")
		}

		return conf.JWTSecret, nil
	})
	return token, err
}
