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
