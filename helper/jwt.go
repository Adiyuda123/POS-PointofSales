package helper

import (
	"POS-PointofSales/app/config"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echoJWT "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	Id int `json:"id"`
	jwt.RegisteredClaims
}

func JWTMiddleware() echo.MiddlewareFunc {
	return echoJWT.WithConfig(echoJWT.Config{
		SigningKey:    []byte(config.JWT),
		SigningMethod: "HS256",
	})
}

func GenerateToken(userId uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userId"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWT))
}

func DecodeToken(e echo.Context) uint {
	user := e.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		userId := claims["userId"].(float64)
		return uint(userId)
	}
	return 0
}

func TrimPrefixHeaderToken(reqToken string) string {
	prefix := "Bearer "
	return strings.TrimPrefix(reqToken, prefix)
}

func ValidateToken(c echo.Context) error {
	reqToken := c.Request().Header.Get("Authorization")
	tokenString := TrimPrefixHeaderToken(reqToken)

	if tokenString == "" {
		return errors.New("request does not contain a valid token")
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&JwtCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.InitConfig().JWT), nil
		},
	)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(*JwtCustomClaims); !ok {
		return errors.New("couldn't parse claims")
	}

	return nil
}

func ClaimsToken(c echo.Context) JwtCustomClaims {
	reqToken := c.Request().Header.Get("Authorization")
	tokenString := TrimPrefixHeaderToken(reqToken)
	claims := &JwtCustomClaims{}
	jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT), nil
	})

	return *claims
}
