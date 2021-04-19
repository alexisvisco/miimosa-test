package config

import "time"

type App struct {
	Port          string        `envconfig:"PORT" default:"3123"`
	JWTSecret     string        `envconfig:"JWT_SECRET" default:"super_secret"`
	JWTExpiration time.Duration `envconfig:"JWT_EXPIRATION" default:"24h"`
}
