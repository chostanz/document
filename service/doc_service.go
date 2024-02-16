package service

import (
	"database/sql"
	"document/database"
	"document/models"
	"encoding/json"
	"fmt"
	"log"
	"time"

	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var db = database.Connection()

type JwtCustomClaims struct {
	UserUUID string `json:"user_uuid"`
	// AppRoleId          int `json:"application_role_id"`
	// DivisionId         int `json:"division_id"`
	jwt.StandardClaims // Embed the StandardClaims struct

}

func DecryptJWE(jweToken string, secretKey string) (string, error) {
	// Dekripsi token JWE
	decrypted, _, err := jose.Decode(jweToken, secretKey)
	if err != nil {
		return "", fmt.Errorf("Gagal mendekripsi token: %s, error: %v", jweToken, err)
	}
	return decrypted, nil
}

func GetUserInfoFromToken(tokenStr string) (string, error) {
	// // Dekripsi token JWE
	// secretKey := "secretJwToken"
	// decrypted, err := DecryptJWE(tokenStr, secretKey)
	// if err != nil {
	// 	return "", fmt.Errorf("Gagal mendekripsi token: %v", err)
	// }

	// Parse token JWT yang telah didekripsi
	log.Print("token str ", tokenStr)
	var claims JwtCustomClaims
	if err := json.Unmarshal([]byte(tokenStr), &claims); err != nil {
		return "", fmt.Errorf("Gagal mengurai klaim: %v", err)
	}

	// Mengambil nilai user_uuid dari klaim
	userUUID := claims.UserUUID
	log.Print("USER UUID : ", userUUID)
	return userUUID, nil
}

// func GetUserInfoFromToken(tokenStr string) (string, error) {
// 	secretKey := "secretJwToken" // Ganti dengan kunci yang benar

// 	// Dekripsi token JWE
// 	decrypted, err := DecryptJWE(tokenStr, secretKey)
// 	if err != nil {
// 		fmt.Println("Gagal mendekripsi token woy:", err)
// 		return "", err
// 	}

// 	// Parse token JWT yang telah didekripsi
// 	var claims JwtCustomClaims
// 	err = json.Unmarshal([]byte(decrypted), &claims)
// 	if err != nil {
// 		fmt.Println("Gagal mengurai klaim:", err)
// 		return "", err
// 	}

// 	// Mengambil nilai user_uuid dari klaim
// 	userUUID, ok := claims["user_uuid"].(string)
// 	if !ok {
// 		return "", errors.New("Tidak dapat menemukan user_uuid dalam token")
// 	}

// 	return userUUID, nil
// }

// // // Mengembalikan user_uuid dari token
// // return claims.UserUUID, nil
// parts := strings.Split(tokenStr, ".")

// // Mengurai bagian terenkripsi token (payload)
// payload, err := base64.RawURLEncoding.DecodeString(parts[1])
// if err != nil {
// 	return "", fmt.Errorf("Gagal mendekode payload token: %v", err)
// }

// fmt.Println("Payload yang diuraikan:", string(payload))

// // Parse payload sebagai struktur yang sesuai
// var claims map[string]interface{}
// if err := json.Unmarshal(payload, &claims); err != nil {
// 	return "", fmt.Errorf("Gagal mengurai klaim: %v", err)
// }

func AddDocument(addDocument models.Document, userUUID string) error {

	_, errP := GetUserInfoFromToken(userUUID)
	if errP != nil {
		return errP
	}
	// username, errP := GetUsernameByID(userUUID)
	// if errP != nil {
	// 	return errP
	// }
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()

	app_id := currentTimestamp + int64(uniqueID)

	uuid := uuid.New()
	uuidString := uuid.String()
	_, err := db.NamedExec("INSERT INTO document_ms (document_id, document_uuid, document_code, document_name, document_format_number, created_by) VALUES (:document_id, :document_uuid, :document_code, :document_name, :document_format_number, :created_by)", map[string]interface{}{
		"document_id":            app_id,
		"document_uuid":          uuidString,
		"document_code":          addDocument.Code,
		"document_name":          addDocument.Name,
		"document_format_number": addDocument.NumberFormat,
		"created_by":             "super admin",
	})
	if err != nil {
		return err
	}
	return nil
}

func GetAllDoc() ([]models.Document, error) {

	document := []models.Document{}
	rows, errSelect := db.Queryx("select document_uuid, document_order, document_code, document_name, document_format_number, created_by, created_at, updated_by, updated_at from document_ms WHERE deleted_at IS NULL")
	if errSelect != nil {
		return nil, errSelect
	}

	for rows.Next() {
		place := models.Document{}
		rows.StructScan(&place)
		document = append(document, place)
	}

	return document, nil
}
func ShowDocById(id string) (models.Document, error) {
	var document models.Document

	err := db.Get(&document, "SELECT document_uuid, document_order, document_code, document_name,document_format_number, created_by, created_at, updated_by, updated_at FROM document_ms WHERE document_uuid = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return models.Document{}, err
	}
	return document, nil

}

func GetDocCodeName(uuid string) (models.DocCodeName, error) {
	var docCodeName models.DocCodeName

	err := db.Get(&docCodeName, "SELECT document_code, document_name FROM document_ms WHERE document_uuid = $1", uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			// Tidak ada baris yang sesuai
			log.Println("No rows found for role_uuid:", uuid)
			return models.DocCodeName{}, err
		}

		// Terjadi kesalahan lain
		log.Println("Error getting role data by role_ms:", err)
		return models.DocCodeName{}, err
	}

	return docCodeName, nil
}

func IsUniqueDoc(uuid, code, name string) (bool, error) {
	var count int

	var exsitingDocCode, exsitingDocName string
	err := db.Get(&exsitingDocCode, "SELECT document_code FROM document_ms WHERE document_uuid = $1", uuid)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	err = db.Get(&exsitingDocName, "SELECT document_name FROM document_ms WHERE document_uuid = $1", uuid)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if code == exsitingDocCode && name == exsitingDocName {
		return true, nil
	}

	err = db.Get(&count, "SELECT COUNT(*) FROM document_ms WHERE (document_code = $1 OR document_name = $2) AND document_uuid != $3 AND deleted_at IS NULL", code, name, uuid)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}
func GetDocumentIDByName(name string) (int, error) {
	var documentID int
	err := db.QueryRow("SELECT document_id FROM document_ms WHERE document_name = $1 AND deleted_at IS NULL", name).Scan(&documentID)
	return documentID, err
}

func GetDocumentIDByCode(code string) (int, error) {
	var documentID int
	err := db.QueryRow("SELECT document_id FROM document_ms WHERE document_code = $1 AND deleted_at IS NULL", code).Scan(&documentID)
	return documentID, err
}
func UpdateDocument(updateDoc models.Document, id string) (models.Document, error) {
	// username, errUser := GetUsernameByID(userUUID)
	// if errUser != nil {
	// 	log.Print(errUser)
	// 	return models.Document{}, errUser

	// }

	currentTime := time.Now()

	_, err := db.NamedExec("UPDATE document_ms SET document_code = :document_code, document_name = :document_name, document_format_number = :document_format_number, updated_by = :updated_by, updated_at = :updated_at WHERE document_uuid = :id", map[string]interface{}{
		"document_code":          updateDoc.Code,
		"document_name":          updateDoc.Name,
		"document_format_number": updateDoc.NumberFormat,
		"updated_by":             updateDoc.Updated_by,
		"updated_at":             currentTime,
		"id":                     id,
	})
	if err != nil {
		log.Print(err)
		return models.Document{}, err
	}
	return updateDoc, nil
}
