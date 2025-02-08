package articleHttpController

import (
	"github.com/labstack/echo/v4"

	articleDomain "github.com/diki-haryadi/go-micro-template/internal/article/domain"
)

type Router struct {
	controller articleDomain.HttpController
}

func NewRouter(controller articleDomain.HttpController) *Router {
	return &Router{
		controller: controller,
	}
}

func (r *Router) Register(e *echo.Group) {
	orderRoutes := e.Group("/orders")
	orderRoutes.Use(r.JWTMiddleware())
	orderRoutes.POST("/checkout", r.controller.Checkout)
}
