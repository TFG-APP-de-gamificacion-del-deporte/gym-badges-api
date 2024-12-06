// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	badgeHandler "gym-badges-api/internal/handler/badge"
	friendsHandler "gym-badges-api/internal/handler/friends"
	loginHandler "gym-badges-api/internal/handler/login"
	statsHandler "gym-badges-api/internal/handler/stats"
	userHandler "gym-badges-api/internal/handler/user"
	badgeDAO "gym-badges-api/internal/repository/badge/postgresql"
	userDAO "gym-badges-api/internal/repository/user/postgresql"
	badgeService "gym-badges-api/internal/service/badge"
	friendsService "gym-badges-api/internal/service/friends"
	loginService "gym-badges-api/internal/service/login"
	sessionService "gym-badges-api/internal/service/session"
	statsService "gym-badges-api/internal/service/stats"
	userService "gym-badges-api/internal/service/user"
	"gym-badges-api/restapi/operations"
	"gym-badges-api/restapi/operations/badges"
	"gym-badges-api/restapi/operations/friends"
	"gym-badges-api/restapi/operations/login"
	"gym-badges-api/restapi/operations/login_with_token"
	"gym-badges-api/restapi/operations/stats"
	"gym-badges-api/restapi/operations/user"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/security"
)

//go:generate swagger generate server --target ../../gym-badges-api --name GymBadges --spec ../swagger.yml --principal interface{} --exclude-main

const (
	securityHeader   = "token"
	authUserIDHeader = "auth_user_id"
	successMsg       = "SUCCESS"
)

func configureFlags(_ *operations.GymBadgesAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.GymBadgesAPI) http.Handler {

	/*******************************************************************
	DEPENDENCY INJECTION
	*******************************************************************/

	// DAO'S
	userDAO := userDAO.NewUserDAO()
	badgeDAO := badgeDAO.NewBadgeDAO()

	// SERVICES
	sessionService := sessionService.NewSessionService()
	loginService := loginService.NewLoginService(userDAO, sessionService)
	userService := userService.NewUserService(userDAO, sessionService)
	statsService := statsService.NewStatsService(userDAO, sessionService)
	friendsService := friendsService.NewFriendsService(userDAO)
	badgeService := badgeService.NewBadgeService(userDAO, badgeDAO)

	// HANDLERS
	loginHandler := loginHandler.NewLoginHandler(loginService)
	userHandler := userHandler.NewUserHandler(userService)
	statsHandler := statsHandler.NewStatsHandler(statsService)
	friendsHandler := friendsHandler.NewFriendsHandler(friendsService)
	badgeHandler := badgeHandler.NewBadgeHandler(badgeService)

	api.ServeError = errors.ServeError

	api.UseSwaggerUI()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.LoginLoginHandler = login.LoginHandlerFunc(func(params login.LoginParams) middleware.Responder {
		return loginHandler.Login(params)
	})

	api.UserCreateUserHandler = user.CreateUserHandlerFunc(func(params user.CreateUserParams) middleware.Responder {
		return userHandler.CreateUser(params)
	})

	api.UserGetUserInfoHandler = user.GetUserInfoHandlerFunc(func(params user.GetUserInfoParams, new interface{}) middleware.Responder {
		return userHandler.GetUser(params)
	})

	api.LoginWithTokenLoginWithTokenHandler = login_with_token.LoginWithTokenHandlerFunc(func(params login_with_token.LoginWithTokenParams, new interface{}) middleware.Responder {
		return login_with_token.NewLoginWithTokenOK()
	})

	// *******************************************************************
	// STATS
	// *******************************************************************

	api.StatsGetWeightHistoryByUserIDHandler = stats.GetWeightHistoryByUserIDHandlerFunc(func(params stats.GetWeightHistoryByUserIDParams, new interface{}) middleware.Responder {
		return statsHandler.GetWeightHistory(params)
	})

	api.StatsAddWeightHandler = stats.AddWeightHandlerFunc(func(params stats.AddWeightParams, new interface{}) middleware.Responder {
		return statsHandler.AddWeight(params)
	})

	api.StatsGetFatHistoryByUserIDHandler = stats.GetFatHistoryByUserIDHandlerFunc(func(params stats.GetFatHistoryByUserIDParams, new interface{}) middleware.Responder {
		return statsHandler.GetFatHistory(params)
	})

	api.StatsGetStreakCalendarByUserIDHandler = stats.GetStreakCalendarByUserIDHandlerFunc(func(params stats.GetStreakCalendarByUserIDParams, new interface{}) middleware.Responder {
		return statsHandler.GetStreakCalendar(params)
	})

	// *******************************************************************
	// FRIENDS
	// *******************************************************************

	api.FriendsGetFriendsByUserIDHandler = friends.GetFriendsByUserIDHandlerFunc(func(params friends.GetFriendsByUserIDParams, new interface{}) middleware.Responder {
		return friendsHandler.GetFriendsByUserID(params)
	})

	api.FriendsAddFriendHandler = friends.AddFriendHandlerFunc(func(params friends.AddFriendParams, new interface{}) middleware.Responder {
		return friendsHandler.AddFriend(params)
	})

	api.FriendsDeleteFriendHandler = friends.DeleteFriendHandlerFunc(func(params friends.DeleteFriendParams, new interface{}) middleware.Responder {
		return friendsHandler.DeleteFriend(params)
	})

	// *******************************************************************
	// BADGES
	// *******************************************************************

	api.BadgesGetBadgesByUserIDHandler = badges.GetBadgesByUserIDHandlerFunc(func(params badges.GetBadgesByUserIDParams, new interface{}) middleware.Responder {
		return badgeHandler.GetBadgesByUserID(params)
	})

	// Authentication Middleware
	api.APIKeyAuthenticator = func(_ string, _ string, authentication security.TokenAuthentication) runtime.Authenticator {
		return Authenticator{sessionService: sessionService}
	}

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

type Authenticator struct {
	sessionService sessionService.ISessionService
}

func (a Authenticator) Authenticate(data interface{}) (bool, interface{}, error) {

	authRequest := data.(*security.ScopedAuthRequest)

	sessionID := authRequest.Request.Header.Get(securityHeader)
	userID := authRequest.Request.Header.Get(authUserIDHeader)

	if err := a.sessionService.ValidateSession(userID, sessionID); err != nil {
		return false, nil, nil
	}

	return true, successMsg, nil
}
