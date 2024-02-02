package main

import (
	"context"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/jmcvetta/neoism"
	"github.com/sirupsen/logrus"
	"log"
	"main/internal/app/store/graph"
	updater "main/internal/app/updater"
	"os"
	"os/signal"
	"path"
	"syscall"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", path.Join("configs", "updater_app.toml"), "path to config file")
}

func main() {
	config := updater.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go handleSignals(cancel)
	go func() {
		if _, err = Start(ctx, config); err != nil {
			log.Fatal(err)
			cancel()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			log.Println("Сервер успешно остановлен")
		}
	}
}

func handleSignals(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	cancel()
}

func Start(ctx context.Context, config *updater.Config) (*updater.Updater, error) {
	log := logrus.New()
	log.Println("New connection...")
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return nil, err
	}
	store := graph.NewGraphStore(db)
	log.Println("Service schedule updater started")
	updater, err := updater.NewUpdater(ctx, 600, log, store)
	updater.Run()
	if err != nil {
		return nil, err
	}
	return updater, nil
}

func newDB(dbURL string) (*neoism.Database, error) {
	db, err := neoism.Connect(dbURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}
