package logging

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
)

func ConfigureLogger(logLvl string) {

	log.SetFormatter(&log.TextFormatter{})

	lvl, err := log.ParseLevel(logLvl)
	if err != nil {
		log.Errorf("error parsing log level: %s", err)
		lvl = log.InfoLevel
	}

	log.SetOutput(os.Stdout)
	log.SetLevel(lvl)
}

func BuildLogger(ctx ...context.Context) *log.Entry {
	if len(ctx) == 0 {
		return log.WithContext(context.Background())
	}
	return log.WithContext(ctx[0])
}
