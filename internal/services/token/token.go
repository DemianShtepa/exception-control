package token

import (
	"github.com/DemianShtepa/exception-control/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Generator struct {
	secret string
}

func NewGenerator(secret string) *Generator {
	return &Generator{secret: secret}
}

func (g *Generator) GenerateToken(user domain.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(g.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
