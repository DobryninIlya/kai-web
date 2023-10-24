package migrations

import (
	"database/sql"
	"embed"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"path"
)

var (
	//go:embed migrations/*.sql
	embedMigrations embed.FS
)

func MakeMigrations(db *sql.DB, log *logrus.Logger) {
	//currentDir, _ := os.Getwd()
	dirMigrations := path.Join("migrations")
	const path = "intenral.app.migrations.migrations.go"
	goose.SetLogger(log)
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(
			logrus.WarnLevel,
			"Ошибка миграций: %v",
			err,
		)
	}
	options := []goose.OptionsFunc{}
	options = append(options, goose.WithNoColor(true))
	options = append(options, goose.WithAllowMissing())
	options = append(options, goose.WithNoVersioning())
	if err := goose.Up(db, dirMigrations, options...); err != nil {
		log.Logf(
			logrus.WarnLevel,
			"%v : Ошибка миграций: %v",
			path,
			err,
		)
		return
	}

	log.Logf(
		logrus.InfoLevel,
		"Миграции успешно применены",
	)
}
