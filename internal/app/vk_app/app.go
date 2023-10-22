package vk_app

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"log"
	"main/internal/app/firebase"
	"main/internal/app/migrations"
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
	migrations.MakeMigrations(db, logrus.New())
	firebaseAPI, err := firebase.NewFirebaseAPI(config.FirebaseServiceAccountKeyPath, config.FirebaseProjectID)
	if err != nil {
		log.Fatalf("Ошибка инициализации Firebase API: %v. Проверьте, находится ли serviceAccountKey.json в папке configs.", err.Error())
	}
	srv := newApp(&store, config.BindAddr, config.Chetnost, firebaseAPI, *config)
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
