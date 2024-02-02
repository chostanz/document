package routes

import (
	"document/controller"

	"github.com/labstack/echo/v4"
)

func Route() *echo.Echo {
	e := echo.New()

	e.GET("/document", controller.GetAllDoc)
	e.GET("/document/:id", controller.ShowDocById)
	e.POST("/document/add", controller.AddDocument)
	e.PUT("/document/update/:id", controller.UpdateDocument)

	e.GET("/form", controller.GetAllForm)
	e.GET("/form/:id", controller.ShowFormById)
	e.POST("/form/add", controller.AddForm)
	e.PUT("/form/update/:id", controller.UpdateForm)

	return e
}
