package configuration

type Api struct {
	Port         string  `env:"PORT" envDefault:"8080"`
	InternalPort string  `env:"INTERNAL_PORT" envDefault:""`
	BaseURL      string  `env:"BASE_URL" envDefault:""`
	MongoDB      MongoDB `envPrefix:"MONGODB_"`
}

type MongoDB struct {
	ConnectionString string `env:"CONNECTION_STRING"`
}
