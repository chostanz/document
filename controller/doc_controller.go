package controller

import (
	"database/sql"
	"document/database"
	"document/models"
	"document/service"
	"log"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
)

var db = database.Connection()

func AddDocument(c echo.Context) error {
	var addDocument models.Document
	if err := c.Bind(&addDocument); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	// Validasi spasi untuk Code, Name, dan NumberFormat
	whitespace := regexp.MustCompile(`^\s`)
	if whitespace.MatchString(addDocument.Code) {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Code tidak boleh dimulai dengan spasi!",
			Status:  false,
		})
	}

	if whitespace.MatchString(addDocument.Name) {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Name tidak boleh dimulai dengan spasi!",
			Status:  false,
		})
	}

	if whitespace.MatchString(addDocument.NumberFormat) {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Format Nomor tidak boleh dimulai dengan spasi!",
			Status:  false,
		})
	}

	errVal := c.Validate(&addDocument)

	if errVal == nil {
		var existingDocumentID int
		err := db.QueryRow("SELECT document_id FROM document_ms WHERE (document_code = $1 OR document_name = $2) AND deleted_at IS NULL", addDocument.Code, addDocument.Name).Scan(&existingDocumentID)

		if err == nil {
			return c.JSON(http.StatusBadRequest, &models.Response{
				Code:    400,
				Message: "Gagal menambahkan document. Document sudah ada!",
				Status:  false,
			})
		} else {
			addroleErr := service.AddDocument(addDocument)
			if addroleErr != nil {
				log.Print(addroleErr)
				return c.JSON(http.StatusInternalServerError, &models.Response{
					Code:    500,
					Message: "Terjadi kesalahan internal pada server. Coba beberapa saat lagi",
					Status:  false,
				})
			}

			return c.JSON(http.StatusCreated, &models.Response{
				Code:    201,
				Message: "Berhasil menambahkan document!",
				Status:  true,
			})
		}
	} else {
		log.Print(errVal)
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}

}

func GetAllDoc(c echo.Context) error {
	documents, err := service.GetAllDoc()
	if err != nil {
		log.Print(err)
		response := models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal server. Mohon coba beberapa saat lagi",
			Status:  false,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	return c.JSON(http.StatusOK, documents)

}

func ShowDocById(c echo.Context) error {
	id := c.Param("id")

	var getDoc models.Document

	getDoc, err := service.ShowDocById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(err)
			response := models.Response{
				Code:    404,
				Message: "Document tidak ditemukan!",
				Status:  false,
			}
			return c.JSON(http.StatusNotFound, response)
		} else {
			log.Print(err)
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}
	}

	return c.JSON(http.StatusOK, getDoc)
}

func UpdateDocument(c echo.Context) error {
	id := c.Param("id")

	perviousContent, errGet := service.ShowDocById(id)
	if errGet != nil {
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Gagal mengupdate document. Document tidak ditemukan!",
			Status:  false,
		})
	}

	var editDoc models.Document
	if err := c.Bind(&editDoc); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data invalid!",
			Status:  false,
		})
	}
	whitespace := regexp.MustCompile(`^\s`)
	if whitespace.MatchString(editDoc.Code) {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Code tidak boleh dimulai dengan spasi!",
			Status:  false,
		})
	}

	if whitespace.MatchString(editDoc.Name) {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Name tidak boleh dimulai dengan spasi!",
			Status:  false,
		})
	}

	if whitespace.MatchString(editDoc.NumberFormat) {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Format Nomor tidak boleh dimulai dengan spasi!",
			Status:  false,
		})
	}

	err := c.Validate(&editDoc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}
	if err == nil {
		var existingDocumentID int
		err := db.QueryRow("SELECT document_id FROM document_ms WHERE (document_name = $1 OR document_code = $2) AND deleted_at IS NULL", editDoc.Name, editDoc.Code).Scan(&existingDocumentID)

		if err == nil {
			return c.JSON(http.StatusBadRequest, &models.Response{
				Code:    400,
				Message: "Document sudah ada! Document tidak boleh sama!",
				Status:  false,
			})
		}

		exsitingDoc, err := service.GetDocCodeName(id)
		if err != nil {
			log.Printf("Error getting existing user data: %v", err)
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server.",
				Status:  false,
			})
		}

		if editDoc.Code != exsitingDoc.Code || editDoc.Name != exsitingDoc.Name {
			isUnique, err := service.IsUniqueDoc(id, editDoc.Code, editDoc.Name)
			if err != nil {
				log.Println("Error checking uniqueness:", err)
				return c.JSON(http.StatusInternalServerError, &models.Response{
					Code:    500,
					Message: "Terjadi kesalahan internal pada server.",
					Status:  false,
				})
			}

			if !isUnique {
				log.Println("Document sudah ada! Document tidak boleh sama!")
				return c.JSON(http.StatusBadRequest, &models.Response{
					Code:    400,
					Message: "Document sudah ada! Document tidak boleh sama!",
					Status:  false,
				})
			}
		}

		_, errService := service.UpdateDocument(editDoc, id)
		if errService != nil {
			log.Println("Kesalahan selama pembaruan:", errService)
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}

		log.Println(perviousContent)
		return c.JSON(http.StatusOK, &models.Response{
			Code:    200,
			Message: "Document berhasil diperbarui!",
			Status:  true,
		})
	} else {
		log.Println("Kesalahan sebelum pembaruan:", err)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
			Status:  false,
		})
	}
}
