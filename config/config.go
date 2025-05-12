package config

type Config struct {
	ServiceName string
	Environment string
	Port        string
	DatabaseDSN string

	PublicKeyPath  string
	PrivateKeyPath string
	Passphrase     string
}

func Load() *Config {
	return &Config{
		ServiceName: "bank",
		Environment: "development",
		Port:        "8080",
		DatabaseDSN: "postgres://postgres:postgres@database:5432/bank?sslmode=disable",
		Passphrase:  "verysecret",
	}
}
