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
	SessionDuration int    `default:"3600" envconfig:"SESSION_DURATION"` // TODO Cambiar duración a un año
	JWTKey          string `default:"GymBadges" envconfig:"JWT_KEY"`
	LogLevel        string `default:"DEBUG" envconfig:"LOG_LEVEL"`
	FriendsPageSize int32  `default:"3" envconfig:"FRIENDS_PAGE_SIZE"`
}

func LoadConfig() {
	toolsConfig.LoadGenericConfig(&Basic)
	postgresql.LoadConfig()
	toolsLogging.ConfigureLogger(Basic.LogLevel)
}
