package models

import (
	"database/sql"
	"time"
)

type BA struct {
	Judul          string `json:"judul"`
	Tanggal        string `json:"tanggal"`
	DokumenDA      string `json:"dokumen_da"`
	DokumenITCM    string `json:"dokumen_itcm"`
	DilakukanOleh  string `json:"dilakukan_oleh"`
	DidampingiOleh string `json:"didampingi_oleh"`
}

type FormsBA struct {
	FormUUID         string         `json:"form_uuid" db:"form_uuid"`
	FormName         string         `json:"form_name" db:"form_name"`
	FormNumber       string         `json:"form_number" db:"form_number"`
	FormTicket       string         `json:"form_ticket" db:"form_ticket"`
	FormStatus       string         `json:"form_status" db:"form_status"`
	DocumentName     string         `json:"document_name" db:"document_name"`
	ProjectName      string         `json:"project_name" db:"project_name"`
	IsApprove        sql.NullBool   `json:"is_approve" db:"is_approve"`
	ApprovalStatus   string         `json:"approval_status"`
	Reason sql.NullString `json:"reason" db:"reason"`
	CreatedBy        string         `json:"created_by" db:"created_by"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedBy        sql.NullString `json:"updated_by" db:"updated_by"`
	UpdatedAt        sql.NullTime   `json:"updated_at" db:"updated_at"`
	DeletedBy        sql.NullString `json:"deleted_by" db:"deleted_by"`
	DeletedAt        sql.NullTime   `json:"deleted_at" db:"deleted_at"`
	Judul            string         `json:"judul" db:"judul"`
	Tanggal          string         `json:"tanggal" db:"tanggal"`
	DokumenDA        string         `json:"dokumen_da" db:"dokumen_da"`
	DokumenITCM      string         `json:"dokumen_itcm" db:"dokumen_itcm"`
	DilakukanOleh    string         `json:"dilakukan_oleh" db:"dilakukan_oleh"`
	DidampingiOleh   string         `json:"didampingi_oleh" db:"didampingi_oleh"`
}
