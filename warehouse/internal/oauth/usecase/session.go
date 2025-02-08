package oauthUseCase

import (
	"encoding/gob"
	"errors"
	"github.com/diki-haryadi/go-micro-template/config"
	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	"github.com/diki-haryadi/go-micro-template/pkg/constant"
	"github.com/diki-haryadi/go-micro-template/pkg/response"
	"github.com/gorilla/sessions"
	"net/http"
)

type Service struct {
	sessionStore   sessions.Store
	sessionOptions *sessions.Options
	session        *sessions.Session
	r              *http.Request
	w              http.ResponseWriter
}

func init() {
	// Register a new datatype for storage in sessions
	gob.Register(new(oauthDomain.UserSession))
}

// NewService returns a new Service instance
func NewService(cnf *config.Config, sessionStore sessions.Store) *Service {
	return &Service{
		// Session cookie storage
		sessionStore: sessionStore,
		// Session options
		sessionOptions: &sessions.Options{
			Path:     cnf.App.ConfigOauth.Session.Path,
			MaxAge:   cnf.App.ConfigOauth.Session.MaxAge,
			HttpOnly: cnf.App.ConfigOauth.Session.HTTPOnly,
		},
	}
}

// SetSessionService sets the request and responseWriter on the session service
func (s *Service) SetSessionService(r *http.Request, w http.ResponseWriter) {
	s.r = r
	s.w = w
}

// StartSession starts a new session. This method must be called before other
// public methods of this struct as it sets the internal session object
func (s *Service) StartSession() error {
	session, err := s.sessionStore.Get(s.r, constant.StorageSessionName)
	if err != nil {
		return err
	}
	s.session = session
	return nil
}

// GetUserSession returns the user session
func (s *Service) GetUserSession() (*oauthDomain.UserSession, error) {
	// Make sure StartSession has been called
	if s.session == nil {
		return nil, response.ErrSessonNotStarted
	}

	// Retrieve our user session struct and type-assert it
	userSession, ok := s.session.Values[constant.UserSessionKey].(*oauthDomain.UserSession)
	if !ok {
		return nil, errors.New("User session type assertion error")
	}

	return userSession, nil
}

// SetUserSession saves the user session
func (s *Service) SetUserSession(userSession *oauthDomain.UserSession) error {
	// Make sure StartSession has been called
	if s.session == nil {
		return response.ErrSessonNotStarted
	}

	// Set a new user session
	s.session.Values[constant.UserSessionKey] = userSession
	return s.session.Save(s.r, s.w)
}

// ClearUserSession deletes the user session
func (s *Service) ClearUserSession() error {
	// Make sure StartSession has been called
	if s.session == nil {
		return response.ErrSessonNotStarted
	}

	// Delete the user session
	delete(s.session.Values, constant.UserSessionKey)
	return s.session.Save(s.r, s.w)
}

// SetFlashMessage sets a flash message,
// useful for displaying an error after 302 redirection
func (s *Service) SetFlashMessage(msg string) error {
	// Make sure StartSession has been called
	if s.session == nil {
		return response.ErrSessonNotStarted
	}

	// Add the flash message
	s.session.AddFlash(msg)
	return s.session.Save(s.r, s.w)
}

// GetFlashMessage returns the first flash message
func (s *Service) GetFlashMessage() (interface{}, error) {
	// Make sure StartSession has been called
	if s.session == nil {
		return nil, response.ErrSessonNotStarted
	}

	// Get the last flash message from the stack
	if flashes := s.session.Flashes(); len(flashes) > 0 {
		// We need to save the session, otherwise the flash message won't be removed
		s.session.Save(s.r, s.w)
		return flashes[0], nil
	}

	// No flash messages in the stack
	return nil, nil
}

// Close stops any running services
func (s *useCase) Close() {}
