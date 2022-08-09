package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORS())

	e.GET("/", hello)

	e.Logger.Fatal(e.Start(":8082"))
}

// ハンドラーを定義
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
