package oauthHttpController

import (
	"github.com/diki-haryadi/go-micro-template/config"
	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain"
	oauthModel "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	oauthDto "github.com/diki-haryadi/go-micro-template/internal/oauth/dto"
	"github.com/diki-haryadi/go-micro-template/pkg/response"
	"github.com/labstack/echo/v4"
)

type controller struct {
	useCase oauthDomain.UseCase
}

func NewController(uc oauthDomain.UseCase) oauthDomain.HttpController {
	return &controller{
		useCase: uc,
	}
}

func (c controller) Tokens(ctx echo.Context) error {
	res := response.NewJSONResponse()
	grantTypes := map[string]func(ctx echo.Context) error{
		"authorization_code": c.AuthorizationCodeGrant,
		"password":           c.PasswordGrant,
		"client_credentials": c.ClientCredentialsGrant,
		"refresh_token":      c.RefreshTokenGrant,
	}

	// Check the grant type
	grantHandler, ok := grantTypes[ctx.Request().FormValue("grant_type")]
	if !ok {
		res.SetError(response.ErrInvalidGrantType).SetMessage(response.ErrInvalidGrantType.Error()).Send(ctx.Response().Writer)
		return nil
	}

	// Grant processing
	err := grantHandler(ctx)
	if err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	return nil
}

func (c controller) AuthorizationCodeGrant(ctx echo.Context) error {
	res := response.NewJSONResponse()

	aDto := new(oauthDto.AuthorizationCodeGrantRequestDto).GetFields(ctx)
	if err := aDto.ValidateAuthorizationCodeDto(); err != nil {
		res.SetError(response.ErrInvalidAuthorizationCodeGrantRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	client, err := c.BasicAuthClient(ctx)
	if err != nil {
		res.SetError(err).Send(ctx.Response().Writer)
		return nil
	}

	acGrant, err := c.useCase.AuthorizationCodeGrant(
		ctx.Request().Context(),
		aDto.Code,
		aDto.RedirectUri,
		aDto.ToModel(client.ID))

	if err != nil {
		res.SetError(err).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	res.APIStatusSuccess().SetData(*acGrant).Send(ctx.Response().Writer)
	return nil
}

func (c controller) PasswordGrant(ctx echo.Context) error {
	res := response.NewJSONResponse()
	aDto := new(oauthDto.PasswordGrantRequestDto).GetFields(ctx)
	if err := aDto.ValidatePasswordDto(); err != nil {
		res.SetError(response.ErrInvalidPasswordGrantRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	client, err := c.BasicAuthClient(ctx)
	if err != nil {
		res.SetError(err).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	acGrant, err := c.useCase.PasswordGrant(
		ctx.Request().Context(),
		aDto.Username,
		aDto.Password,
		aDto.Scope,
		aDto.ToModel(client.ID))

	if err != nil {
		res.SetError(err).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	res.APIStatusSuccess().SetData(*acGrant).Send(ctx.Response().Writer)
	return nil
}

func (c controller) ClientCredentialsGrant(ctx echo.Context) error {
	res := response.NewJSONResponse()
	aDto := new(oauthDto.ClientCredentialsGrantRequestDto).GetFields(ctx)
	if err := aDto.ValidateClientCredentialsDto(); err != nil {
		res.SetError(response.ErrInvalidClientCredentialsGrantRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	client, err := c.BasicAuthClient(ctx)
	if err != nil {
		res.SetError(err).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	acGrant, err := c.useCase.GrantAccessToken(
		ctx.Request().Context(),
		client,
		nil,
		config.BaseConfig.App.ConfigOauth.Oauth.AccessTokenLifetime,
		aDto.Scope)

	if err != nil {
		res.SetError(err).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	res.APIStatusSuccess().SetData(&acGrant).Send(ctx.Response().Writer)
	return nil
}

func (c controller) RefreshTokenGrant(ctx echo.Context) error {
	res := response.NewJSONResponse()
	aDto := new(oauthDto.RefreshTokenRequestDto).GetFields(ctx)
	if err := aDto.ValidateRefreshTokenDto(); err != nil {
		res.SetError(response.ErrInvalidClientCredentialsGrantRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	client, err := c.BasicAuthClient(ctx)
	if err != nil {
		res.SetError(err).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	acGrant, err := c.useCase.ClientCredentialsGrant(
		ctx.Request().Context(),
		aDto.RefreshToken,
		aDto.Scope,
		aDto.ToModel(client.ID))

	if err != nil {
		res.SetError(err).SetMessage(err.Error()).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	res.APIStatusSuccess().SetData(*acGrant).Send(ctx.Response().Writer)
	return nil
}

func (c controller) Introspect(ctx echo.Context) error {
	res := response.NewJSONResponse()
	aDto := new(oauthDto.IntrospectRequestDto).GetFields(ctx)
	if err := aDto.ValidateIntrospectDto(); err != nil {
		res.SetError(response.ErrInvalidIntrospectRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	client, err := c.BasicAuthClient(ctx)
	if err != nil {
		res.SetError(err).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	introspect, err := c.useCase.IntrospectToken(
		ctx.Request().Context(),
		aDto.Token,
		aDto.TokenTypeHint,
		aDto.ToModel(client.ID))

	if err != nil {
		res.SetError(err).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	res.APIStatusSuccess().SetData(*introspect).Send(ctx.Response().Writer)
	return nil
}

func (c controller) BasicAuthClient(ctx echo.Context) (*oauthModel.Client, error) {
	// Get client credentials from basic auth
	clientID, secret, ok := ctx.Request().BasicAuth()
	if !ok {
		return nil, response.ErrInvalidClientIDOrSecret
	}

	// Authenticate the client
	client, err := c.useCase.AuthClient(ctx.Request().Context(), clientID, secret)
	if err != nil {
		return nil, response.ErrInvalidClientIDOrSecret
	}

	return client, nil
}
