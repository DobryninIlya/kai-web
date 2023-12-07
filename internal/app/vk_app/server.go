package vk_app

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"log"
	"main/internal/app/authorization"
	"main/internal/app/firebase"
	api "main/internal/app/handlers/api"
	"main/internal/app/handlers/api/auth"
	pay_handler "main/internal/app/handlers/api/payments"
	portal "main/internal/app/handlers/api/portal"
	"main/internal/app/handlers/api/schedule"
	task "main/internal/app/handlers/api/tasks"
	"main/internal/app/handlers/web_app"
	"main/internal/app/mailer"
	"main/internal/app/openai"
	"main/internal/app/store/influxdb"
	"main/internal/app/store/sqlstore"
	"main/internal/app/tg_api"
	"main/internal/app/tools"
	"main/internal/payments"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var secretKey = os.Getenv("SECRET_KEY")

type App struct {
	router      *chi.Mux
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
	auth        authorization.AuthorizationInterface
	pay         payments.Yokassa
}

func newApp(ctx context.Context, store sqlstore.StoreInterface, bindAddr string, weekParity int, firebaseAPI *firebase.FirebaseAPI, config Config) *App {
	router := chi.NewRouter()
	server := &http.Server{
		Addr:    bindAddr,
		Handler: router,
	}
	logger := logrus.New()
	a := &App{
		router:      router,
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
		auth:        authorization.NewAuthorization(logger),
		pay:         payments.NewYokassa(config.ShopID, config.APIKey, logger),
	}
	a.openai.WithPrompt(openai.NEWS_PROMPT)
	a.configureRouter()
	return a
}

