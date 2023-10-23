package vk_app

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"main/internal/app/firebase"
	api "main/internal/app/handlers/api"
	"main/internal/app/handlers/web_app"
	"main/internal/app/mailer"
	"main/internal/app/openai"
	"main/internal/app/store/influxdb"
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
	router      *chi.Mux
	done        chan os.Signal
	server      *http.Server
	store       sqlstore.StoreInterface
	tgApi       *tg_api.APItg
	logger      *logrus.Logger
	mailer      *mailer.Mailing
	firebaseAPI *firebase.FirebaseAPI
	ctx         context.Context
	weekParity  int
	metrics     *influxdb.Metrics
	openai      *openai.ChatGPT
}

func newApp(store sqlstore.StoreInterface, bindAddr string, weekParity int, firebaseAPI *firebase.FirebaseAPI, config Config) *App {
	router := chi.NewRouter()
	server := &http.Server{
		Addr:    bindAddr,
		Handler: router,
	}
	logger := logrus.New()
	ctx := context.Background()
	a := &App{
		router:      router,
		done:        make(chan os.Signal, 1),
		server:      server,
		store:       store,
		logger:      logger,
		tgApi:       tg_api.NewAPItg(),
		mailer:      mailer.NewMailing(store, logger),
		weekParity:  weekParity,
		firebaseAPI: firebaseAPI,
		metrics:     influxdb.NewMetrics(ctx, config.InfluxDBToken, config.InfluxDBURL, config.InfluxDBName, logger),
		ctx:         ctx,
		openai:      openai.NewChatGPT(ctx, logger, "gpt-3.5-turbo", 0.7, "user"),
	}
	a.openai.WithPrompt(openai.NEWS_PROMPT)
	a.configureRouter()
	signal.Notify(a.done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return a
}

func (a *App) configureRouter() {
	a.router.Use(a.logRequest)
	//a.router.Use(a.APIMetricsMiddleware)
	//a.router.Use(imageStatusCodeHandler)
	a.router.Route("/api", func(r chi.Router) {
		r.Route("/private", func(r chi.Router) {
			r.Route("/sendMessage", func(r chi.Router) {
				r.Use(a.authorizationBySecretPhrase)
				r.Get("/vk", api.NewSendMailVKHandler(a.store, a.logger, a.mailer)) // Отправка сообщения в ВК
			})
		})
		r.Get("/week", api.NewWeekParityHandler(a.weekParity)) // Текущая четность недели
		r.Route("/schedule", func(r chi.Router) {
			r.Use(a.authorizationByToken)
			r.Get("/{groupid}", api.NewScheduleHandler(a.store, a.logger))                        // Расписание полностью
			r.Get("/{groupid}/by_margin", api.NewLessonsHandler(a.store, a.logger, a.weekParity)) // На день с отступом margin от текущего дня
			r.Get("/{groupid}/teachers", api.NewTeachersHandler(a.store, a.logger))               // Список преподавателей
		})
		r.Route("/groups", func(r chi.Router) {
			r.Use(a.authorizationByToken)
			r.Get("/{group}", api.NewIdByGroupHandler(a.store, a.logger)) // ID группы по ее номеру
		})
		//r.Get("/", api.NewLessonsHandler(a.store, a.logger))
		r.Route("/feedback", func(r chi.Router) {
			r.Use(a.authorizationByToken)
			r.Post("/", api.NewFeedbackHandler(a.store, a.logger, a.tgApi)) // Отправка отзыва в чат
		})
		r.Route("/attestation", func(r chi.Router) {
			r.Use(a.authorizationByToken)
			r.Get("/", api.NewScoreHandler(a.logger))        // Список оценок БРС
			r.Get("/faculties", api.NewFacHandler(a.logger)) // Список факультетов
			r.Get("/groups", api.NewGroupsHandler(a.logger)) // Список групп
			r.Get("/person", api.NewPersonHandler(a.logger)) // Список фио
		})
		r.Route("/news", func(r chi.Router) {
			r.Post("/", api.NewMakeNewsHandler(a.store, a.logger))
			r.Get("/", api.NewNewsPageHandler(a.store, a.logger))
			r.Get("/{newsId}", api.NewNewsHandler(a.store, a.logger))
			r.Get("/previews", api.NewNewsPreviewsHandler(a.store, a.logger))
			r.Get("/create", api.NewNewsCreateFormHandler(a.store, a.logger))
			r.Post("/collect", api.NewHandleVKUpdateHandler(a.store, a.logger, a.openai))
			r.Get("/collect", api.NewHandleVKUpdateHandler(a.store, a.logger, a.openai))
		})

		r.Get("/doc", api.NewDocumentationPageHandler())                                     // Главная страница документации
		r.Get("/doc/{page}", api.NewDocumentationOtherPageHandler())                         // Страница документации
		r.Get("/token", api.NewRegistrationHandler(a.ctx, a.store, a.logger, a.firebaseAPI)) // Получение токена
		r.Get("/token/whoiam", api.NewWhoIAmHandler(a.store, a.logger))                      // Информация из токена
		r.Get("/confirmation_code/{code}", api.NewSetCodeHandler(&a.store))                  // Установка токена верификации до перезагрузки
		r.Get("/echo", api.NewEchoHandler())                                                 // Эхо  метод для теста
	})
	a.router.Route("/web", func(r chi.Router) {
		r.Use(a.checkSign)
		r.Get("/", web_app.New(a.store, a.logger))
		r.Post("/registration", web_app.NewRegistrationHandler(a.store, a.logger))
		r.Get("/delete_user", web_app.NewDeleteUserHandler(a.store, a.logger))
		r.Get("/verification", web_app.NewVerificationTemplate())
		r.Post("/verification/done", web_app.NewVerificationDoneTemplate(a.store, a.logger))
		r.Get("/get_lesson/{uId}", web_app.NewLessonsHandler(a.store, a.logger, a.weekParity))
		r.Get("/exam", web_app.NewExamHandler(a.store, a.logger))
		r.Get("/teacher", web_app.NewTeacherHandler(a.store, a.logger))
		r.Get("/scoretable", web_app.NewScoreListHandler(a.store, a.logger))
		r.Post("/delete_lesson", web_app.NewDeleteLessonHandler(a.store, a.logger))
		r.Post("/return_lesson", web_app.NewReturnLessonHandler(a.store, a.logger))
		r.Get("/stylesheet", web_app.NewStyleSheetHandler())
		//
		r.Route("/attestation", func(r chi.Router) {
			r.Get("/get_groups", web_app.NewGroupsHandler(a.logger))
			r.Get("/get_person", web_app.NewPersonHandler(a.logger))
			r.Get("/get_fac", web_app.NewFacHandler(a.logger))
			r.Get("/get_score", web_app.NewScoreHandler(a.logger))

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
			web_app.ErrorHandlerAPI(w, r, code, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *App) authorizationBySecretPhrase(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query()
		_, err, code := a.store.API().CheckSecret(url.Get("secret"))
		if err != nil {
			web_app.ErrorHandlerAPI(w, r, code, err)
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
		ok := web_app.VerifyLaunchParams(r.RequestURI, secretKey)
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

func (a *App) APIMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем имя вызываемого метода из URL
		a.metrics.IncrementAPI("test")

		// Вызываем следующий обработчик
		next.ServeHTTP(w, r)
	})
}
