package main

import (
	"net/http"
	"qr_auth/auth"
	"qr_auth/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.Setup()
	e := echo.New()
	// e.Use(middleware.CORS())
	// OR: Use custom CORS config (recommended for production)
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"*"}, // your frontend origin
        AllowMethods: []string{echo.GET, echo.POST},
    }))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/auth/qr-login", auth.SendQRLogin)
	e.GET("/auth/verify", auth.VerifyQRLogin)
	e.Logger.Fatal(e.Start(":1323"))
}
