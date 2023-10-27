package image_host_app

type Config struct {
	StorePath   string `toml:"store_path"`
	BindAddr    string `toml:"bind_addr"`
	DatabaseURL string `json:"database_url"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8283",
	}
}
