package logging

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
)

func ConfigureLogger() {
	log.SetFormatter(&log.JSONFormatter{
		PrettyPrint: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func BuildLogger(ctx ...context.Context) *log.Entry {
	if len(ctx) == 0 {
		return log.New().WithContext(context.Background())
	}
	return log.New().WithContext(ctx[0])
}
