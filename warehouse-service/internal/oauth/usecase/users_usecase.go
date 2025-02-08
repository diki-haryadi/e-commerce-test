package oauthUseCase

import (
	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	"github.com/diki-haryadi/go-micro-template/pkg"
	"github.com/diki-haryadi/go-micro-template/pkg/response"
)

// UserExists returns true if user exists
func (uc *useCase) UserExists(username string) bool {
	_, err := uc.repository.FindUserByUsername(username)
	return err == nil
}

// CreateUser saves a new user to database
func (uc *useCase) CreateUser(roleID, username, password string) (*oauthDomain.Users, error) {
	return uc.repository.CreateUserCommon(roleID, username, password)
}

// CreateUserTx saves a new user to database using injected db object
func (uc *useCase) CreateUserTx(roleID, username, password string) (*oauthDomain.Users, error) {
	return uc.repository.CreateUserCommon(roleID, username, password)
}

// SetPassword sets a user password
func (uc *useCase) SetPassword(user *oauthDomain.Users, password string) error {
	return uc.repository.SetPasswordCommon(user, password)
}

// SetPasswordTx sets a user password in a transaction
func (uc *useCase) SetPasswordTx(user *oauthDomain.Users, password string) error {
	return uc.repository.SetPasswordCommon(user, password)
}

// AuthUser authenticates user
func (uc *useCase) AuthUser(username, password string) (*oauthDomain.Users, error) {
	// Fetch the user
	user, err := uc.repository.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}

	// Check that the password is set
	if !user.Password.Valid {
		return nil, err
	}

	// Verify the password
	if pkg.VerifyPassword(user.Password.String, password) != nil {
		return nil, err
	}

	role, err := uc.repository.FindRoleByID(user.RoleID.String)
	if err != nil {
		return nil, err
	}
	user.Role = role
	return user, nil
}

// UpdateUsername ...
func (uc *useCase) UpdateUsername(user *oauthDomain.Users, username string) error {
	if username == "" {
		return response.ErrCannotSetEmptyUsername
	}

	return uc.repository.UpdateUsernameCommon(user, username)
}

// UpdateUsernameTx ...
func (uc *useCase) UpdateUsernameTx(user *oauthDomain.Users, username string) error {
	return uc.repository.UpdateUsernameCommon(user, username)
}
