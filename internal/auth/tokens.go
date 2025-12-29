package auth

import (
	"fmt"
	"time"

	"github.com/NurulloMahmud/habits/internal/user"
	"github.com/golang-jwt/jwt/v5"
)

func (s *JWTService) GenerateAccessToken(user user.User) (string, error) {
	claims := jwt.MapClaims{
		"email":      user.Email,
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"role":       user.UserRole,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *JWTService) VerifyToken(tokenString string) (*user.UserTokenClaim, error) {
	claims := &user.UserTokenClaim{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
