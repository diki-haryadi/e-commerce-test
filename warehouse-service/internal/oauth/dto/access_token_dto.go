package oauthDto

import (
	"fmt"
	oauthModel "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	"github.com/diki-haryadi/go-micro-template/pkg"
	"github.com/google/uuid"
	"time"
)

// AccessTokenResponse ...
type AccessTokenResponse struct {
	UserID       string `json:"user_id,omitempty"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// NewAccessTokenResponse ...
func NewAccessTokenResponse(accessToken *oauthModel.AccessToken, refreshToken *oauthModel.RefreshToken, lifetime int, theTokenType string) (*AccessTokenResponse, error) {
	response := &AccessTokenResponse{
		AccessToken: accessToken.Token,
		ExpiresIn:   lifetime,
		TokenType:   theTokenType,
		Scope:       accessToken.Scope,
	}
	if accessToken.UserID.Valid {
		response.UserID = accessToken.UserID.String
	}
	if refreshToken != nil {
		response.RefreshToken = refreshToken.Token
	}
	return response, nil
}

func NewOauthAccessToken(client *oauthModel.Client, user *oauthModel.Users, expiresIn int, scope string) *oauthModel.AccessToken {
	accessToken := &oauthModel.AccessToken{
		Common: oauthModel.Common{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
		},
		ClientID:  pkg.StringOrNull(fmt.Sprint(client.ID)),
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().UTC().Add(time.Duration(expiresIn) * time.Second),
		Scope:     scope,
	}
	if user != nil {
		accessToken.UserID = pkg.StringOrNull(fmt.Sprint(user.ID))
	}
	return accessToken
}
