package config

type Config struct {
	ListenAddr     string
	RedisAddr      string
	DbURL          string
	MigrationsPath string
}

func Load() *Config {

}
