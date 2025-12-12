package config

type Config struct {
	Port string
	Host string
}

func Load() *Config {
	return &Config{
		Port: "8080",
		Host: "0.0.0.0",
	}
}

