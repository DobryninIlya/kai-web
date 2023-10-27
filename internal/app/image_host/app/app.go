package image_host_app

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"main/internal/app/store/sqlstore"
)

func Start(ctx context.Context, config *Config) (*App, error) {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return nil, err
	}
	store := sqlstore.New(db)
	//migrations.MakeMigrations(db, logrus.New())
	srv := newApp(ctx, &store, config.BindAddr, *config)
	//return http.ListenAndServe(config.BindAddr, srv)
	return srv, nil
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
