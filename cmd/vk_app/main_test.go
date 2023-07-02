package main

import (
	"database/sql"
	"github.com/BurntSushi/toml"
	app "main/internal/app/vk_app"
	"os"
	"testing"
)

func TestDBConnection(t *testing.T) {
	config := app.NewConfig()
	_, err := toml.DecodeFile(configPathTest, config)
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	t.Logf("Successfully connected to database")
}

func TestMain(m *testing.M) {
	// Setup before running tests
	// e.g. initialize database connection

	// Run tests
	exitCode := m.Run()

	// Teardown after running tests
	// e.g. close database connection

	os.Exit(exitCode)
}
