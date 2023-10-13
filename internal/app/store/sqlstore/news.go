package sqlstore

import (
	"main/internal/app/model"
	_ "main/internal/app/store"
)

// NewsRepository реализует работу раздела новостей с хранилищем базы данных
type NewsRepository struct {
	store *Store
}

func (r *NewsRepository) GetNews(limit int) ([]*model.News, error) {
	if limit == 0 {
		limit = 10
	}
	rows, err := r.store.db.Query("SELECT id, header, description, body, date FROM news ORDER BY id DESC LIMIT $1",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	resultList := make([]*model.News, 0)
	for rows.Next() {
		var news model.News
		err := rows.Scan(&news.Id, &news.Header, &news.Description, &news.Body, &news.Date)
		if err != nil {
			return nil, err
		}
		resultList = append(resultList, &news)
	}
	return resultList, nil
}

func (r *NewsRepository) CreateNews() {

}
