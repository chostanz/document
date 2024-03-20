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

func generateFormNumberDA(documentID int64, recursionCount int) (string, error) {
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
	formNumberWithDivision := fmt.Sprintf("%s/%s/%s/%s/%d", formNumberString, "PED", "F", romanMonth, year)

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

func AddDA(addDA models.Form, isPublished bool, username string, userID int, divisionCode string, recrusionCount int, data models.DampakAnalisa, signatories []models.Signatory) error {
	var documentCode string
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()

	app_id := currentTimestamp + int64(uniqueID)

	uuidObj := uuid.New()
	uuidString := uuidObj.String()

	formStatus := "Draft"
	if isPublished {
		formStatus = "Published"
	}
	var documentID int64
	err := db.Get(&documentID, "SELECT document_id FROM document_ms WHERE document_uuid = $1", addDA.DocumentUUID)
	if err != nil {
		log.Println("Error getting document_id:", err)
		return err
	}

	err = db.Get(&documentCode, "SELECT document_code FROM document_ms WHERE document_uuid = $1", addDA.DocumentUUID)
	if err != nil {
		log.Println("Error getting document code:", err)
		return err
	}

	var projectID int64
	err = db.Get(&projectID, "SELECT project_id FROM project_ms WHERE project_uuid = $1", addDA.ProjectUUID)
	if err != nil {
		log.Println("Error getting project_id:", err)
		return err
	}
	// Generate form number based on document code
	formNumberDA, err := generateFormNumberDA(documentID, recrusionCount+1)
	if err != nil {
		// Handle error
		log.Println("Error generating form number:", err)
		return err
	}
	// Marshal ITCM struct to JSON
	daJSON, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling ITCM struct:", err)
		return err
	}

	_, err = db.NamedExec("INSERT INTO form_ms (form_id, form_uuid, document_id, user_id, form_name, form_number, form_ticket, form_status, form_data, is_project, project_id, created_by) VALUES (:form_id, :form_uuid, :document_id, :user_id, :form_name, :form_number, :form_ticket, :form_status, :form_data, :is_project, :project_id, :created_by)", map[string]interface{}{
		"form_id":     app_id,
		"form_uuid":   uuidString,
		"document_id": documentID,
		"user_id":     userID,
		"form_name":   addDA.FormName,
		"form_number": formNumberDA,
		"form_ticket": addDA.FormTicket,
		"form_status": formStatus,
		"form_data":   daJSON,
		"project_id":  projectID,
		"created_by":  username,
	})

	if err != nil {
		return err
	}

	for _, signatory := range signatories {
		uuidString := uuid.New().String() // Gunakan uuid.New() dari paket UUID yang diimpor
		_, err := db.NamedExec("INSERT INTO sign_form (sign_uuid, form_id, name, position, role_sign, created_by) VALUES (:sign_uuid, :form_id, :name, :position, :role_sign, :created_by)", map[string]interface{}{
			"sign_uuid":  uuidString,
			"form_id":    app_id,
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

func GetAllFormDA() ([]models.Formss, error) {
	rows, err := db.Query(`
		SELECT 
			f.form_uuid, f.form_name, f.form_number, f.form_ticket, f.form_status, f.is_approved, f.reason, f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
			d.document_name,
			p.project_name,
			f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
			(f.form_data->>'nama_analis')::text AS nama_analis,
			(f.form_data->>'jabatan')::text AS jabatan,
			(f.form_data->>'departemen')::text AS departemen,
			(f.form_data->>'jenis_perubahan')::text AS jenis_perubahan,
			(f.form_data->>'detail_dampak_perubahan')::text AS detail_dampak_perubahan,
			(f.form_data->>'rencana_pengembangan_perubahan')::text AS rencana_pengembangan_perubahan,
			(f.form_data->>'rencana_pengujian_perubahan_sistem')::text AS rencana_pengujian_perubahan_sistem,
			(f.form_data->>'rencana_rilis_perubahan_dan_implementasi')::text AS rencana_rilis_perubahan_dan_implementasi
			CASE
			WHEN f.is_approved IS NULL THEN 'Menunggu Disetujui'
			WHEN f.is_approved = false THEN 'Tidak Disetujui'
			WHEN f.is_approved = true THEN 'Disetujui'
		END AS ApprovalStatus -- Alias the CASE expression as ApprovalStatus
			FROM 
			form_ms f
		LEFT JOIN 
			document_ms d ON f.document_id = d.document_id
		LEFT JOIN 
			project_ms p ON f.project_id = p.project_id
			WHERE
			d.document_code = 'DA' 
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Slice to hold all form data
	var forms []models.Formss

	// Iterate through the rows
	for rows.Next() {
		// Scan the row into the Forms struct
		var form models.Formss
		err := rows.Scan(
			&form.FormUUID,
			&form.FormName,
			&form.FormNumber,
			&form.FormTicket,
			&form.FormStatus,
			&form.DocumentName,
			&form.ProjectName,
			&form.IsApprove,
			&form.Reason,
			&form.ApprovalStatus,
			&form.CreatedBy,
			&form.CreatedAt,
			&form.UpdatedBy,
			&form.UpdatedAt,
			&form.DeletedBy,
			&form.DeletedAt,
			&form.NamaAnalis,
			&form.Jabatan,
			&form.Departemen,
			&form.JenisPerubahan,
			&form.DetailDampakPerubahan,
			&form.RencanaPengembanganPerubahan,
			&form.RencanaPengujianPerubahanSistem,
			&form.RencanaRilisPerubahanDanImplementasi,
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

func GetSpecDA(id string) (models.Formss, error) {
	var specDA models.Formss

	err := db.Get(&specDA, `SELECT 
		f.form_uuid, f.form_name, f.form_number, f.form_ticket, f.form_status, f.is_approved, f.reason, f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
		d.document_name,
		p.project_name,
		f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
		(f.form_data->>'nama_analis')::text AS nama_analis,
		(f.form_data->>'jabatan')::text AS jabatan,
		(f.form_data->>'departemen')::text AS departemen,
		(f.form_data->>'jenis_perubahan')::text AS jenis_perubahan,
		(f.form_data->>'detail_dampak_perubahan')::text AS detail_dampak_perubahan,
		(f.form_data->>'rencana_pengembangan_perubahan')::text AS rencana_pengembangan_perubahan,
		(f.form_data->>'rencana_pengujian_perubahan_sistem')::text AS rencana_pengujian_perubahan_sistem,
		(f.form_data->>'rencana_rilis_perubahan_dan_implementasi')::text AS rencana_rilis_perubahan_dan_implementasi,
		CASE
		WHEN f.is_approved IS NULL THEN 'Menunggu Disetujui'
		WHEN f.is_approved = false THEN 'Tidak Disetujui'
		WHEN f.is_approved = true THEN 'Disetujui'
		FROM 
		form_ms f
	LEFT JOIN 
		document_ms d ON f.document_id = d.document_id
	LEFT JOIN 
		project_ms p ON f.project_id = p.project_id
	WHERE f.form_uuid = $1 and d.document_code = 'DA' 
`, id)

	if err != nil {
		return models.Formss{}, err
	}

	return specDA, nil
}

func GetSpecFormDA(id string) ([]models.Formss, error) {
	var signatories []models.Formss

	err := db.Select(&signatories, `SELECT 
		f.form_uuid, f.form_name, f.form_number, f.form_ticket, f.form_status,
		d.document_name,
		p.project_name,
		f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
		(f.form_data->>'nama_analis')::text AS nama_analis,
		(f.form_data->>'jabatan')::text AS jabatan,
		(f.form_data->>'departemen')::text AS departemen,
		(f.form_data->>'jenis_perubahan')::text AS jenis_perubahan,
		(f.form_data->>'detail_dampak_perubahan')::text AS detail_dampak_perubahan,
		(f.form_data->>'rencana_pengembangan_perubahan')::text AS rencana_pengembangan_perubahan,
		(f.form_data->>'rencana_pengujian_perubahan_sistem')::text AS rencana_pengujian_perubahan_sistem,
		(f.form_data->>'rencana_rilis_perubahan_dan_implementasi')::text AS rencana_rilis_perubahan_dan_implementasi
	FROM 
		form_ms f
	LEFT JOIN 
		document_ms d ON f.document_id = d.document_id
	LEFT JOIN 
		project_ms p ON f.project_id = p.project_id
		WHERE
		f.form_uuid = $1 and d.document_code = 'DA' 
`)
	if err != nil {
		return nil, err
	}
	return signatories, nil
}

func UpdateFormDA(updateDA models.Form, data models.DampakAnalisa, username string, userID int, isPublished bool, id string) (models.Form, error) {
	currentTime := time.Now()
	formStatus := "Draft"
	if isPublished {
		formStatus = "Published"
	}

	var projectID int64
	err := db.Get(&projectID, "SELECT project_id FROM project_ms WHERE project_uuid = $1", updateDA.ProjectUUID)
	if err != nil {
		log.Println("Error getting project_id:", err)
		return models.Form{}, err
	}

	daJSON, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling DampakAnalisa struct:", err)
		return models.Form{}, err
	}
	log.Println("DampakAnalisa JSON:", string(daJSON)) // Periksa hasil marshaling

	_, err = db.NamedExec("UPDATE form_ms SET user_id = :user_id, form_name = :form_name, form_ticket = :form_ticket, form_status = :form_status, form_data = :form_data, updated_by = :updated_by, updated_at = :updated_at WHERE form_uuid = :id AND form_status = 'Draft'", map[string]interface{}{
		"user_id":     userID,
		"form_name":   updateDA.FormName,
		"form_ticket": updateDA.FormTicket,
		"project_id":  projectID,
		"form_status": formStatus,
		"form_data":   daJSON,
		"updated_by":  username,
		"updated_at":  currentTime,
		"id":          id,
	})
	if err != nil {
		return models.Form{}, err
	}
	return updateDA, nil
}
