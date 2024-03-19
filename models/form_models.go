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
	IsProject    bool           `json:"is_project" db:"is_project"`
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

// type FormName struct {
// 	UUID string `json:"document_uuid" db:"document_uuid"`
// 	Code string `json:"document_code" db:"document_code"`
// 	Name string `json:"document_name" db:"document_name"`
// }

// models.Signatory
type Signatory struct {
	Name     string `json:"name" db:"name"`
	Position string `json:"position" db:"position"`
	Role     string `json:"role_sign" db:"role_sign"`
}

type Signatories struct {
	UUID             string         `json:"sign_uuid" db:"sign_uuid"`
	Name             string         `json:"name" db:"name"`
	Position         string         `json:"position" db:"position"`
	Role             string         `json:"role_sign" db:"role_sign"`
	IsSign           bool           `json:"is_sign" db:"is_sign"`
	IsApprove        sql.NullBool   `json:"is_approve" db:"is_approve"`
	ApprovalStatus   string         `json:"approval_status"`
	ReasonNotApprove string         `json:"reason_not_approve" db:"reason_not_approve"`
	Created_by       sql.NullString `json:"created_by" db:"created_by"`
	Created_at       time.Time      `json:"created_at" db:"created_at"`
	Updated_by       sql.NullString `json:"updated_by" db:"updated_by"`
	Updated_at       sql.NullTime   `json:"updated_at" db:"updated_at"`
	Deleted_by       sql.NullString `json:"deleted_by" db:"deleted_by"`
	Deleted_at       sql.NullTime   `json:"deleted_at" db:"deleted_at"`
}

type Signatorie struct {
	UUID             string         `json:"sign_uuid" db:"sign_uuid"`
	Name             string         `json:"name" db:"name"`
	Position         string         `json:"position" db:"position"`
	Role             string         `json:"role_sign" db:"role_sign"`
	IsSign           bool           `json:"is_sign" db:"is_sign"`
	IsApprove        sql.NullBool   `json:"is_approve" db:"is_approve"`
	ReasonNotApprove string         `json:"reason_not_approve" db:"reason_not_approve"`
	Created_by       sql.NullString `json:"created_by" db:"created_by"`
	Created_at       time.Time      `json:"created_at" db:"created_at"`
	Updated_by       sql.NullString `json:"updated_by" db:"updated_by"`
	Updated_at       sql.NullTime   `json:"updated_at" db:"updated_at"`
	Deleted_by       sql.NullString `json:"deleted_by" db:"deleted_by"`
	Deleted_at       sql.NullTime   `json:"deleted_at" db:"deleted_at"`
}

type UpdateSign struct {
	IsSign           bool      `json:"is_sign" db:"is_sign" validate:"required"`
	IsApprove        bool      `json:"is_approve" db:"is_approve"`
	ReasonNotApprove string    `json:"reason_not_approve" db:"reason_not_approve"`
	Updated_by       string    `json:"updated_by" db:"updated_by"`
	Updated_at       time.Time `json:"updated_at" db:"updated_at"`
}

// models.Form
