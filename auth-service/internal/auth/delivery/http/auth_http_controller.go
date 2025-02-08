package articleHttpController

import (
	"github.com/diki-haryadi/go-micro-template/config"
	"github.com/diki-haryadi/go-micro-template/pkg"
	"github.com/diki-haryadi/go-micro-template/pkg/response"
	"github.com/labstack/echo/v4"
	
	authDomain "github.com/diki-haryadi/go-micro-template/internal/auth/domain"
	authDto "github.com/diki-haryadi/go-micro-template/internal/auth/dto"
	authException "github.com/diki-haryadi/go-micro-template/internal/auth/exception"
)

type controller struct {
	useCase authDomain.UseCase
}

func NewController(uc authDomain.UseCase) authDomain.HttpController {
	return &controller{
		useCase: uc,
	}
}

func (c controller) SignUp(ctx echo.Context) error {
	res := response.NewJSONResponse()
	aDto := new(authDto.SignUpRequestDto)
	if err := ctx.Bind(aDto); err != nil {
		return authException.ArticleBindingExc()
	}

	if err := aDto.ValidateSignUpDto(); err != nil {
		return authException.CreateArticleValidationExc(err)
	}

	user, err := c.useCase.SignUp(ctx.Request().Context(), aDto)
	if err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	res.APIStatusSuccess().SetData(user).Send(ctx.Response().Writer)
	return nil
}

func (c controller) SignIn(ctx echo.Context) error {
	res := response.NewJSONResponse()
	aDto := new(authDto.SignInRequestDto)
	if err := ctx.Bind(aDto); err != nil {
		return authException.ArticleBindingExc()
	}

	if err := aDto.ValidateSignInDto(); err != nil {
		return authException.CreateArticleValidationExc(err)
	}
	user, err := c.useCase.SignIn(ctx.Request().Context(), aDto.Username, aDto.Password)
	if err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	res.APIStatusSuccess().SetData(user).Send(ctx.Response().Writer)
	return nil
}

func (c controller) RefreshToken(ctx echo.Context) error {
	res := response.NewJSONResponse()
	refreshToken := ctx.FormValue("refresh_token")
	if refreshToken == "" {
		res.SetError(response.ErrBadRequest).SetMessage("refresh token is required").Send(ctx.Response().Writer)
		return nil
	}

	newTokenPair, err := pkg.RefreshTokenPair(refreshToken, config.BaseConfig.App.JWTSecret)
	if err != nil {
		res.SetError(response.ErrUnauthorized).SetMessage("refresh token is required").Send(ctx.Response().Writer)
		return nil
	}

	res.APIStatusSuccess().SetData(newTokenPair).Send(ctx.Response().Writer)
	return nil
}

func (c controller) GetProfile(ctx echo.Context) error {
	res := response.NewJSONResponse()
	// Get user ID from JWT claims
	claims := ctx.Get("claims").(*pkg.TokenClaims)
	if claims == nil {
		res.SetError(response.ErrUnauthorized).SetMessage("invalid token claims").Send(ctx.Response().Writer)
		return nil
	}

	userId := claims.Sub
	if userId == "" {
		res.SetError(response.ErrUnauthorized).SetMessage("user id not found in token").Send(ctx.Response().Writer)
		return nil
	}

	// Get profile from usecase
	profile, err := c.useCase.GetProfile(ctx.Request().Context(), userId)
	if err != nil {
		res.SetError(response.ErrInternalServerError).SetMessage("failed to get profile").Send(ctx.Response().Writer)
		return nil
	}

	res.APIStatusSuccess().SetData(profile).Send(ctx.Response().Writer)
	return nil
}
