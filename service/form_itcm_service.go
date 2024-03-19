package service

import (
	"database/sql"
	"document/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func generateFormNumberITCM(documentID int64, divisionCode string, recursionCount int) (string, error) {
	const maxRecursionCount = 1000

	// Check if the maximum recursion count is reached
	if recursionCount > maxRecursionCount {
		return "", errors.New("Maximum recursion count exceeded")
	}

	// Get the latest form number for the given document ID
	var latestFormNumber sql.NullString
	err := db.Get(&latestFormNumber, "SELECT MAX(form_number) FROM form_ms WHERE document_id = $1", documentID)
	if err != nil {
		return "", fmt.Errorf("Error getting latest form number: %v", err)
	}

	documentCode, err := GetDocumentCode(documentID)
	if err != nil {
		return "", fmt.Errorf("failed to get document code: %v", err)
	}
	// Initialize formNumber to 1 if latestFormNumber is NULL
	formNumber := 1
	if latestFormNumber.Valid {
		// Parse the latest form number
		var latestFormNumberInt int
		_, err := fmt.Sscanf(latestFormNumber.String, "%d", &latestFormNumberInt)
		if err != nil {
			return "", fmt.Errorf("Error parsing latest form number: %v", err)
		}
		// Increment the latest form number
		formNumber = latestFormNumberInt + 1
	}

	// Get current year and month
	year := time.Now().Year()
	month := time.Now().Month()

	// Convert month to Roman numeral
	romanMonth, err := convertToRoman(int(month))
	if err != nil {
		return "", fmt.Errorf("Error converting month to Roman numeral: %v", err)
	}

	// Format the form number according to the specified format
	formNumberString := fmt.Sprintf("%04d", formNumber)
	formNumberWithDivision := fmt.Sprintf("%s/%s/%s/%s/%d", formNumberString, "F", documentCode, romanMonth, year)

	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM form_ms WHERE form_number = $1 and document_id = $2", formNumberString, documentID)
	if err != nil {
		return "", fmt.Errorf("Error checking existing form number: %v", err)
	}
	if count > 0 {
		// If the form number already exists, recursively call the function again
		return generateFormNumberDA(documentID, recursionCount+1)
	}

	return formNumberWithDivision, nil
}
func AddITCM(addForm models.Form, itcm models.ITCM, isPublished bool, userID int, username string, divisionCode string, recursionCount int, signatories []models.Signatory) error {
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()
	appID := currentTimestamp + int64(uniqueID)
	uuidObj := uuid.New()
	uuidString := uuidObj.String()

	formStatus := "Draft"
	if isPublished {
		formStatus = "Published"
	}

	var documentID int64
	err := db.Get(&documentID, "SELECT document_id FROM document_ms WHERE document_uuid = $1", addForm.DocumentUUID)
	if err != nil {
		log.Println("Error getting document_id:", err)
		return err
	}

	var projectID int64
	err = db.Get(&projectID, "SELECT project_id FROM project_ms WHERE project_uuid = $1", addForm.ProjectUUID)
	if err != nil {
		log.Println("Error getting project_id:", err)
		return err
	}

	// var projectID int64
	// if isProject {
	// 	err = db.Get(&projectID, "SELECT project_id FROM project_ms WHERE project_code = $1", projectCode)
	// 	if err != nil {
	// 		log.Println("Error getting project_id:", err)
	// 		return err
	// 	}
	// }

	var documentCode string
	err = db.Get(&documentCode, "SELECT document_code FROM document_ms WHERE document_uuid = $1", addForm.DocumentUUID)
	if err != nil {
		log.Println("Error getting document code:", err)
		return err
	}

	// Generate form number based on document code
	// var formNumber string
	// if isProject {
	// 	// Format nomor formulir sesuai dengan proyek
	// 	formNumber, err = generateProjectFormNumber(documentID, projectID, recursionCount)
	// 	if err != nil {
	// 		log.Println("Error generating project form number:", err)
	// 		return err
	// 	}
	// } else {
	// 	// Format nomor formulir sesuai dengan divisi
	// 	formNumber, err = generateFormNumber(documentID, divisionCode, recursionCount+1)
	// 	if err != nil {
	// 		log.Println("Error generating division form number:", err)
	// 		return err
	// 	}
	// }

	formNumber, err := generateFormNumberITCM(documentID, divisionCode, recursionCount+1)
	if err != nil {
		log.Println("Error generating project form number:", err)
		return err
	}

	// Marshal ITCM struct to JSON
	itcmJSON, err := json.Marshal(itcm)
	if err != nil {
		log.Println("Error marshaling ITCM struct:", err)
		return err
	}

	_, err = db.NamedExec("INSERT INTO form_ms (form_id, form_uuid, document_id, user_id, project_id, form_name, form_number, form_ticket, form_status, form_data, is_project, created_by) VALUES (:form_id, :form_uuid, :document_id, :user_id, :project_id, :form_name, :form_number, :form_ticket, :form_status, :form_data, :is_project, :created_by)", map[string]interface{}{
		"form_id":     appID,
		"form_uuid":   uuidString,
		"document_id": documentID,
		"user_id":     userID,
		"project_id":  projectID,
		"form_name":   addForm.FormName,
		"form_number": formNumber,
		"form_ticket": addForm.FormTicket,
		"form_status": formStatus,
		"form_data":   itcmJSON, // Convert JSON to string
		"is_project":  addForm.IsProject,
		"created_by":  username,
	})

	if err != nil {
		return err
	}

	for _, signatory := range signatories {
		uuidString := uuid.New().String() // Gunakan uuid.New() dari paket UUID yang diimpor
		_, err := db.NamedExec("INSERT INTO sign_form (sign_uuid, form_id, name, position, role_sign, created_by) VALUES (:sign_uuid, :form_id, :name, :position, :role_sign, :created_by)", map[string]interface{}{
			"sign_uuid":  uuidString,
			"form_id":    appID,
			"name":       signatory.Name,
			"position":   signatory.Position,
			"role_sign":  signatory.Role,
			"created_by": username,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func GetAllFormITCM() ([]models.FormITCM, error) {
	rows, err := db.Query(`SELECT 
    f.form_uuid, f.form_name, 
    REPLACE(f.form_number, '/ITCM/', '/') AS formatted_form_number,
    f.form_ticket, f.form_status,
    d.document_name,
    p.project_name,
    f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
    (f.form_data->>'no_da')::text AS no_da,
    (f.form_data->>'nama_pemohon')::text AS nama_pemohon,
    (f.form_data->>'intansi')::text AS intansi,
    (f.form_data->>'tanggal')::text AS tanggal,
    (f.form_data->>'perubahan_aset')::text AS perubahan_aset,
    (f.form_data->>'deskripsi')::text AS deskripsi
   FROM 
    form_ms f
	LEFT JOIN 
    document_ms d ON f.document_id = d.document_id
	LEFT JOIN 
    project_ms p ON f.project_id = p.project_id
	WHERE
    d.document_code = 'ITCM'
`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Slice to hold all form data
	var forms []models.FormITCM

	// Iterate through the rows
	for rows.Next() {
		// Scan the row into the Forms struct
		var form models.FormITCM
		err := rows.Scan(
			&form.FormUUID,
			&form.FormName,
			&form.FormNumber,
			&form.FormTicket,
			&form.FormStatus,
			&form.DocumentName,
			&form.ProjectName,
			&form.CreatedBy,
			&form.CreatedAt,
			&form.UpdatedBy,
			&form.UpdatedAt,
			&form.DeletedBy,
			&form.DeletedAt,
			&form.NoDA,
			&form.NamaPemohon,
			&form.Intansi,
			&form.Tanggal,
			&form.PerubahanAset,
			&form.Deskripsi,
		)
		if err != nil {
			return nil, err
		}

		// Append the form data to the slice
		forms = append(forms, form)
	}
	// Return the forms as JSON response
	return forms, nil
}

func GetSpecITCM(id string) (models.FormITCM, error) {
	var specITCM models.FormITCM

	err := db.Get(&specITCM, `SELECT
	f.form_uuid, f.form_name, 
    REPLACE(f.form_number, '/ITCM/', '/') AS formatted_form_number,
    f.form_ticket, f.form_status,
    d.document_name,
    p.project_name,
    f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
    (f.form_data->>'no_da')::text AS no_da,
    (f.form_data->>'nama_pemohon')::text AS nama_pemohon,
    (f.form_data->>'intansi')::text AS intansi,
    (f.form_data->>'tanggal')::text AS tanggal,
    (f.form_data->>'perubahan_aset')::text AS perubahan_aset,
    (f.form_data->>'deskripsi')::text AS deskripsi
   FROM 
    form_ms f
	LEFT JOIN 
    document_ms d ON f.document_id = d.document_id
	LEFT JOIN 
    project_ms p ON f.project_id = p.project_id
	WHERE
    f.form_uuid = $1 AND d.document_code = 'ITCM'
	`, id)
	if err != nil {
		return models.FormITCM{}, err
	}

	return specITCM, nil

}
