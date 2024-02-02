package graph

import "github.com/jmcvetta/neoism"

type Store struct {
	db       *neoism.Database
	schedule *Schedule
}

func NewGraphStore(db *neoism.Database) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Schedule() ScheduleInterface {
	if s.schedule != nil {
		return s.schedule
	}

	s.schedule = &Schedule{
		store: s,
	}

	return s.schedule
}
