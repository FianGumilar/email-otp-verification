package helpers

import (
	"strings"
	"time"

	"github.com/FianGumilar/email-otp-verification/config"
	"github.com/FianGumilar/email-otp-verification/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type JwtClaims struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJwtToken(user models.User) (string, error) {
	claims := JwtClaims{
		user.ID,
		user.Email,
		user.FullName,
		jwt.StandardClaims{
			Id:        "12345678",
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := rawToken.SignedString([]byte(config.GetEnv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParseToken(ctx echo.Context) (*JwtClaims, error) {
	token := ctx.Request().Header.Get("Authorization")
	token = strings.ReplaceAll(token, "Bearer ", "")

	// Parse Token Jwt
	claims := &JwtClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv("JWT_KEY")), nil
	})

	if err != nil || !parsedToken.Valid {
		return nil, err
	}

	return claims, nil
}
