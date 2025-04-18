package main

import (
	configs "gym-badges-api/config/gym-badges-server"
	"gym-badges-api/restapi"
	"gym-badges-api/restapi/operations"
	toolsLogging "gym-badges-api/tools/logging"
	"os"

	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"
)

func main() {

	configs.LoadConfig()

	ctxLog := toolsLogging.BuildLogger()

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		ctxLog.Fatalln(err)
	}

	api := operations.NewGymBadgesAPI(swaggerSpec)
	server := restapi.NewServer(api)

	defer func(server *restapi.Server) {
		if err := server.Shutdown(); err != nil {
			ctxLog.Error(err)
		}
	}(server)

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "Gym Badges API"
	parser.LongDescription = swaggerSpec.Spec().Info.Description
	server.ConfigureFlags()

	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			ctxLog.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	server.Host = "0.0.0.0"
	server.Port = configs.Basic.Port

	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		ctxLog.Fatalln(err)
	}

}
