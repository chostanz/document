package models

type DampakAnalisa struct {
	//NamaProyek                           string    `json:"nama_proyek"`
	NamaAnalis                           string `json:"nama_analis"`
	Jabatan                              string `json:"jabatan"`
	Departemen                           string `json:"departemen"`
	JenisPerubahan                       string `json:"jenis_perubahan"`
	DetailDampakPerubahan                string `json:"detail_dampak_perubahan"`
	RencanaPengembanganPerubahan         string `json:"rencana_pengembangan_perubahan"`
	RencanaPengujianPerubahanSistem      string `json:"rencana_pengujian_perubahan_sistem"`
	RencanaRilisPerubahanDanImplementasi string `json:"rencana_rilis_perubahan_dan_implementasi"`
}
