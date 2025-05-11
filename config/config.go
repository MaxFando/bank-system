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
		ServiceName:    "my-service",
		Environment:    "development",
		Port:           ":8080",
		DatabaseDSN:    "user:password@tcp(localhost:3306)/dbname",
		PublicKeyPath:  "./keys/public.pem",
		PrivateKeyPath: "./keys/private.pem",
		Passphrase:     "your-passphrase",
	}
}
