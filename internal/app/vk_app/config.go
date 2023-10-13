package vk_app

type Config struct {
	BindAddr                      string `toml:"bind_addr"`
	DatabaseURL                   string `toml:"database_url"`
	Chetnost                      int    `toml:"chetnost"`
	FirebaseProjectID             string `toml:"firebase_project_id"`
	FirebaseServiceAccountKeyPath string `toml:"firebase_service_account_key_path""`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
	}
}
