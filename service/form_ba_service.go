package service

import (
	"document/models"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

func AddBA(addForm models.Form, ba models.BA, isPublished bool, userID int, username string, divisionCode string, recursionCount int, signatories []models.Signatory) error {
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

	formNumber, err := generateFormNumber(documentID, divisionCode, recursionCount+1)
	if err != nil {
		log.Println("Error generating project form number:", err)
		return err
	}

	// Marshal ITCM struct to JSON
	baJSON, err := json.Marshal(ba)
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
		"form_data":   baJSON, // Convert JSON to string
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

func GetAllFormBA() ([]models.FormsBA, error) {
	rows, err := db.Query(`
		SELECT 
			f.form_uuid, f.form_name, f.form_number, f.form_ticket, f.form_status,
			d.document_name,
			p.project_name,
			f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
			(f.form_data->>'judul')::text AS judul,
			(f.form_data->>'tanggal')::text AS tanggal,
			(f.form_data->>'dokumen_da')::text AS dokumen_da,
			(f.form_data->>'dokumen_itcm')::text AS dokumen_itcm,
			(f.form_data->>'dilakukan_oleh')::text AS dilakukan_oleh,
			(f.form_data->>'didampingi_oleh')::text AS didampingi_oleh
			FROM 
			form_ms f
		LEFT JOIN 
			document_ms d ON f.document_id = d.document_id
		LEFT JOIN 
			project_ms p ON f.project_id = p.project_id
			WHERE
			d.document_code = 'BA' 
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Slice to hold all form data
	var forms []models.FormsBA

	// Iterate through the rows
	for rows.Next() {
		// Scan the row into the Forms struct
		var form models.FormsBA
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
			&form.Judul,
			&form.Tanggal,
			&form.DokumenDA,
			&form.DokumenITCM,
			&form.DilakukanOleh,
			&form.DidampingiOleh,
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

func GetSpecBA(id string) (models.FormsBA, error) {
	var specBA models.FormsBA
	err := db.Get(&specBA, `SELECT 
	f.form_uuid, f.form_name, f.form_number, f.form_ticket, f.form_status,
	d.document_name,
	p.project_name,
	f.created_by, f.created_at, f.updated_by, f.updated_at, f.deleted_by, f.deleted_at,
	(f.form_data->>'judul')::text AS judul,
	(f.form_data->>'tanggal')::text AS tanggal,
	(f.form_data->>'dokumen_da')::text AS dokumen_da,
	(f.form_data->>'dokumen_itcm')::text AS dokumen_itcm,
	(f.form_data->>'dilakukan_oleh')::text AS dilakukan_oleh,
	(f.form_data->>'didampingi_oleh')::text AS didampingi_oleh
	FROM 
	form_ms f
LEFT JOIN 
	document_ms d ON f.document_id = d.document_id
LEFT JOIN 
	project_ms p ON f.project_id = p.project_id
	WHERE
	d.document_code = 'BA' 
	`, id)

	if err != nil {
		return models.FormsBA{}, err
	}

	return specBA, nil
}

func UpdateBA(updateBA models.Form, data models.BA, username string, userID int, isPublished bool, id string) (models.Form, error) {
	currentTime := time.Now()
	formStatus := "Draft"
	if isPublished {
		formStatus = "Published"
	}

	var projectID int64
	err := db.Get(&projectID, "SELECT project_id FROM project_ms WHERE project_uuid = $1", updateBA.ProjectUUID)
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
		"form_name":   updateBA.FormName,
		"form_ticket": updateBA.FormTicket,
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
	return updateBA, nil
}
