package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SharedKeyAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method == "OPTIONS" {
			return next(c)
		}
		userKey := c.Request().Header.Get("X-Secret-Key")
		if userKey == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing X-Secret-Key header")
		}
		if userKey != secretKey {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid X-Secret-Key")
		}
		return next(c)
	}
}
