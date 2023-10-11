package sqlstore

import (
	"database/sql"
)

// Store ...
type Store struct {
	db                     *sql.DB
	userRepository         *UserRepository
	scheduleRepository     *ScheduleRepository
	verificationRepository *VerificationRepository
	apiRepository          *ApiRepository
	mailingRepository      *MailingRepository
}

// New ...
func New(db *sql.DB) Store {
	return Store{
		db: db,
	}
}

// User ...
func (s *Store) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

func (s *Store) Schedule() *ScheduleRepository {
	if s.scheduleRepository != nil {
		return s.scheduleRepository
	}

	s.scheduleRepository = &ScheduleRepository{
		store: s,
	}

	return s.scheduleRepository
}

func (s *Store) Verification() *VerificationRepository {
	if s.verificationRepository != nil {
		return s.verificationRepository
	}

	s.verificationRepository = &VerificationRepository{
		store: s,
	}

	return s.verificationRepository
}

func (s *Store) API() *ApiRepository {
	if s.apiRepository != nil {
		return s.apiRepository
	}

	s.apiRepository = &ApiRepository{
		store: s,
	}

	return s.apiRepository
}

func (s *Store) Mail() *MailingRepository {
	if s.mailingRepository != nil {
		return s.mailingRepository
	}

	s.mailingRepository = &MailingRepository{
		store: s,
	}

	return s.mailingRepository
}
