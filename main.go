package main

import (
	"document/routes"
	"document/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4/middleware"
)

// type Error struct {
// 	Code    int    // Kode status HTTP
// 	Message string // Pesan kesalahan
// }

// // Handler adalah tipe fungsi penanganan yang mengembalikan Error
// type Handler func(http.ResponseWriter, *http.Request) *Error

// // ServeHTTP menerapkan fungsi penanganan kustom ke Echo
// // ServeHTTP menerapkan fungsi penanganan kustom ke Echo
// func (fn Handler) ServeHTTP(c echo.Context) error {
// 	w := c.Response().Writer
// 	r := c.Request()
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
// 	if e := fn(w, r); e != nil { // e is *Error
// 		return c.String(e.Code, e.Message)
// 	}
// 	return nil
// }

func main() {
	e := routes.Route()

	e.Use(middleware.Logger())
	//header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	// Deklarasikan variabel methods sebagai slice dari string
	methods := []string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}

	// Terapkan middleware CORS ke Echo
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		// Menggunakan variabel methods yang sudah dideklarasikan
		AllowMethods:     methods,
		AllowCredentials: true,
	}))
	customValidator := &utils.CustomValidator{Validator: validator.New()}
	e.Validator = customValidator
	e.Logger.Fatal(e.Start(":1234"))

}
