package vk_app

import (
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	handler "main/internal/app/handlers"
	"main/internal/app/store/sqlstore"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var secretKey = os.Getenv("SECRET_KEY")

type App struct {
	router *chi.Mux
	done   chan os.Signal
	server *http.Server
	store  sqlstore.StoreInterface
	logger *logrus.Logger
}

func newApp(store sqlstore.StoreInterface, bindAddr string) *App {
	router := chi.NewRouter()
	server := &http.Server{
		Addr:    bindAddr,
		Handler: router,
	}
	a := &App{
		router: router,
		done:   make(chan os.Signal, 1),
		server: server,
		store:  store,
		logger: logrus.New(),
	}
	a.configureRouter()
	signal.Notify(a.done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return a
}

//func (a *App) run() {
//	a.configureRouter()
//	go func() {
//		log.Println("Starting worker")
//		log.Fatal(http.ListenAndServe(":8282", a.router))
//		//log.Fatal(a.server.ListenAndServeTLS("", ""))
//
//	}()
//	<-a.done
//	log.Println("Exiting")
//}

func (a *App) configureRouter() {
	a.router.Use(a.logRequest)
	a.router.Route("/web", func(r chi.Router) {
		r.Use(a.checkSign)
		r.Get("/", handler.New(a.store))
		r.Post("/registration", handler.NewRegistrationHandler(a.store))
		r.Get("/verification", handler.NewVerificationTemplate())
		r.Post("/verification/done", handler.NewVerificationDoneTemplate(a.store))
		r.Get("/get_lesson/{uId}", handler.NewLessonsHandler(a.store))
		r.Get("/exam", handler.NewExamHandler(a.store))
		r.Get("/teacher", handler.NewTeacherHandler(a.store))
		r.Get("/scoretable", handler.NewScoreListHandler(a.store))
		r.Get("/stylesheet", handler.NewStyleSheetHandler())
		//
		r.Route("/attestation", func(r chi.Router) {
			r.Get("/get_groups", handler.NewGroupsHandler())
			r.Get("/get_person", handler.NewPersonHandler())
			r.Get("/get_fac", handler.NewFacHandler())
			r.Get("/get_score", handler.NewScoreHandler())
		})

	})
	//a.router.Post("/brs.php", tools.NewExamHandler(a.store))
	//a.router.Handle("/static/css/*", http.StripPrefix("/static/css/", http.FileServer(http.Dir(filepath.Join("internal", "templates", "css")))))
	a.router.Handle("/static/css/*", http.StripPrefix("/static/css/", cssHandler(http.FileServer(http.Dir(filepath.Join("internal", "app", "templates", "css"))))))
	a.router.Handle("/static/js/*", http.StripPrefix("/static/js/", http.FileServer(http.Dir(filepath.Join("internal", "app", "templates", "js")))))
	a.router.Handle("/static/img/*", http.StripPrefix("/static/img/", http.FileServer(http.Dir(filepath.Join("internal", "app", "templates", "img")))))
	a.router.Handle("/static/json/*", http.StripPrefix("/static/json/", http.FileServer(http.Dir(filepath.Join("internal", "app", "templates", "json")))))
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func cssHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		next.ServeHTTP(w, r)
	})
}

func (a *App) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := a.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}
		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (a *App) checkSign(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := a.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
		})
		ok := handler.VerifyLaunchParams(r.RequestURI, secretKey)
		if ok != nil {
			logger.Log(
				logrus.WarnLevel,
				"the signature didn't match.",
			)
			return
		}
		rw := &responseWriter{w, http.StatusOK}
		h.ServeHTTP(rw, r)
		return
	})
}