func (a *App) Close() error {
	err := a.server.Close()
	if err != nil {
		return err
	}
	return a.server.Close()
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
		r.Post("/registration", auth.NewRegistrationByPasswordHandler(a.ctx, a.store, a.logger, a.auth)) // Регистрация по логину и паролю
		r.Route("/auth", func(r chi.Router) {
			r.Get("/personal", auth.NewAboutInfoHandler(a.ctx, a.store, a.logger, a.firebaseAPI, a.auth))              // Номер группы")
			r.Get("/profile_photo", auth.NewProfilePhotoURLHandler(a.ctx, a.store, a.logger, a.firebaseAPI, a.auth))   // Номер группы")
			r.Post("/profile_photo", auth.NewPostProfilePhotoHandler(a.ctx, a.store, a.logger, a.firebaseAPI, a.auth)) // Номер группы")
			r.Get("/group", auth.NewGroupInfoHandler(a.ctx, a.store, a.logger, a.firebaseAPI, a.auth))                 // Номер группы")
			r.Get("/attestation", auth.NewAttestationHandler(a.ctx, a.store, a.logger, a.firebaseAPI, a.auth))         // Номер группы")
		})
		r.Get("/week", api.NewWeekParityHandler(a.weekParity)) // Текущая четность недели
		r.Route("/schedule", func(r chi.Router) {
			r.Use(a.authorizationByToken)
			r.Get("/{groupid}", schedule.NewScheduleHandler(a.store, a.logger))                        // Расписание полностью
			r.Get("/{groupid}/by_margin", schedule.NewLessonsHandler(a.store, a.logger, a.weekParity)) // На день с отступом margin от текущего дня
			r.Get("/{groupid}/teachers", schedule.NewTeachersHandler(a.store, a.logger))
			r.Get("/{groupid}/ical", schedule.NewIcalHandler(a.store, a.logger, a.weekParity)) // Импорт файла календаря в виджеты
		})
		r.Route("/icalendar", func(r chi.Router) {
			r.Use(a.parseURLParamsFromTelegramStart)
			r.Use(a.checkSignTelegram)
			r.Get("/", schedule.NewIcalHandler(a.store, a.logger, a.weekParity))
		})
		r.Route("/groups", func(r chi.Router) {
			r.Use(a.authorizationByToken)
			r.Get("/{group}", schedule.NewIdByGroupHandler(a.store, a.logger)) // ID группы по ее номеру
		})
		//r.Get("/", api.NewLessonsHandler(a.store, a.logger))
		r.Route("/feedback", func(r chi.Router) {
			r.Use(a.authorizationByToken)
			r.Post("/", api.NewFeedbackHandler(a.store, a.logger, a.tgApi)) // Отправка отзыва в чат
		})
		r.Route("/attestation", func(r chi.Router) { // DEPRECATED
			r.Use(a.authorizationByToken)
			r.Get("/", api.NewScoreHandler(a.logger))        // DEPRECATED Список оценок БРС
			r.Get("/faculties", api.NewFacHandler(a.logger)) // DEPRECATED Список факультетов
			r.Get("/groups", api.NewGroupsHandler(a.logger)) // DEPRECATED Список групп
			r.Get("/person", api.NewPersonHandler(a.logger)) // DEPRECATED Список фио
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
		r.Route("/tasks", func(r chi.Router) {
			//r.Use(a.authorizationByToken)
			r.Get("/{groupname}", task.NewGetTaskHandler(a.store, a.logger))
			r.Post("/", task.NewTaskHandler(a.store, a.logger))
			r.Delete("/{ID}", task.NewDeleteTaskHandler(a.store, a.logger))
		})
		r.Get("/doc", api.NewDocumentationPageHandler())             // Главная страница документации
		r.Get("/doc/{page}", api.NewDocumentationOtherPageHandler()) // Страница документации
		//r.Post("/token", api.NewRegistrationHandler(a.ctx, a.store, a.logger, a.firebaseAPI)) // DEPRECATED Получение токена
		//r.Get("/token", api.NewRegistrationHandler(a.ctx, a.store, a.logger, a.firebaseAPI))  // DEPRECATED Получение токена
		r.Get("/token/whoiam", api.NewWhoIAmHandler(a.store, a.logger))     // Информация из токена
		r.Get("/confirmation_code/{code}", api.NewSetCodeHandler(&a.store)) // Установка токена верификации до перезагрузки
		r.Get("/echo", api.NewEchoHandler())                                // Эхо  метод для теста
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
	a.router.Route("/portal", func(r chi.Router) { // Портал для сервисов ботов
		r.Route("/sign", func(r chi.Router) {
			r.Use(a.authorizationBySecretPhrase)
			r.Get("/", portal.NewAuthLinkHandler(secretKey)) // Получение подписи (доступно для авторизованных сервисов)
		})
		r.Route("/authorization", func(r chi.Router) {
			r.Use(ParallelHandlerMiddleware)
			r.Use(a.parseURLParamsFromTelegramStart)
			r.Use(a.checkSignTelegram)
			r.Get("/", portal.NewPortalPageHandler(secretKey))
			r.Route("/telegram", func(r chi.Router) {
				r.Post("/", portal.NewAuthTelegramHandler(a.store, a.logger, secretKey, a.auth))
			})
		})
		r.Route("/attestation", func(r chi.Router) {
			r.Use(ParallelHandlerMiddleware)
			r.Use(a.parseURLParamsFromTelegramStart)
			r.Use(a.loadingMiddleware)
			r.Use(a.checkSignTelegram)
			r.Get("/", portal.NewAttestationPageHandler(a.store, a.logger, a.auth, secretKey))
			r.Get("/{id}", portal.NewAttestationElementPageHandler(a.store, a.logger, a.auth))
		})
		r.Route("/exam", func(r chi.Router) {
			r.Use(a.parseURLParamsFromTelegramStart)
			r.Get("/", portal.NewExamPageHandler(a.logger))
		})

	})
	a.router.Route("/payments", func(r chi.Router) {
		r.Use(a.parseURLParamsFromTelegramStart)
		r.Get("/request", pay_handler.NewPaymentRequestHandler(a.logger, a.store, a.pay))
		r.Get("/check/{payment_id}", pay_handler.NewCheckPaymentRequestHandler(a.logger, a.store, a.pay))
		r.Get("/done/{payment_id}", pay_handler.NewDonePaymentPageHandler(a.logger, a.store, a.pay))
		r.Post("/notifications", pay_handler.NewNotificationsPaymentRequestHandler(a.logger, a.store, a.pay))
		r.Get("/subscribe", pay_handler.NewMakePaymentPageHandler())
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

func (a *App) loadingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loading := r.URL.Query().Get("loading")
		if loading == "true" {
			// Если параметр loading=true, показываем страницу загрузки
			w.Write(tools.GetLoadingPage())
			w.WriteHeader(http.StatusOK)
			return
		} else {
			// Иначе, обрабатываем запрос как обычно
			next.ServeHTTP(w, r)
		}
	})
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

func (a *App) parseURLParamsFromTelegramStart(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tgWebAppStartParam := r.URL.Query().Get("tgWebAppStartParam")
		if tgWebAppStartParam != "" {
			resultString, err := processQueryString(tgWebAppStartParam)
			if err != nil {
				log.Println(err)
			} else {
				updateRequestParams(r, resultString)
			}
		}
		rw := &responseWriter{w, http.StatusOK}
		h.ServeHTTP(rw, r)

		return
	})
}
func (a *App) checkSignTelegram(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := a.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
		})
		urlSign := r.FormValue("sign")
		if urlSign == "" {
			urlSign = r.URL.Query().Get("sign")
		}
		sign := portal.GetSignForURLParams(r.URL.Query(), secretKey)
		if sign != urlSign {
			log.Log(
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

func processQueryString(input string) (string, error) {
	// Декодируем URL-параметры
	decoded, err := url.QueryUnescape(input)
	if err != nil {
		return "", err
	}

	// Заменяем тройной знак ___ на знак &
	processed := strings.ReplaceAll(decoded, "___", "&")
	processed = strings.ReplaceAll(processed, "---", "/")

	return processed, nil
}

func updateRequestParams(r *http.Request, queryString string) (*http.Request, error) {
	// Разбираем строку с URL-параметрами в url.Values
	queryValues, err := url.ParseQuery(queryString)
	if err != nil {
		return nil, err
	}

	// Обновляем параметры запроса в r.URL
	r.URL.RawQuery = queryValues.Encode()

	// Обновляем параметры запроса в r.Form
	r.Form = queryValues

	return r, nil
}

// Middleware для параллельной обработки хэндлеров
func ParallelHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			next.ServeHTTP(w, r)
		}()
		wg.Wait()
	})
}
