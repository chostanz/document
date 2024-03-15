package service

import (
	"document/models"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

type ITCM struct {
	PerubahanAset    string   `json:"perubahan_aset" validate:"required"`
	Deskripsi        string   `json:"deskripsi" validate:"required"`
	isApproved       bool     `json:"is_approved"`
	AlasanNotApprove string   `json:"alasan_not_approve,omitempty"`
	Pemohon          string   `json:"pemohon"`
	Penerima         string   `json:"penerima"`
	AtasanPemohon    []string `json:"atasan_pemohon"`
	AtasanPenerima   string   `json:"atasan_penerima"`
}

// type ITCM struct {
// 	FormName         string `json:"form_name"`
// 	FormTicket       string `json:"form_ticket"`
// 	DocumentUUID     string `json:"document_uuid"`
// 	IsPublished      bool   `json:"isPublished"`
// 	PerubahanAset    string `json:"perubahan_aset"`
// 	Deskripsi        string `json:"deskripsi"`
// 	Status           bool   `json:"status"`                       // Approve: true, Not approve: false
// 	AlasanNotApprove string `json:"alasan_not_approve,omitempty"` // Alasan ketika tidak di-approve
// 	Pemohon          string `json:"pemohon"`                      // Dapat diambil dari token
// 	//TTDPemohon       bool   `json:"ttd_pemohon"`
// 	Penerima string `json:"penerima"`
// 	//	TTDPenerima       bool     `json:"ttd_penerima"`
// 	AtasanPemohon []string `json:"atasan_pemohon"` // List dari atasan pemohon
// 	//	TTDatasanPemohon  bool     `json:"ttd_atasan_pemohon"`
// 	AtasanPenerima string `json:"atasan_penerima"`
// 	//TTDatasanPenerima bool     `json:"ttd_atasan_penerima"`
// }

func AddITCM(addForm models.Form, itcm ITCM, isPublished bool, userID int, username string, divisionCode string, recursionCount int) error {
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

	var documentCode string
	err = db.Get(&documentCode, "SELECT document_code FROM document_ms WHERE document_uuid = $1", addForm.DocumentUUID)
	if err != nil {
		log.Println("Error getting document code:", err)
		return err
	}
	// Generate form number based on document code
	formNumber, err := generateFormNumber(documentID, divisionCode, recursionCount+1)
	if err != nil {
		log.Println("Error generating form number:", err)
		return err
	}

	// Marshal ITCM struct to JSON
	itcmJSON, err := json.Marshal(itcm)
	if err != nil {
		log.Println("Error marshaling ITCM struct:", err)
		return err
	}

	//	log.Println("ITCM:", itcm)
	_, err = db.NamedExec("INSERT INTO form_ms (form_id, form_uuid, document_id, user_id, form_name, form_number, form_ticket, form_status, form_data, created_by) VALUES (:form_id, :form_uuid, :document_id, :user_id, :form_name, :form_number, :form_ticket, :form_status, :form_data, :created_by)", map[string]interface{}{
		"form_id":     appID,
		"form_uuid":   uuidString,
		"document_id": documentID,
		"user_id":     userID,
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
