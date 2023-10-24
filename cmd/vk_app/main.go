package main

import (
	"context"
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	app "main/internal/app/vk_app"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
)

var (
	configPath     string
	configPathTest string
)

func init() {
	flag.StringVar(&configPath, "config-path", path.Join("configs", "vk_app.toml"), "path to config file")
	flag.StringVar(&configPathTest, "config-path-test", path.Join("..", "..", "configs", "vk_app.toml"), "path to config file")
}

func main() {
	flag.Parse()
	config := app.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	var server *app.App
	var srv *http.Server
	go handleSignals(cancel)
	go func() {
		if server, err = app.Start(ctx, config); err != nil {
			log.Fatal(err)
			cancel()
		}
		srv = &http.Server{
			Addr:    config.BindAddr,
			Handler: server,
		}
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
			cancel()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			if err = srv.Shutdown(ctx); err != nil {
				log.Printf("Ошибка при остановке сервера: %v", err)
				if err = server.Close(); err != nil {
					log.Printf("Ошибка при закрытии сервера: %v", err)
				}
				return
			} else {
				if err = server.Close(); err != nil {
					log.Printf("Ошибка при закрытии сервера: %v", err)
				}
				log.Println("Сервер успешно остановлен")
				return
			}
		}
	}
}

func handleSignals(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	cancel()
}
