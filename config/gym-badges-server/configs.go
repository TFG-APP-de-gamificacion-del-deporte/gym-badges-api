package configs

import "gym-badges-api/internal/utils"

var (
	Basic BasicConfiguration
)

type BasicConfiguration struct {
	Port int `default:"8080" envconfig:"APP_PORT"`
}

func LoadConfig() {
	utils.LoadGenericConfig(&Basic)
}
