package oauthHttpController

import (
	"github.com/labstack/echo/v4"

	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain"
)

type Router struct {
	controller oauthDomain.HttpController
}

func NewRouter(controller oauthDomain.HttpController) *Router {
	return &Router{
		controller: controller,
	}
}

func (r *Router) Register(e *echo.Group) {
	oauth := e.Group("/oauth")
	{
		oauth.POST("/tokens", r.controller.Tokens)
		oauth.POST("/introspect", r.controller.Introspect)
	}

}
