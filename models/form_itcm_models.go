package models

import (
	"database/sql"
	"time"
)

type ITCM struct {
	NoDA          string `json:"no_da" validate:"required"`
	NamaPemohon   string `json:"nama_pemohon" validate:"required"`
	Intansi       string `json:"intansi" validate:"required"`
	Tanggal       string `json:"tanggal" validate:"required"`
	PerubahanAset string `json:"perubahan_aset" validate:"required"`
	Deskripsi     string `json:"deskripsi" validate:"required"`
}

type FormITCM struct {
	FormUUID      string         `json:"form_uuid" db:"form_uuid"`
	FormName      string         `json:"form_name" db:"form_name"`
	FormNumber    string         `json:"form_number" db:"form_number"`
	FormTicket    string         `json:"form_ticket" db:"form_ticket"`
	FormStatus    string         `json:"form_status" db:"form_status"`
	DocumentName  string         `json:"document_name" db:"document_name"`
	ProjectName   string         `json:"project_name" db:"project_name"`
	CreatedBy     string         `json:"created_by" db:"created_by"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedBy     sql.NullString `json:"updated_by" db:"updated_by"`
	UpdatedAt     sql.NullTime   `json:"updated_at" db:"updated_at"`
	DeletedBy     sql.NullString `json:"deleted_by" db:"deleted_by"`
	DeletedAt     sql.NullTime   `json:"deleted_at" db:"deleted_at"`
	NoDA          string         `json:"no_da" db:"no_da"`
	NamaPemohon   string         `json:"nama_pemohon" db:"nama_pemohon"`
	Intansi       string         `json:"intansi" db:"intansi"`
	Tanggal       string         `json:"tanggal" db:"tanggal"`
	PerubahanAset string         `json:"perubahan_aset" db:"perubahan_aset"`
	Deskripsi     string         `json:"deskripsi" db:"deskripsi"`
}
