package configuration

type Api struct {
	Port         string `env:"PORT" envDefault:"8080"`
	InternalPort string `env:"INTERNAL_PORT" envDefault:""`
	BaseURL      string `env:"BASE_URL" envDefault:""`
}
