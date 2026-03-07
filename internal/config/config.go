package config

type Config struct {
	ServerAdress string
}

func NewConfig() *Config {
	return &Config{
		ServerAdress: ":8080",
	}
}
