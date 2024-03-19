package service

import (
	"document/models"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

func AddITCM(addForm models.Form, itcm models.ITCM, isPublished bool, userID int, username string, divisionCode string, recursionCount int, isProject bool, projectCode string, projectUUID string) error {
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
