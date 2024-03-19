package models

type ITCM struct {
	PerubahanAset string `json:"perubahan_aset" validate:"required"`
	Deskripsi     string `json:"deskripsi" validate:"required"`
}
