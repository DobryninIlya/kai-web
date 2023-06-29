package sqlstore

import (
	"main/internal/app/model"
	_ "main/internal/app/store"
)

// VerificationRepository ...
type VerificationRepository struct {
	store *Store
}

func (r VerificationRepository) GetPersonInfoScore(id int) (model.UserVerification, error) {
	var result model.UserVerification
	url := "SELECT id, date_update, faculty, course, group_id, idcard, groupname, studentid FROM user_verification WHERE id=$1"
	if err := r.store.db.QueryRow(
		url,
		id,
	).Scan(
		&result.Id,
		&result.DateUpdate,
		&result.Faculty,
		&result.Course,
		&result.GroupId,
		&result.Idcard,
		&result.Groupname,
		&result.Studentid,
	); err != nil {
		return model.UserVerification{}, err
	}

	return result, nil
}
