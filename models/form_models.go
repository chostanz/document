package models

import (
	"database/sql"
	"time"
)

type Form struct {
	UUID         string         `json:"form_uuid" db:"form_uuid"`
	DocumentUUID string         `json:"document_uuid" db:"document_uuid"`
	DocumentID   int64          `json:"document_id" db:"document_id"`
	UserID       int            `json:"user_id" db:"user_id" validate:"required"`
	FormName     string         `json:"form_name" db:"form_name" validate:"required"`
	FormNumber   string         `json:"form_number" db:"form_number"`
	FormTicket   string         `json:"form_ticket" db:"form_ticket" validate:"required"`
	FormStatus   string         `json:"form_status" db:"form_status"`
	Created_by   string         `json:"created_by" db:"created_by"`
	ProjectUUID  string         `json:"project_uuid" db:"project_uuid"`
	ProjectID    int64          `json:"project_id" db:"project_id"`
	Created_at   time.Time      `json:"created_at" db:"created_at"`
	Updated_by   sql.NullString `json:"updated_by" db:"updated_by"`
	Updated_at   sql.NullTime   `json:"updated_at" db:"updated_at"`
	Deleted_by   sql.NullString `json:"deleted_by" db:"deleted_by"`
	Deleted_at   sql.NullTime   `json:"deleted_at" db:"deleted_at"`
}

type Forms struct {
	FormUUID     string         `json:"form_uuid" db:"form_uuid"`
	FormName     string         `json:"form_name" db:"form_name"`
	FormNumber   string         `json:"form_number" db:"form_number"`
	FormTicket   string         `json:"form_ticket" db:"form_ticket"`
	FormStatus   string         `json:"form_status" db:"form_status"`
	DocumentName string         `json:"document_name" db:"document_name"`
	Created_by   string         `json:"created_by" db:"created_by"`
	Created_at   time.Time      `json:"created_at" db:"created_at"`
	Updated_by   sql.NullString `json:"updated_by" db:"updated_by"`
	Updated_at   sql.NullTime   `json:"updated_at" db:"updated_at"`
	Deleted_by   sql.NullString `json:"deleted_by" db:"deleted_by"`
	Deleted_at   sql.NullTime   `json:"deleted_at" db:"deleted_at"`
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

type Formss struct {
	FormUUID                             string         `json:"form_uuid" db:"form_uuid"`
	FormName                             string         `json:"form_name" db:"form_name"`
	FormNumber                           string         `json:"form_number" db:"form_number"`
	FormTicket                           string         `json:"form_ticket" db:"form_ticket"`
	FormStatus                           string         `json:"form_status" db:"form_status"`
	DocumentName                         string         `json:"document_name" db:"document_name"`
	ProjectName                          string         `json:"project_name" db:"project_name"`
	IsApprove                            sql.NullBool   `json:"is_approve" db:"is_approve"`
	ApprovalStatus                       string         `json:"approval_status"`
	Reason                               sql.NullString `json:"reason" db:"reason"`
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
