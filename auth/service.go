package auth

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(userId int) (string, error)
	ValidateToken(encodeToken string) (*jwt.Token, error)
}

type jwtService struct {
}

func NewSevice() *jwtService {
	return &jwtService{}
}

var SECREAT_KEY = []byte("bw4_g0l4n9_secreat_k3y")

func (j *jwtService) GenerateToken(userId int) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = userId

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(SECREAT_KEY)
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil

}
func (j *jwtService) ValidateToken(encodeToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodeToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Invalid token")
		}
		return []byte(SECREAT_KEY), nil
	})

	if err != nil {
		return token, err
	}
	return token, nil
}
