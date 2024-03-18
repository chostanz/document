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

type ITCM struct {
	PerubahanAset string `json:"perubahan_aset" validate:"required"`
	Deskripsi     string `json:"deskripsi" validate:"required"`
}
type DampakAnalisa struct {
	//NamaProyek                           string    `json:"nama_proyek"`
	NamaAnalis                           string `json:"nama_analis"`
	Jabatan                              string `json:"jabatan"`
	Departemen                           string `json:"departemen"`
	JenisPerubahan                       string `json:"jenis_perubahan"`
	DetailDampakPerubahan                string `json:"detail_dampak_perubahan"`
	RencanaPengembanganPerubahan         string `json:"rencana_pengembangan_perubahan"`
	RencanaPengujianPerubahanSistem      string `json:"rencana_pengujian_perubahan_sistem"`
	RencanaRilisPerubahanDanImplementasi string `json:"rencana_rilis_perubahan_dan_implementasi"`
}

type FormWithApprovalStatus struct {
	FormID              int64         `json:"form_id"`
	FormUUID            string        `json:"form_uuid"`
	FormNumber          string        `json:"form_number"`
	FormTicket          string        `json:"form_ticket"`
	FormStatus          string        `json:"form_status"`
	FormData            DampakAnalisa `json:"form_data"`
	SignatoriesName     []string      `json:"signatories_name"`
	SignatoriesPosition []string      `json:"signatories_position"`
	SignatoriesRole     []string      `json:"signatories_role"`
	IsSigns             []bool        `json:"is_signs"`
	IsApproves          []bool        `json:"is_approves"`
	ApprovalStatus      string        `json:"approval_status"`
}

type Forms struct {
	FormUUID                             string         `json:"form_uuid" db:"form_uuid"`
	FormName                             string         `json:"form_name" db:"form_name"`
	FormNumber                           string         `json:"form_number" db:"form_number"`
	FormTicket                           string         `json:"form_ticket" db:"form_ticket"`
	FormStatus                           string         `json:"form_status" db:"form_status"`
	DocumentName                         string         `json:"document_name" db:"document_name"`
	ProjectName                          string         `json:"project_name" db:"project_name"`
	CreatedBy                            string         `json:"created_by" db:"created_by"`
	CreatedAt                            time.Time      `json:"created_at" db:"created_at"`
	UpdatedBy                            sql.NullString `json:"updated_by" db:"updated_by"`
	UpdatedAt                            sql.NullTime   `json:"updated_at" db:"updated_at"`
	DeletedBy                            sql.NullString `json:"deleted_by" db:"deleted_by"`
	DeletedAt                            sql.NullTime   `json:"deleted_at" db:"deleted_at"`
	NamaAnalis                           string         `json:"nama_analis" db:"nama_analis"`
	Jabatan                              string         `json:"jabatan" db:"jabatan"`
	Departemen                           string         `json:"departemen" db:"departemen"`
	JenisPerubahan                       string         `json:"jenis_perubahan" db:"jenis_perubahan"`
	DetailDampakPerubahan                string         `json:"detail_dampak_perubahan" db:"detail_dampak_perubahan"`
	RencanaPengembanganPerubahan         string         `json:"rencana_pengembangan_perubahan" db:"rencana_pengembangan_perubahan"`
	RencanaPengujianPerubahanSistem      string         `json:"rencana_pengujian_perubahan_sistem" db:"rencana_pengujian_perubahan_sistem"`
	RencanaRilisPerubahanDanImplementasi string         `json:"rencana_rilis_perubahan_dan_implementasi" db:"rencana_rilis_perubahan_dan_implementasi"`
}

