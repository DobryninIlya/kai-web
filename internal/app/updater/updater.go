package updater_app

import (
	"context"
	"github.com/sirupsen/logrus"
	"main/internal/app/store/graph"
	"main/internal/app/updater/get_request"
	"time"
)

const deltaDays = 1

var oldDateUpdate = time.Now().AddDate(-1, 0, 0)         //last year
var newDateUpdate = time.Now().AddDate(0, 0, -deltaDays) //week date ago

type Updater struct {
	groupsParsed []get_request.GroupInfo
	timeout      time.Duration
	store        *graph.Store
	log          *logrus.Logger
	ctx          context.Context
}

func (s *Updater) Run() {
	parsedRes := get_request.GetGroupsList()
	s.log.Info("Groups parsed, length: ", len(parsedRes))
	UpdateSchedule(s.store, parsedRes)

}

func NewUpdater(ctx context.Context, timeout time.Duration, log *logrus.Logger, store *graph.Store) (*Updater, error) {
	return &Updater{
		ctx:     ctx,
		timeout: timeout,
		store:   store,
		log:     log,
	}, nil
}
