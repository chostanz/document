package routes

import (
	"document/controller"
	"document/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Error adalah tipe kustom yang mewakili kesalahan aplikasi
type Error struct {
	Code    int    // Kode status HTTP
	Message string // Pesan kesalahan
}

// Handler adalah tipe fungsi penanganan yang mengembalikan Error
type Handler func(http.ResponseWriter, *http.Request) *Error

// ServeHTTP menerapkan fungsi penanganan kustom ke Echo
func (fn Handler) ServeHTTP(c echo.Context) error {
	w := c.Response().Writer
	r := c.Request()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if e := fn(w, r); e != nil { // e is *Error
		return c.String(e.Code, e.Message)
	}
	return nil
}
func Route() *echo.Echo {
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			return next(c)
		}
	})
	superAdmin := e.Group("/superadmin")
	superAdmin.Use(middleware.SuperAdminMiddleware)

	adminMember := e.Group("/api")
	adminMember.Use(middleware.AdminMemberMiddleware)

	//admin
	adminGroup := e.Group("/admin")
	adminGroup.Use(middleware.AdminMiddleware)
	adminGroup.GET("/my/form/division", controller.FormByDivision)

	e.GET("/document", controller.GetAllDoc)
	e.GET("/document/:id", controller.ShowDocById)
	superAdmin.POST("/document/add", controller.AddDocument)
	superAdmin.PUT("/document/update/:id", controller.UpdateDocument)

	e.GET("/form", controller.GetAllForm)
	e.GET("/form/:id", controller.ShowFormById)
	adminMember.POST("/form/add", controller.AddForm)
	adminMember.GET("/my/form", controller.MyForm)
	adminMember.PUT("/form/update/:id", controller.UpdateForm)

	//dampak analisa
	adminMember.POST("/add/da", controller.AddDA)
	e.GET("/dampak/analisa", controller.GetAllFormDA)
	e.GET("/dampak/analisa/:id", controller.GetSpecDA)
	e.GET("/da/:id", controller.GetSpecAllDA)
	adminMember.PUT("/dampak/analisa/update/:id", controller.UpdateFormDA)
	adminMember.GET("/my/form/da", controller.GetAllFormDAbyUser)
	adminGroup.GET("/da/all", controller.GetAllDAbyAdmin)

	//tandatangan
	e.GET("/signatory/:id", controller.GetSpecSignatureByID)
	adminMember.PUT("/signature/update/:id", controller.UpdateSignature)
	e.GET("/form/signatories/:id", controller.GetSignatureForm)

	//FORM itcm
	adminMember.POST("/add/itcm", controller.AddITCM)
	e.GET("/form/itcm", controller.GetAllFormITCM)
	e.GET("/form/itcm/:id", controller.GetSpecITCM)
	e.GET("/itcm/:id", controller.GetSpecAllITCM)
	adminMember.PUT("/form/itcm/update/:id", controller.UpdateFormITCM)
	adminMember.GET("/my/form/itcm", controller.GetAllFormITCMbyUserID)
	adminGroup.GET("/itcm/all", controller.GetAllFormITCMAdmin)

	//add approval
	adminMember.PUT("/form/approval/:id", controller.AddApproval)
	//form BA
	adminMember.POST("/add/ba", controller.AddBA)
	e.GET("/form/ba", controller.GetAllFormBA)
	e.GET("/form/ba/:id", controller.GetSpecBA)
	e.GET("/ba/:id", controller.GetSpecAllBA)
	adminMember.GET("/my/form/ba", controller.GetAllFormBAbyUserID)
	adminGroup.GET("/ba/all", controller.GetAllFormBAAdmin)
	adminMember.PUT("/form/ba/update/:id", controller.UpdateFormBA)

	//product
	e.GET("/product", controller.GetAllProduct)
	e.GET("/product/:id", controller.ShowProductById)
	superAdmin.POST("/product/add", controller.AddProduct)
	superAdmin.PUT("/product/update/:id", controller.UpdateProdcut)

	//project
	e.GET("/project", controller.GetAllProject)
	e.GET("/project/:id", controller.ShowProjectById)
	superAdmin.POST("/project/add", controller.AddProject)
	superAdmin.PUT("/project/update/:id", controller.UpdateProject)

	//delete form
	adminMember.PUT("/form/delete/:id", controller.DeleteForm)

	return e
}
