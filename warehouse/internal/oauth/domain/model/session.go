package oauthDomain

import "errors"

var (
	// StorageSessionName ...
	StorageSessionName = "go_oauth2_server_session"
	// UserSessionKey ...
	UserSessionKey = "go_oauth2_server_user"
	// ErrSessonNotStarted ...
	ErrSessonNotStarted = errors.New("Session not started")
)

type UserSession struct {
	ClientID     string
	Username     string
	AccessToken  string
	RefreshToken string
}
