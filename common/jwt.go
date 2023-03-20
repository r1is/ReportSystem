package common

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/r1is/ReportSystem/model"
	"time"
)

var jwtKey = []byte("iCODE@2023")

type Claims struct {
	UserId uint
	Phone  string
	jwt.StandardClaims
}

func ReleaseToken(user model.User) (string, error) {
	expTime := time.Now().Add(1 * 24 * time.Hour)
	claims := &Claims{
		UserId: user.ID,
		Phone:  user.Telephone,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "ReportSystem",
			Subject:   "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}
