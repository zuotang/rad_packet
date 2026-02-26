package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"red_packet/backend/internal/config"
)

const CtxUserID = "user_id"

func JWT(cfg config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"code":    "UNAUTHORIZED",
					"message": "missing token",
				})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWT.Secret), nil
			})
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"code":    "UNAUTHORIZED",
					"message": "invalid token",
				})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"code":    "UNAUTHORIZED",
					"message": "invalid claims",
				})
			}
			rawUID, ok := claims["uid"]
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"code":    "UNAUTHORIZED",
					"message": "missing uid",
				})
			}

			uid := uint(rawUID.(float64))
			c.Set(CtxUserID, uid)
			return next(c)
		}
	}
}
