package controller

import (
	"database/sql"
	"document/models"
	"document/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
)

func AddDA(c echo.Context) error {
	const maxRecursionCount = 1000
	recursionCount := 0 // Set nilai awal untuk recursionCount
	var addFormRequest struct {
		IsPublished bool                  `json:"isPublished"`
		FormData    models.Form           `json:"formData"`
		DA          service.DampakAnalisa `json:"data_da"` // Tambahkan ITCM ke dalam struct request
		Signatory   []models.Signatory    `json:"signatories"`
	}

	if err := c.Bind(&addFormRequest); err != nil {
		log.Print(err)
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	if len(addFormRequest.Signatory) == 0 || addFormRequest.DA == (service.DampakAnalisa{}) {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data boleh kosong!",
			Status:  false,
		})
	}

	fmt.Println("Nilai isPublished yang diterima di backend:", addFormRequest.IsPublished)

	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	// Periksa apakah tokenString mengandung "Bearer "
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Hapus "Bearer " dari tokenString
	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

	//dekripsi token JWE
	decrypted, err := DecryptJWE(tokenOnly, secretKey)
	if err != nil {
		fmt.Println("Gagal mendekripsi token:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		fmt.Println("Gagal mengurai klaim:", errJ)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}
	divisionCode := c.Get("division_code").(string)
	userID := c.Get("user_id").(int) // Mengambil userUUID dari konteks
	userName := c.Get("user_name").(string)
	addFormRequest.FormData.UserID = userID
	addFormRequest.FormData.Created_by = userName
	// addFormRequest.FormData.isProject = false
	// addFormRequest.FormData.projectCode =
	// Token yang sudah dideskripsi
	fmt.Println("Token yang sudah dideskripsi:", decrypted)
	fmt.Println("User ID:", userID)
	fmt.Println("User Name:", userName)
	fmt.Println("Division Code:", divisionCode)
	// Lakukan validasi token
	if userID == 0 && userName == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}

	// Validasi spasi untuk Code, Name, dan NumberFormat
	whitespace := regexp.MustCompile(`^\s`)
	if whitespace.MatchString(addFormRequest.FormData.FormTicket) || whitespace.MatchString(addFormRequest.FormData.FormNumber) {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Ticket atau Nomor tidak boleh dimulai dengan spasi!",
			Status:  false,
		})
	}

	errVal := c.Validate(&addFormRequest.FormData)
	//	addFormRequest.FormData.UserID = userID
	if errVal == nil {
		// Gunakan addFormRequest.IsPublished untuk menentukan apakah menyimpan sebagai draft atau mempublish
		addroleErr := service.AddDA(addFormRequest.FormData, addFormRequest.IsPublished, userName, userID, divisionCode, recursionCount, addFormRequest.DA, addFormRequest.Signatory)

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
			Message: "Berhasil menambahkan formulir!",
			Status:  true,
		})

	} else {
		fmt.Println(errVal)
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}
}

func GetAllFormDA(c echo.Context) error {
	form, err := service.GetAllFormDA()
	if err != nil {
		log.Print(err)
		response := models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal server. Mohon coba beberapa saat lagi",
			Status:  false,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	return c.JSON(http.StatusOK, form)

}

func GetSpecDA(c echo.Context) error {
	id := c.Param("id")

	var getDoc service.Forms

	getDoc, err := service.GetSpecDA(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(err)
			response := models.Response{
				Code:    404,
				Message: "Formulir DA tidak ditemukan!",
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

func GetSignatureForm(c echo.Context) error {
	id := c.Param("id")

	var getAppRole []models.Signatories

	getAppRole, err := service.GetSignatureForm(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(err)
			response := models.Response{
				Code:    404,
				Message: "Signatory tidak ditemukan!",
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

	return c.JSON(http.StatusOK, getAppRole)
}

func GetSpecSignatureByID(c echo.Context) error {
	id := c.Param("id")

	var getAppRole models.Signatorie

	getAppRole, err := service.GetSpecSignatureByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(err)
			response := models.Response{
				Code:    404,
				Message: "Signatory tidak ditemukan!",
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

	return c.JSON(http.StatusOK, getAppRole)
}
