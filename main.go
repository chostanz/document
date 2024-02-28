package main

import (
	"document/routes"
	"document/utils"
	"time"

	cache "github.com/SporkHubr/echo-http-cache"
	"github.com/SporkHubr/echo-http-cache/adapter/memory"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := routes.Route()

	// Inisialisasi adapter memory
	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(10000000),
	)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Inisialisasi client caching
	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(10*time.Minute),
		cache.ClientWithRefreshKey("opn"),
	)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Gunakan middleware caching
	e.Use(cacheClient.Middleware())

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	//header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	// Deklarasikan variabel methods sebagai slice dari string
	// methods := []string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}

	// // Terapkan middleware CORS ke Echo
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"*"},
	// 	// Menggunakan variabel methods yang sudah dideklarasikan
	// 	AllowMethods:     methods,
	// 	AllowCredentials: true,
	// }))
	customValidator := &utils.CustomValidator{Validator: validator.New()}
	e.Validator = customValidator
	e.Logger.Fatal(e.Start(":1234"))

}
