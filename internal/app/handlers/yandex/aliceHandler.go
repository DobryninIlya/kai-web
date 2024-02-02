package alice

import (
	"context"
	"github.com/sirupsen/logrus"
	"main/internal/app/alice"
	"main/internal/app/handlers/web_app"
	"main/internal/app/store/sqlstore"
	"net/http"
)

func NewAliceHandler(ctx context.Context, store sqlstore.StoreInterface, log *logrus.Logger, alice *alice.Alice) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Здесь обработка запросов от Алисы, надо демаршалить в новую структуру (хранятся в фолдере models, надо прописать там под их ответы)
		// И реализовать обработчик alice.Alice.RunWithCtx(ctx)
		log.Logf(logrus.InfoLevel, "Request to AliceHandler: %v", r.URL.Path)
		web_app.RespondAPI(w, r, http.StatusNotFound, "You may answer here")
	}
}