func AddITCM(addForm models.Form, itcm ITCM, isPublished bool, userID int, username string, divisionCode string, recursionCount int, isProject bool, projectCode string, projectUUID string) error {
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()
	appID := currentTimestamp + int64(uniqueID)
	uuid := uuid.New()
	uuidString := uuid.String()

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
	if isProject {
		err = db.Get(&projectID, "SELECT project_id FROM project_ms WHERE project_code = $1", projectCode)
		if err != nil {
			log.Println("Error getting project_id:", err)
			return err
		}
	}

	var documentCode string
	err = db.Get(&documentCode, "SELECT document_code FROM document_ms WHERE document_uuid = $1", addForm.DocumentUUID)
	if err != nil {
		log.Println("Error getting document code:", err)
		return err
	}

	// Generate form number based on document code
	var formNumber string
	if isProject {
		// Format nomor formulir sesuai dengan proyek
		formNumber, err = generateProjectFormNumber(documentID, projectID, recursionCount)
		if err != nil {
			log.Println("Error generating project form number:", err)
			return err
		}
	} else {
		// Format nomor formulir sesuai dengan divisi
		formNumber, err = generateFormNumber(documentID, divisionCode, recursionCount+1)
		if err != nil {
			log.Println("Error generating division form number:", err)
			return err
		}
	}

	// Marshal ITCM struct to JSON
	itcmJSON, err := json.Marshal(itcm)
	if err != nil {
		log.Println("Error marshaling ITCM struct:", err)
		return err
	}

	_, err = db.NamedExec("INSERT INTO form_ms (form_id, form_uuid, document_id, user_id, project_id, form_name, form_number, form_ticket, form_status, form_data, created_by) VALUES (:form_id, :form_uuid, :document_id, :user_id, :project_id, :form_name, :form_number, :form_ticket, :form_status, :form_data, :created_by)", map[string]interface{}{
		"form_id":     appID,
		"form_uuid":   uuidString,
		"document_id": documentID,
		"user_id":     userID,
		"project_id":  projectID,
		"form_name":   addForm.FormName,
		"form_number": formNumber,
		"form_ticket": addForm.FormTicket,
		"form_status": formStatus,
		"form_data":   string(itcmJSON), // Convert JSON to string
		"created_by":  username,
	})

	if err != nil {
		return err
	}

	return nil
}

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

func AddDA(addDA models.Form, isPublished bool, username string, userID int, divisionCode string, recrusionCount int, data DampakAnalisa, signatories []models.Signatory) error {
	var documentCode string
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()

	app_id := currentTimestamp + int64(uniqueID)

	uuid := uuid.New()
	uuidString := uuid.String()

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
		"is_project":  addDA.IsProject,
		"project_id":  projectID,
		"created_by":  username,
	})

	if err != nil {
		return err
	}

	for _, signatory := range signatories {
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

func GetAllFormDA() ([]Forms, error) {
	rows, err := db.Query(`
		SELECT 
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
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Slice to hold all form data
	var forms []Forms

	// Iterate through the rows
	for rows.Next() {
		// Scan the row into the Forms struct
		var form Forms
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

func GetSpecDA(id string) (Forms, error) {
	var specDA Forms

	err := db.Get(&specDA, `SELECT 
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
	WHERE f.form_uuid = $1
`, id)

	if err != nil {
		return Forms{}, err
	}

	return specDA, nil
}

func GetSpecFormDA(id string) ([]Forms, error) {
	var signatories []Forms

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
		f.form_uuid = $1
`)
	if err != nil {
		return nil, err
	}
	return signatories, nil
}

func GetSignatureForm(id string) ([]models.Signatories, error) {
	var signatories []models.Signatories

	err := db.Select(&signatories, `SELECT 
	sf.sign_uuid, 
	sf.name, 
	sf.position, 
	sf.role_sign, 
	sf.is_sign, 
	sf.is_approve, 
	sf.created_by, 
	sf.created_at, 
	sf.updated_by, 
	sf.updated_at, 
	sf.deleted_by, 
	sf.deleted_at,
	CASE
		WHEN sf.is_approve IS NULL THEN 'Menunggu Disetujui'
		WHEN sf.is_approve = false THEN 'Tidak Disetujui'
		WHEN sf.is_approve = true THEN 'Disetujui'
	END AS ApprovalStatus -- Alias the CASE expression as ApprovalStatus
FROM 
	sign_form sf 
	JOIN form_ms fm ON sf.form_id = fm.form_id 
WHERE 
	fm.form_uuid = $1`, id)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return signatories, nil

}

func GetSpecSignatureByID(id string) (models.Signatorie, error) {
	var signatories models.Signatorie
	err := db.Get(&signatories, "SELECT sign_uuid, name, position, role_sign, is_sign, is_approve, created_by, created_at, updated_by, updated_at, deleted_by, deleted_at FROM sign_form sf WHERE sign_uuid = $1", id)
	if err != nil {
		log.Print(err)
		return models.Signatorie{}, err
	}

	return signatories, nil
}

func UpdateFormSignature(updateSign models.UpdateSign, id string, username string) error {
	currentTime := time.Now()

	_, err := db.NamedExec("UPDATE sign_form SET is_sign = :is_sign, is_approve = :is_approve, reason_not_approve = :reason_not_approve, updated_by = :updated_by, updated_at = :updated_at, WHERE sign_uuid = :id", map[string]interface{}{
		"is_sign":            updateSign.IsSign,
		"is_approve":         updateSign.IsApprove,
		"reason_not_approve": updateSign.ReasonNotApprove,
		"updated_by":         username,
		"updated_at":         currentTime,
		"id":                 id,
	})
	if err != nil {
		return err
	}
	return nil
}
