package image_host_app

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	image_handler "main/internal/app/handlers/image_host"
	"main/internal/app/store/sqlstore"
	"net/http"
	"path/filepath"
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
			r.Post("/groups/tasks", image_handler.NewPostTaskPhotoHandler(a.logger, a.filePath, a.store))
			r.Post("/users", image_handler.NewPostUserProfilePhotoHandler(a.logger, a.filePath, a.store))
		})
	})
	a.router.Handle("/image/users/*", http.StripPrefix("/image/", http.FileServer(http.Dir(filepath.Join(a.filePath)))))
	a.router.Handle("/image/groups/tasks/*", http.StripPrefix("/image/", http.FileServer(http.Dir(filepath.Join(a.filePath)))))
}
