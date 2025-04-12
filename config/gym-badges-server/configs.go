package configs

import (
	"gym-badges-api/internal/repository/config/postgresql"
	toolsConfig "gym-badges-api/tools/config"
	toolsLogging "gym-badges-api/tools/logging"
)

var (
	Basic BasicConfiguration
)

type BasicConfiguration struct {
	Port             int    `default:"8080" envconfig:"APP_PORT"`
	SessionDuration  int    `default:"31536000" envconfig:"SESSION_DURATION"` // One year
	JWTKey           string `default:"GymBadges" envconfig:"JWT_KEY"`
	LogLevel         string `default:"DEBUG" envconfig:"LOG_LEVEL"`
	FriendsPageSize  int32  `default:"3" envconfig:"FRIENDS_PAGE_SIZE"`
	RankingsPageSize int32  `default:"10" envconfig:"RANKINGS_PAGE_SIZE"`
}

func LoadConfig() {
	toolsConfig.LoadGenericConfig(&Basic)
	postgresql.LoadConfig()
	toolsLogging.ConfigureLogger(Basic.LogLevel)
}
