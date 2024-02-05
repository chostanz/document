package models

import (
	"database/sql"
	"time"
)

type Form struct {
	UUID         string         `json:"form_uuid" db:"form_uuid"`
	DocumentUUID string         `json:"document_uuid" db:"document_uuid" validate:"required"`
	UserID       string         `json:"user_id" db:"user_id" validate:"required"`
	FormNumber   string         `json:"form_number" db:"form_number" validate:"required"`
	FormTicket   string         `json:"form_ticket" db:"form_ticket" validate:"required"`
	FormStatus   string         `json:"form_status" db:"form_status"`
	Created_by   string         `json:"created_by" db:"created_by"`
	Created_at   time.Time      `json:"created_at" db:"created_at"`
	Updated_by   sql.NullString `json:"updated_by" db:"updated_by"`
	Updated_at   sql.NullTime   `json:"updated_at" db:"updated_at"`
	Deleted_by   sql.NullString `json:"deleted_by" db:"deleted_by"`
	Deleted_at   sql.NullTime   `json:"deleted_at" db:"deleted_at"`
}

type Forms struct {
	FormUUID   string         `json:"form_uuid" db:"form_uuid"`
	FormNumber string         `json:"form_number" db:"form_number"`
	FormTicket string         `json:"form_ticket" db:"form_ticket"`
	FormStatus string         `json:"form_status" db:"form_status"`
	DocumentID string         `json:"document_id" db:"document_id"`
	UserID     int64          `json:"user_id" db:"user_id"`
	Created_by string         `json:"created_by" db:"created_by"`
	Created_at time.Time      `json:"created_at" db:"created_at"`
	Updated_by sql.NullString `json:"updated_by" db:"updated_by"`
	Updated_at sql.NullTime   `json:"updated_at" db:"updated_at"`
	Deleted_by sql.NullString `json:"deleted_by" db:"deleted_by"`
	Deleted_at sql.NullTime   `json:"deleted_at" db:"deleted_at"`
}

type FormName struct {
	UUID string `json:"document_uuid" db:"document_uuid"`
	Code string `json:"document_code" db:"document_code"`
	Name string `json:"document_name" db:"document_name"`
}
