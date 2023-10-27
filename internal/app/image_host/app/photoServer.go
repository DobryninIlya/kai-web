package image_host_app

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	image_handler "main/internal/app/handlers/image_host"
	"main/internal/app/store/sqlstore"
	"net/http"
)

type App struct {
	router   *chi.Mux
	server   *http.Server
	store    sqlstore.StoreInterface
	logger   *logrus.Logger
	ctx      context.Context
	filePath string
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func newApp(ctx context.Context, store sqlstore.StoreInterface, bindAddr string, config Config) *App {
	router := chi.NewRouter()
	server := &http.Server{
		Addr:    bindAddr,
		Handler: router,
	}
	logger := logrus.New()
	a := &App{
		router:   router,
		server:   server,
		store:    store,
		logger:   logger,
		ctx:      ctx,
		filePath: config.StorePath,
	}
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
	a.router.Route("/api", func(r chi.Router) {
		r.Route("/image", func(r chi.Router) {
			r.Post("/tasks", image_handler.NewPostPhotoHandler(a.logger, a.filePath, a.store))
		})
	})

}
