package main

import (
	"github.com/go-chi/chi"
	"github.com/jackc/pgx"
	"log"
	service "main/internal/database"
	web "main/internal/web_app"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

type App struct {
	router *chi.Mux
	done   chan os.Signal
	//store  map[int]city.City
}

func NewApp() *App {
	ret := &App{
		router: chi.NewRouter(),
		done:   make(chan os.Signal, 1),
	}
	signal.Notify(ret.done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return ret
}

func (a *App) run() {
	dbService := service.NewService(pgx.ConnConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,
		Database: os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	a.router.Route("/web", func(r chi.Router) {
		r.Get("/", web.New(dbService))
		r.Get("/get_lesson/{uId}", web.NewLessonsHandler(dbService))
		r.Get("/stylesheet", web.NewStyleSheetHandler())

	})
	a.router.Handle("/static/css/*", http.StripPrefix("/static/css/", http.FileServer(http.Dir(filepath.Join("internal", "templates", "css")))))
	a.router.Handle("/static/js/*", http.StripPrefix("/static/js/", http.FileServer(http.Dir(filepath.Join("internal", "templates", "js")))))
	go func() {
		log.Println("Starting worker")
		log.Fatal(http.ListenAndServe(":8282", a.router))

	}()
	<-a.done
	log.Println("Exiting")
}
func main() {
	var app = NewApp()
	app.run()
}
