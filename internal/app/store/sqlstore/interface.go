package sqlstore

type StoreInterface interface {
	User() *UserRepository
	Schedule() *ScheduleRepository
	Verification() *VerificationRepository
	API() *ApiRepository
	Mail() *MailingRepository
}
