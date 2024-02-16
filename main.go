package main

import (
	"document/routes"
	"document/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := routes.Route()

	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4200"},                    // Izinkan akses dari origin ini
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE}, // Izinkan metode permintaan tertentu
	}))
	customValidator := &utils.CustomValidator{Validator: validator.New()}
	e.Validator = customValidator
	e.Logger.Fatal(e.Start(":1234"))

}
