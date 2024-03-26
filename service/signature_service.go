package service

import (
	"document/models"
	"log"
	"time"
)

func GetAllPersonalName() ([]models.Personal, error) {
	getUserAppRole := []models.Personal{}

	// Lakukan query ke database lain
	rows, err := db2.Queryx("SELECT u.user_id, pd.personal_name FROM user_ms u JOIN personal_data_ms pd ON u.user_id = pd.user_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Proses data hasil query
	for rows.Next() {
		place := models.Personal{}
		err := rows.StructScan(&place)
		if err != nil {
			log.Println("Error scanning row to struct:", err)
			continue
		}
		getUserAppRole = append(getUserAppRole, place)
	}

	return getUserAppRole, nil
}

func GetSignatureForm(id string) ([]models.Signatories, error) {
	var signatories []models.Signatories

	err := db.Select(&signatories, `SELECT 
	sf.sign_uuid, 
	sf.name, 
	sf.position, 
	sf.role_sign, 
	sf.is_sign, 
	sf.created_by, 
	sf.created_at, 
	sf.updated_by, 
	sf.updated_at, 
	sf.deleted_by, 
	sf.deleted_at
FROM 
	sign_form sf 
	JOIN form_ms fm ON sf.form_id = fm.form_id 
WHERE 
	fm.form_uuid = $1`, id)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return signatories, nil

}

func GetSpecSignatureByID(id string) (models.Signatorie, error) {
	var signatories models.Signatorie
	err := db.Get(&signatories, "SELECT sign_uuid, name, position, role_sign, is_sign, created_by, created_at, updated_by, updated_at, deleted_by, deleted_at FROM sign_form sf WHERE sign_uuid = $1", id)
	if err != nil {
		log.Print(err)
		return models.Signatorie{}, err
	}

	return signatories, nil
}

func UpdateFormSignature(updateSign models.UpdateSign, id string, username string) error {
	currentTime := time.Now()

	_, err := db.NamedExec("UPDATE sign_form SET is_sign = :is_sign, updated_by = :updated_by, updated_at = :updated_at WHERE sign_uuid = :id", map[string]interface{}{
		"is_sign":    updateSign.IsSign,
		"updated_by": username,
		"updated_at": currentTime,
		"id":         id,
	})
	if err != nil {
		return err
	}
	return nil
}

func AddApproval(addApproval models.AddApproval, id string, username string) error {
	currentTime := time.Now()

	_, err := db.NamedExec("UPDATE form_ms SET is_approve = :is_approve, reason = :reason, updated_by = :updated_by, updated_at = :updated_at WHERE form_uuid = :id", map[string]interface{}{
		"is_approve": addApproval.IsApproval,
		"reason":     addApproval.Reason,
		"updated_by": username,
		"updated_at": currentTime,
		"id":         id,
	})
	if err != nil {
		return err
	}
	return nil
}
