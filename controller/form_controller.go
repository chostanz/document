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

func AddForm(c echo.Context) error {
	var addFormRequest struct {
		IsPublished bool        `json:"isPublished"`
		FormData    models.Form `json:"formData"`
	}

	if err := c.Bind(&addFormRequest); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

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
	userID := c.Get("user_id").(int) // Mengambil userUUID dari konteks
	userName := c.Get("user_name").(string)
	addFormRequest.FormData.UserID = userID
	addFormRequest.FormData.Created_by = userName
	// Token yang sudah dideskripsi
	fmt.Println("Token yang sudah dideskripsi:", decrypted)
	fmt.Println("User ID:", userID)
	fmt.Println("User Name:", userName)
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
		addroleErr := service.AddForm(addFormRequest.FormData, addFormRequest.IsPublished, userName, userID)
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

func GetAllForm(c echo.Context) error {
	form, err := service.GetAllForm()
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

func MyForm(c echo.Context) error {
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
	userID := c.Get("user_id").(int)

	fmt.Println("User ID :", userID)

	myform, err := service.MyForm(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(err)
			response := models.Response{
				Code:    404,
				Message: "Form tidak ditemukan!",
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
	return c.JSON(http.StatusOK, myform)
}

func ShowFormById(c echo.Context) error {
	id := c.Param("id")

	var getDoc models.Forms

	getDoc, err := service.ShowFormById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(err)
			response := models.Response{
				Code:    404,
				Message: "Formulir tidak ditemukan!",
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

func UpdateForm(c echo.Context) error {
	id := c.Param("id")
	var updateFormRequest struct {
		IsPublished bool        `json:"isPublished"`
		FormData    models.Form `json:"formData"`
	}

	if err := c.Bind(&updateFormRequest); err != nil {
		log.Print("error saat binding:", err)
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

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
	userID := c.Get("user_id").(int) // Mengambil userUUID dari konteks
	userName := c.Get("user_name").(string)

	updateFormRequest.FormData.UserID = userID

	var updatedBy sql.NullString
	if userName != "" {
		updatedBy.String = userName
		updatedBy.Valid = true
	} else {
		updatedBy.Valid = false
	}

	updateFormRequest.FormData.Updated_by = updatedBy

	// Token yang sudah dideskripsi
	fmt.Println("Token yang sudah dideskripsi:", decrypted)
	fmt.Println("User ID:", userID)
	fmt.Println("user name: ", userName)

	// Lakukan validasi token
	if userID == 0 && userName == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}

	whitespace := regexp.MustCompile(`^\s`)
	if whitespace.MatchString(updateFormRequest.FormData.FormTicket) {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Ticket tidak boleh dimulai dengan spasi!",
			Status:  false,
		})
	}

	if whitespace.MatchString(updateFormRequest.FormData.FormNumber) {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Name tidak boleh dimulai dengan spasi!",
			Status:  false,
		})
	}

	if err := c.Validate(&updateFormRequest.FormData); err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}

	// Tidak perlu melakukan pemeriksaan err lagi di sini

	previousContent, errGet := service.ShowFormById(id)
	if errGet != nil {
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Gagal mengupdate formulir. Formulir tidak ditemukan!",
			Status:  false,
		})
	}
	if previousContent.FormStatus == "Published" {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Tidak dapat memperbarui dokumen yang sudah dipublish",
			Status:  false,
		})
	}

	_, errService := service.UpdateForm(updateFormRequest.FormData, id, updateFormRequest.IsPublished, userName, userID)
	if errService != nil {
		log.Println("Kesalahan selama pembaruan:", errService)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
			Status:  false,
		})
	}

	log.Println(previousContent)
	return c.JSON(http.StatusOK, &models.Response{
		Code:    200,
		Message: "Formulir berhasil diperbarui!",
		Status:  true,
	})
}
