package session_service

type ISessionService interface {
	GenerateSession(userID string) (string, error)
	ValidateSession(userID string, sessionID string) error
}
