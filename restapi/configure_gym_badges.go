// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	loginHandler "gym-badges-api/internal/handler/login"
	userDAO "gym-badges-api/internal/repository/user/postgresql"
	loginService "gym-badges-api/internal/service/login"
	"gym-badges-api/restapi/operations"
	"gym-badges-api/restapi/operations/login"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

//go:generate swagger generate server --target ../../gym-badges-api --name GymBadges --spec ../swagger.yml --principal interface{} --exclude-main

func configureFlags(_ *operations.GymBadgesAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.GymBadgesAPI) http.Handler {

	/*******************************************************************
	DEPENDENCY INJECTION
	*******************************************************************/

	// DAO'S
	userDAO := userDAO.NewUserDAO()

	// SERVICES
	loginService := loginService.NewLoginService(userDAO)

	// HANDLERS
	loginHandler := loginHandler.NewLoginHandler(loginService)

	api.ServeError = errors.ServeError

	api.UseSwaggerUI()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.LoginLoginHandler = login.LoginHandlerFunc(func(params login.LoginParams) middleware.Responder {
		return loginHandler.Login(params)
	})

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(_ *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(_ *http.Server, _, _ string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}