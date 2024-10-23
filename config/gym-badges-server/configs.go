package configs

import (
	"gym-badges-api/internal/repository/user/postgresql"
	toolsConfig "gym-badges-api/tools/config"
	toolsLogging "gym-badges-api/tools/logging"
)

var (
	Basic BasicConfiguration
)

type BasicConfiguration struct {
	Port            int    `default:"8080" envconfig:"APP_PORT"`
	SessionDuration int    `default:"3600" envconfig:"SESSION_DURATION"`
	JWTKey          string `default:"GymBadges" envconfig:"JWT_KEY"`
}

func LoadConfig() {
	toolsConfig.LoadGenericConfig(&Basic)
	postgresql.LoadConfig()
	toolsLogging.ConfigureLogger()
}
