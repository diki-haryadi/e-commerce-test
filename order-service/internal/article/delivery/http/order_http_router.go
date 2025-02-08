package orderHttpController

import (
	"github.com/labstack/echo/v4"

	orderDomain "github.com/diki-haryadi/go-micro-template/internal/article/domain"
)

type Router struct {
	controller orderDomain.HttpController
}

func NewRouter(controller orderDomain.HttpController) *Router {
	return &Router{
		controller: controller,
	}
}

func (r *Router) Register(e *echo.Group) {
	orderRoutes := e.Group("/orders")
	orderRoutes.Use(r.JWTMiddleware())
	orderRoutes.POST("/checkout", r.controller.Checkout)
}
