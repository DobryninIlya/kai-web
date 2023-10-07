package vk_app

import (
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	handler "main/internal/app/handlers"
	api "main/internal/app/handlers/api"
	"main/internal/app/store/sqlstore"
	"main/internal/app/tg_api"
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
	tgApi  *tg_api.APItg
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
		tgApi:  tg_api.NewAPItg(),
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
	//a.router.Use(imageStatusCodeHandler)
	a.router.Route("/api", func(r chi.Router) {
		r.Route("/schedule", func(r chi.Router) {
			r.Use(a.authorizationByToken)
			r.Get("/{groupid}", api.NewScheduleHandler(a.store, a.logger))          // Расписание полностью
			r.Get("/{groupid}/by_margin", api.NewLessonsHandler(a.store, a.logger)) // На день с отступом margin от текущего дня
			r.Get("/{groupid}/teachers", api.NewTeachersHandler(a.store, a.logger)) // На день с отступом margin от текущего дня
		})
		r.Route("/groups", func(r chi.Router) {
			r.Use(a.authorizationByToken)
			r.Get("/{group}", api.NewIdByGroupHandler(a.store, a.logger)) // ID группы по ее номеру
		})
		//r.Get("/", api.NewLessonsHandler(a.store, a.logger))
		r.Route("/feedback", func(r chi.Router) {
			r.Use(a.authorizationByToken)
			r.Post("/", api.NewFeedbackHandler(a.store, a.logger, a.tgApi)) // ID группы по ее номеру
		}) // Главная страница документации
		r.Get("/doc", api.NewDocumentationPageHandler())                // Главная страница документации
		r.Get("/doc/{page}", api.NewDocumentationOtherPageHandler())    // Страница документации
		r.Get("/token", api.NewRegistrationHandler(a.store, a.logger))  // ID группы по ее номеру
		r.Get("/token/whoiam", api.NewWhoIAmHandler(a.store, a.logger)) // ID группы по ее номеру
	})
	a.router.Route("/web", func(r chi.Router) {
		r.Use(a.checkSign)
		r.Get("/", handler.New(a.store, a.logger))
		r.Post("/registration", handler.NewRegistrationHandler(a.store, a.logger))
		r.Get("/delete_user", handler.NewDeleteUserHandler(a.store, a.logger))
		r.Get("/verification", handler.NewVerificationTemplate())
		r.Post("/verification/done", handler.NewVerificationDoneTemplate(a.store, a.logger))
		r.Get("/get_lesson/{uId}", handler.NewLessonsHandler(a.store, a.logger))
		r.Get("/exam", handler.NewExamHandler(a.store, a.logger))
		r.Get("/teacher", handler.NewTeacherHandler(a.store, a.logger))
		r.Get("/scoretable", handler.NewScoreListHandler(a.store, a.logger))
		r.Post("/delete_lesson", handler.NewDeleteLessonHandler(a.store, a.logger))
		r.Post("/return_lesson", handler.NewReturnLessonHandler(a.store, a.logger))
		r.Get("/stylesheet", handler.NewStyleSheetHandler())
		//
		r.Route("/attestation", func(r chi.Router) {
			r.Get("/get_groups", handler.NewGroupsHandler(a.logger))
			r.Get("/get_person", handler.NewPersonHandler(a.logger))
			r.Get("/get_fac", handler.NewFacHandler(a.logger))
			r.Get("/get_score", handler.NewScoreHandler(a.logger))
		})

	})
	a.router.Handle("/static/css/*", http.StripPrefix("/static/css/", cssHandler(http.FileServer(http.Dir(filepath.Join("internal", "app", "templates", "css"))))))
	a.router.Handle("/static/js/*", http.StripPrefix("/static/js/", http.FileServer(http.Dir(filepath.Join("internal", "app", "templates", "js")))))
	a.router.Handle("/static/img/*", http.StripPrefix("/static/img/", http.FileServer(http.Dir(filepath.Join("internal", "app", "templates", "img")))))
	a.router.Handle("/static/json/*", http.StripPrefix("/static/json/", http.FileServer(http.Dir(filepath.Join("internal", "app", "templates", "json")))))
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func (a *App) authorizationByToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query()
		_, err, code := a.store.API().CheckToken(url.Get("token"))
		if err != nil {
			handler.ErrorHandlerAPI(w, r, code, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

//func imageStatusCodeHandler(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		rec := &responseWriter{w, http.StatusOK}
//
//		next.ServeHTTP(rec, r)
//		if rec.code >= 400 && rec.code < 500 {
//			fmt.Sprintf("<div><img src=\"https://http.cat/%d\"></div>", rec.code)
//			w.Header().Set("Content-Type", "text/html; charset=utf-8")
//			w.WriteHeader(rec.code)
//		}
//	})
//}

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
			w.WriteHeader(http.StatusForbidden)
			return
		}
		rw := &responseWriter{w, http.StatusOK}
		h.ServeHTTP(rw, r)

		return
	})
}
