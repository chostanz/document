package service

import (
	"document/models"
	"log"
	"time"

	"github.com/google/uuid"
)

func AddForm(addFrom models.Form, isPublished bool, userUUID string) error {
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
	err := db.Get(&documentID, "SELECT document_id FROM document_ms WHERE document_uuid = $1", addFrom.DocumentUUID)
	if err != nil {
		log.Println("Error getting document_id:", err)
		return err
	}

	_, err = db.NamedExec("INSERT INTO form_ms (form_id, form_uuid, document_id, user_id, form_number, form_ticket, form_status, created_by) VALUES (:form_id, :form_uuid, :document_id, :user_id, :form_number, :form_ticket, :form_status, :created_by)", map[string]interface{}{
		"form_id":     app_id,
		"form_uuid":   uuidString,
		"document_id": documentID,
		"user_id":     addFrom.UserID,
		"form_number": addFrom.FormNumber,
		"form_ticket": addFrom.FormTicket,
		"form_status": formStatus,
		"created_by":  userUUID,
	})

	if err != nil {
		return err
	}
	return nil
}

func GetAllForm() ([]models.Forms, error) {

	form := []models.Forms{}

	//rows, errSelect := db.Queryx("SELECT f.form_uuid, f.form_number, f.form_ticket, f.form_status, f.user_id, f.created_by, f.created_at, f.updated_by, f.updated_at, d.document_name FROM form_ms f JOIN  document_ms d ON f.document_id = d.document_id WHERE f.deleted_at IS NULL")
	rows, errSelect := db.Queryx("select form_uuid, form_number, form_ticket, form_status, document_id, user_id, created_by, created_at, updated_by, updated_at from form_ms WHERE deleted_at IS NULL")
	if errSelect != nil {
		return nil, errSelect
	}

	for rows.Next() {
		place := models.Forms{}
		rows.StructScan(&place)
		form = append(form, place)
	}

	return form, nil
}
func ShowFormById(id string) (models.Forms, error) {
	var form models.Forms

	//err := db.Get(&form, "SELECT f.form_uuid, f.form_number, f.form_ticket, f.form_status, f.user_id, f.created_by, f.created_at, f.updated_by, f.updated_at, d.document_name FROM form_ms f JOIN  document_ms d ON f.document_id = d.document_id WHERE f.form_uuid = $1 AND f.deleted_at IS NULL", id)
	err := db.Get(&form, "select form_uuid, form_number, form_ticket, form_status, document_id, user_id, created_by, created_at, updated_by, updated_at from form_ms WHERE form_uuid = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return models.Forms{}, err
	}
	return form, nil

}

func UpdateForm(updateForm models.Form, id string, isPublished bool, userUUID string) (models.Form, error) {
	// username, errUser := GetUsernameByID(userUUID)
	// if errUser != nil {
	// 	log.Print(errUser)
	// 	return models.Document{}, errUser

	// }

	currentTime := time.Now()
	formStatus := "Draft"
	if isPublished {
		formStatus = "Published"
	}

	var documentID int64
	err := db.Get(&documentID, "SELECT document_id FROM document_ms WHERE document_uuid = $1", updateForm.DocumentUUID)
	if err != nil {
		log.Println("Error getting application_id:", err)
		return models.Form{}, err
	}
	_, err = db.NamedExec("UPDATE form_ms SET form_number = :form_number, form_ticket = :form_ticket, form_status = :form_status, document_id = :document_id, user_id = :user_id, updated_by = :updated_by, updated_at = :updated_at WHERE form_uuid = :id and form_status='Draft'", map[string]interface{}{
		"form_number": updateForm.FormNumber,
		"form_ticket": updateForm.FormTicket,
		"form_status": formStatus,
		"document_id": documentID,
		"user_id":     userUUID,
		"updated_by":  updateForm.Updated_by,
		"updated_at":  currentTime,
		"id":          id,
	})
	if err != nil {
		log.Print(err)
		return models.Form{}, err
	}
	return updateForm, nil
}
