package configs

import "os"

type Config struct {
	JWTSecret              string
	JWTExpirationinSeconds int64
}

var Configs = initConfig()

func initConfig() Config {
	return Config{
		JWTSecret:              os.Getenv("JWT_Secret"),
		JWTExpirationinSeconds: 3600 * 24 * 7,
	}
}
