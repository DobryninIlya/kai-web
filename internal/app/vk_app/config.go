package vk_app

type Config struct {
	BindAddr    string `toml:"bind_addr"`
	DatabaseURL string `toml:"database_url"`
	Chetnost    int    `toml:"chetnost"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		//LogLevel: "debug",
	}
}
