package config

type Config struct {
	Port     string `envconfig:"PREVIEWER_PORT"`
	Capacity int    `envconfig:"PREVIEWER_CAPACITY"`
}
