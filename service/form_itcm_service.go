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

	var documentCode string
	err = db.Get(&documentCode, "SELECT document_code FROM document_ms WHERE document_uuid = $1", addForm.DocumentUUID)
	if err != nil {
		log.Println("Error getting document code:", err)
		return err
	}

	formNumber, err := generateFormNumberITCM(documentID, divisionCode, recursionCount+1)
	if err != nil {
		log.Println("Error generating project form number:", err)
		return err
	}

	itcmJSON, err := json.Marshal(itcm)
	if err != nil {
		log.Println("Error marshaling ITCM struct:", err)
		return err
	}

	_, err = db.NamedExec("INSERT INTO form_ms (form_id, form_uuid, document_id, user_id, project_id, form_number, form_ticket, form_status, form_data, created_by) VALUES (:form_id, :form_uuid, :document_id, :user_id, :project_id,:form_number, :form_ticket, :form_status, :form_data, :created_by)", map[string]interface{}{
		"form_id":     appID,
		"form_uuid":   uuidString,
		"document_id": documentID,
		"user_id":     userID,
		"project_id":  projectID,
		"form_number": formNumber,
		"form_ticket": addForm.FormTicket,
		"form_status": formStatus,
		"form_data":   itcmJSON,
		"created_by":  username,
	})

	if err != nil {
		return err
	}

	for _, signatory := range signatories {
		uuidString := uuid.New().String()
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

// menampilkan form tanpa token
func GetAllFormITCM() ([]models.FormsITCM, error) {
	rows, err := db.Query(`SELECT
		f.form_uuid,
		REPLACE(f.form_number, '/ITCM/', '/') AS formatted_form_number,
		f.form_ticket, f.form_status,
		d.document_name,
		p.project_name,
		p.project_manager,
		CASE
			WHEN f.is_approve IS NULL THEN 'Menunggu Disetujui'
			WHEN f.is_approve = false THEN 'Tidak Disetujui'
			WHEN f.is_approve = true THEN 'Disetujui'
		END AS ApprovalStatus,
		f.reason, f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
		(f.form_data->>'no_da')::text AS no_da,
		(f.form_data->>'nama_pemohon')::text AS nama_pemohon,
		(f.form_data->>'instansi')::text AS instansi,
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
	var forms []models.FormsITCM
	//rows, err := db.Query(&forms, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var form models.FormsITCM
		err := rows.Scan(
			&form.FormUUID,
			&form.FormNumber,
			&form.FormTicket,
			&form.FormStatus,
			&form.DocumentName,
			&form.ProjectName,
			&form.ProjectManager,
			&form.ApprovalStatus,
			&form.Reason,
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

		forms = append(forms, form)
	}

	return forms, nil
}

// menampilkan form berdasar user/ milik dia sendiri
func GetAllITCMbyUserID(userID int) ([]models.FormsITCM, error) {
	rows, err := db.Query(`SELECT
		f.form_uuid,
		REPLACE(f.form_number, '/ITCM/', '/') AS formatted_form_number,
		f.form_ticket, f.form_status,
		d.document_name,
		p.project_name,
		p.project_manager,
		CASE
			WHEN f.is_approve IS NULL THEN 'Menunggu Disetujui'
			WHEN f.is_approve = false THEN 'Tidak Disetujui'
			WHEN f.is_approve = true THEN 'Disetujui'
		END AS ApprovalStatus,
		f.reason, f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
		(f.form_data->>'no_da')::text AS no_da,
		(f.form_data->>'nama_pemohon')::text AS nama_pemohon,
		(f.form_data->>'instansi')::text AS instansi,
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
			f.user_id = $1 AND d.document_code = 'ITCM'
			`, userID)
	var forms []models.FormsITCM
	//rows, err := db.Query(&forms, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var form models.FormsITCM
		err := rows.Scan(
			&form.FormUUID,
			&form.FormNumber,
			&form.FormTicket,
			&form.FormStatus,
			&form.DocumentName,
			&form.ProjectName,
			&form.ProjectManager,
			&form.ApprovalStatus,
			&form.Reason,
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

		forms = append(forms, form)
	}

	return forms, nil

}

// menampilkan form dari admin
func GetAllFormITCMAdmin() ([]models.FormsITCM, error) {
	rows, err := db.Query(`SELECT
		f.form_uuid,
		REPLACE(f.form_number, '/ITCM/', '/') AS formatted_form_number,
		f.form_ticket, f.form_status,
		d.document_name,
		p.project_name,
		p.project_manager,
		CASE
			WHEN f.is_approve IS NULL THEN 'Menunggu Disetujui'
			WHEN f.is_approve = false THEN 'Tidak Disetujui'
			WHEN f.is_approve = true THEN 'Disetujui'
		END AS ApprovalStatus,
		f.reason, f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
		(f.form_data->>'no_da')::text AS no_da,
		(f.form_data->>'nama_pemohon')::text AS nama_pemohon,
		(f.form_data->>'instansi')::text AS instansi,
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
	var forms []models.FormsITCM
	//rows, err := db.Query(&forms, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var form models.FormsITCM
		err := rows.Scan(
			&form.FormUUID,
			&form.FormNumber,
			&form.FormTicket,
			&form.FormStatus,
			&form.DocumentName,
			&form.ProjectName,
			&form.ProjectManager,
			&form.ApprovalStatus,
			&form.Reason,
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

		forms = append(forms, form)
	}

	return forms, nil
}

// untuk mengambil data spesifik hanya data form_ms
func GetSpecITCM(id string) (models.FormITCM, error) {
	var specITCM models.FormITCM

	err := db.Get(&specITCM, `SELECT 
	f.form_uuid, 
	REPLACE(f.form_number, '/ITCM/', '/') AS formatted_form_number,
	f.form_ticket, f.form_status,
	d.document_name,
	p.project_name,
	p.project_manager,
	CASE
		WHEN f.is_approve IS NULL THEN 'Menunggu Disetujui'
		WHEN f.is_approve = false THEN 'Tidak Disetujui'
		WHEN f.is_approve = true THEN 'Disetujui'
	END AS ApprovalStatus,
	f.reason, f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
	(f.form_data->>'no_da')::text AS no_da,
	(f.form_data->>'nama_pemohon')::text AS nama_pemohon,
	(f.form_data->>'instansi')::text AS instansi,
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

// untuk mengambil data spesifik dari frm_ms dan sign_form
func GetSpecAllITCM(id string) ([]models.FormITCMAll, error) {
	var speccITCM []models.FormITCMAll

	err := db.Select(&speccITCM, `SELECT
    f.form_uuid,
    REPLACE(f.form_number, '/ITCM/', '/') AS formatted_form_number,
    f.form_ticket,
    f.form_status,
    d.document_name,
    p.project_name,
    p.project_manager,
    CASE
        WHEN f.is_approve IS NULL THEN 'Menunggu Disetujui'
        WHEN f.is_approve = false THEN 'Tidak Disetujui'
        WHEN f.is_approve = true THEN 'Disetujui'
    END AS ApprovalStatus,
    f.reason,
    f.created_by,
    f.created_at,
    f.updated_by,
    f.updated_at,
    f.deleted_by,
    f.deleted_at,
    (f.form_data->>'no_da')::text AS no_da,
    (f.form_data->>'nama_pemohon')::text AS nama_pemohon,
    (f.form_data->>'instansi')::text AS instansi,
    (f.form_data->>'tanggal')::text AS tanggal,
    (f.form_data->>'perubahan_aset')::text AS perubahan_aset,
    (f.form_data->>'deskripsi')::text AS deskripsi,
    sf.sign_uuid AS sign_uuid,
    sf.name AS name,
    sf.position AS position,
    sf.role_sign AS role_sign
FROM
    form_ms f
LEFT JOIN 
    document_ms d ON f.document_id = d.document_id
LEFT JOIN 
    project_ms p ON f.project_id = p.project_id
LEFT JOIN
    sign_form sf ON f.form_id = sf.form_id
WHERE
    f.form_uuid = $1 AND d.document_code = 'ITCM'
	`, id)

	if err != nil {
		return nil, err
	}

	return speccITCM, nil
}

func UpdateFormITCM(updateITCM models.Form, data models.ITCM, username string, userID int, isPublished bool, id string) (models.Form, error) {
	currentTime := time.Now()
	formStatus := "Draft"
	if isPublished {
		formStatus = "Published"
	}

	var projectID int64
	err := db.Get(&projectID, "SELECT project_id FROM project_ms WHERE project_uuid = $1", updateITCM.ProjectUUID)
	if err != nil {
		log.Println("Error getting project_id:", err)
		return models.Form{}, err
	}

	daJSON, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling DampakAnalisa struct:", err)
		return models.Form{}, err
	}
	log.Println("ITCM JSON:", string(daJSON))

	_, err = db.NamedExec("UPDATE form_ms SET user_id = :user_id, form_ticket = :form_ticket, form_status = :form_status, form_data = :form_data, updated_by = :updated_by, updated_at = :updated_at WHERE form_uuid = :id AND form_status = 'Draft'", map[string]interface{}{
		"user_id":     userID,
		"form_ticket": updateITCM.FormTicket,
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
	return updateITCM, nil

}
