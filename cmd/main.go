package main

import (
	"context"
	"crypto/tls"
	"github.com/go-chi/chi"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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
	server *http.Server
	//store  map[int]city.City
}

func NewApp() *App {
	cert, err := tls.LoadX509KeyPair(filepath.Join("openssl", "cert.pem"), filepath.Join("openssl", "key.pem"))
	if err != nil {
		log.Fatal(err)
	}
	router := chi.NewRouter()
	server := &http.Server{
		Addr:    ":8282",
		Handler: router,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}
	ret := &App{
		router: router,
		done:   make(chan os.Signal, 1),
		server: server,
	}
	signal.Notify(ret.done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return ret
}

func (a *App) run() {
	pgConn := pgconn.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,
		Database: os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	connConfig := pgx.ConnConfig{Config: pgConn}
	dbService := service.NewService(context.Background(), pgxpool.Config{
		ConnConfig:            &connConfig,
		BeforeConnect:         nil,
		AfterConnect:          nil,
		BeforeAcquire:         nil,
		AfterRelease:          nil,
		MaxConnLifetime:       0,
		MaxConnLifetimeJitter: 0,
		MaxConnIdleTime:       0,
		MaxConns:              0,
		MinConns:              0,
		HealthCheckPeriod:     0,
		LazyConnect:           false,
	})

	a.router.Route("/web", func(r chi.Router) {
		r.Get("/", web.New(dbService))
		r.Post("/registration", web.NewRegistrationHandler(dbService))
		r.Get("/get_lesson/{uId}", web.NewLessonsHandler(dbService))
		r.Get("/exam", web.NewExamHandler(dbService))
		r.Get("/teacher", web.NewTeacherHandler(dbService))
		r.Get("/stylesheet", web.NewStyleSheetHandler())
		r.Get("/cas", web.NewCas(dbService))

	})
	a.router.Handle("/static/css/*", http.StripPrefix("/static/css/", http.FileServer(http.Dir(filepath.Join("internal", "templates", "css")))))
	a.router.Handle("/static/js/*", http.StripPrefix("/static/js/", http.FileServer(http.Dir(filepath.Join("internal", "templates", "js")))))
	a.router.Handle("/static/img/*", http.StripPrefix("/static/img/", http.FileServer(http.Dir(filepath.Join("internal", "templates", "img")))))
	go func() {
		log.Println("Starting worker")
		log.Fatal(http.ListenAndServe(":8282", a.router))
		//log.Fatal(a.server.ListenAndServeTLS("", ""))

	}()
	<-a.done
	log.Println("Exiting")
}
func main() {
	var app = NewApp()
	app.run()
}
