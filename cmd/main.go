package main

import (
	"crypto/tls"
	"github.com/go-chi/chi"
	"log"
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
	//pgConn := pgconn.Config{
	//	Host:     "localhost",
	//	Port:     5432,
	//	Database: "",
	//	User:     "",
	//	Password: "",
	//}
	//connConfig := pgx.ConnConfig{Config: pgConn}
	//dbService := service.NewService(context.Background(), pgxpool.Config{
	//	ConnConfig:            &connConfig,
	//	BeforeConnect:         nil,
	//	AfterConnect:          nil,
	//	BeforeAcquire:         nil,
	//	AfterRelease:          nil,
	//	MaxConnLifetime:       0,
	//	MaxConnLifetimeJitter: 0,
	//	MaxConnIdleTime:       0,
	//	MaxConns:              0,
	//	MinConns:              0,
	//	HealthCheckPeriod:     0,
	//	LazyConnect:           false,
	//})

	//a.router.Route("/web", func(r chi.Router) {
	//	r.Get("/", tools.New(dbService))
	//	r.Post("/registration", tools.NewRegistrationHandler(dbService))
	//	r.Get("/verification", tools.NewVerificationTemplate(dbService))
	//	r.Post("/verification/done", tools.NewVerificationDoneTemplate(dbService))
	//	r.Get("/get_lesson/{uId}", tools.NewLessonsHandler(dbService))
	//	r.Get("/exam", tools.NewExamHandler(dbService))
	//	r.Get("/teacher", tools.NewTeacherHandler(dbService))
	//	r.Get("/scoretable", tools.NewScoreListHandler(dbService))
	//	r.Get("/stylesheet", tools.NewStyleSheetHandler())
	//	r.Get("/cas", tools.NewCas(dbService))
	//
	//	r.Route("/attestation", func(r chi.Router) {
	//		r.Get("/get_groups", tools.NewGroupsHandler(dbService))
	//		r.Get("/get_person", tools.NewPersonHandler(dbService))
	//		r.Get("/get_fac", tools.NewFacHandler(dbService))
	//		r.Get("/get_score", tools.NewScoreHandler(dbService))
	//	})
	//
	//})
	//a.router.Post("/brs.php", tools.NewExamHandler(dbService))
	////a.router.Handle("/static/css/*", http.StripPrefix("/static/css/", http.FileServer(http.Dir(filepath.Join("internal", "templates", "css")))))
	//a.router.Handle("/static/css/*", http.StripPrefix("/static/css/", cssHandler(http.FileServer(http.Dir(filepath.Join("internal", "templates", "css"))))))
	//a.router.Handle("/static/js/*", http.StripPrefix("/static/js/", http.FileServer(http.Dir(filepath.Join("internal", "templates", "js")))))
	//a.router.Handle("/static/img/*", http.StripPrefix("/static/img/", http.FileServer(http.Dir(filepath.Join("internal", "templates", "img")))))
	//go func() {
	//	log.Println("Starting worker")
	//	log.Fatal(http.ListenAndServe(":8282", a.router))
	//	//log.Fatal(a.server.ListenAndServeTLS("", ""))
	//
	//}()
	//<-a.done
	//log.Println("Exiting")
}

func cssHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		next.ServeHTTP(w, r)
	})
}

func main() {
	var app = NewApp()
	app.run()
}
