package web

import (
	"github.com/Project-Planner/backend/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"strconv"
	"time"
)

// createToken takes the userid and generates a JWT authentication token for it. The token is returned as model.JWT
// and signed string ready for transmission to the client's browser, and an error, if one occurs
func createToken(userid string) (model.JWT, string, error) {
	authorized := true
	id, err := uuid.NewRandom()
	if err != nil {
		return model.JWT{}, "", err
	}
	tokenID := id.String()
	exp := time.Now().Add(time.Hour * 365 * 24).Unix()

	c := jwt.MapClaims{}
	c["authorized"] = authorized
	c["token_id"] = tokenID
	c["user_id"] = userid
	c["expiry"] = exp

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenStr, err := t.SignedString(conf.JWTSecret)

	return model.JWT{
		Text: "",
		Authorized: struct {
			Text string `xml:",chardata"`
			Val  string `xml:"val,attr"`
		}{Text: "", Val: strconv.FormatBool(authorized)},
		TokenID: struct {
			Text string `xml:",chardata"`
			Val  string `xml:"val,attr"`
		}{Text: "", Val: tokenID},
		UserID: struct {
			Text string `xml:",chardata"`
			Val  string `xml:"val,attr"`
		}{Text: "", Val: userid},
		Expiry: struct {
			Text string `xml:",chardata"`
			Val  string `xml:"val,attr"`
		}{Text: "", Val: strconv.FormatInt(exp, 10)},
	}, tokenStr, err
}
