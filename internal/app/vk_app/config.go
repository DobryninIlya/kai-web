package vk_app

type Config struct {
	BindAddr                      string `toml:"bind_addr"`
	DatabaseURL                   string `toml:"database_url"`
	Chetnost                      int    `toml:"chetnost"`
	FirebaseProjectID             string `toml:"firebase_project_id"`
	FirebaseServiceAccountKeyPath string `toml:"firebase_service_account_key_path"`
	InfluxDBName                  string `toml:"influxdb_name"`
	InfluxDBToken                 string `toml:"influxdb_token"`
	InfluxDBURL                   string `toml:"influxdb_url"`
	ShopID                        int    `toml:"yookassa_shop_id"`
	APIKey                        string `toml:"yookassa_api_key"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
	}
}
