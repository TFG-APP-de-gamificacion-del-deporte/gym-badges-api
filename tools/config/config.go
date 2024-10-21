package config

import (
	"gym-badges-api/internal/constants"
	"log"

	"github.com/kelseyhightower/envconfig"
)

func LoadGenericConfig(config interface{}) {
	err := envconfig.Process(constants.EmptyString, config)
	if err != nil {
		log.Panicf("error loading environment configuration: %s", err.Error())
	}
}
