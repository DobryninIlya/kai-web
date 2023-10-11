package sqlstore

type MailingRepository struct {
	store *Store
}

func (r MailingRepository) GetVKRecipients() ([]int64, error) {
	vkIds := make([]int64, 0)
	rows, err := r.store.db.Query("SELECT id_vk FROM users WHERE id_vk IS NOT NULL")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		vkIds = append(vkIds, id)
	}
	return vkIds, nil
}
