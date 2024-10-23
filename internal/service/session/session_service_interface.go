package session_service

type ISessionService interface {
	GenerateSession(username string) (string, error)
	ValidateSession(sessionID string) error
}
