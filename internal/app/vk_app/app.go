package vk_app

import (
	"database/sql"
	_ "github.com/lib/pq"
	"main/internal/app/store/sqlstore"
	"net/http"
)

type APIServer struct {
	store *sqlstore.StoreInterface
}

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	srv := newApp(&store, config.BindAddr)

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
