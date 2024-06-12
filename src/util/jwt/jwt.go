package jwt

import (
	"strings"
	"time"

	"github.com/Artexus/api-widyabhuvana/src/constant"
	"github.com/Artexus/api-widyabhuvana/src/util/aes"
	jwt "github.com/form3tech-oss/jwt-go"
)

type TokenResponse struct {
	EncID     string
	Username  string
	Email     string
	ExpiredIn int64
}

type GenerateTokenPayload struct {
	EncUserID string
	Username  string
	Email     string
}

func (tr TokenResponse) IsExpired() bool {
	return !time.Unix(tr.ExpiredIn, 0).After(time.Now())
}

func GenerateToken(payload GenerateTokenPayload) (tokenString string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	expired := time.Now().Add(constant.JWTInterval).Unix()

	claims := make(jwt.MapClaims)
	claims["user_id"] = payload.EncUserID
	claims["email"] = payload.Email
	claims["username"] = payload.Username
	claims["exp"] = expired
	token.Claims = claims

	tokenString, err = token.SignedString([]byte(constant.JWTSignedKey))
	return
}

func ExtractIDToken(token string) (id string, err error) {
	claims := jwt.MapClaims{}

	if strings.Contains(token, "Bearer ") {
		token = strings.Split(token, "Bearer ")[1]
	}
	_, err = jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(constant.JWTSignedKey), nil
	})
	if err != nil {
		err = constant.ErrInvalid
		return
	}

	id, err = aes.DecryptID(claims["user_id"].(string))
	return
}

func ExtractToken(token string) (resp TokenResponse, err error) {
	claims := jwt.MapClaims{}

	if strings.Contains(token, "Bearer ") {
		token = strings.Split(token, "Bearer ")[1]
	}
	_, err = jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(constant.JWTSignedKey), nil
	})

	if err != nil && strings.Contains(err.Error(), "expired") {
		err = constant.ErrExpired
		return
	} else if err != nil {
		err = constant.ErrExpired
		return
	}

	resp.EncID = claims["user_id"].(string)
	resp.Email = claims["email"].(string)
	resp.Username = claims["username"].(string)
	resp.ExpiredIn = int64(claims["exp"].(float64))
	return
}
