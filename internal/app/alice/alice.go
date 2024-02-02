package alice

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

const apiGateway = "https://api.yandex...."

type Alice struct {
	log *logrus.Logger
	ctx context.Context // У них там вроде щас вместо 3 сек можно 5 секунд обрабатываться, надо контекст на все вещи с Алисой делать
}

func NewAlice(ctx context.Context, log *logrus.Logger) *Alice {
	return &Alice{
		log: log,
		ctx: ctx,
	}
}

func (a *Alice) RunWithCtx(ctx context.Context) {
	ctx, cancel := context.WithTimeout(a.ctx, time.Minute*3)
	defer cancel()
	a.log.Info("Start Alice")
	a.run()
}

func (a *Alice) run() {

}
