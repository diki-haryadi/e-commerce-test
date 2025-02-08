package productHttpController

import (
	productDomain "github.com/diki-haryadi/go-micro-template/internal/product/domain"
	"github.com/labstack/echo/v4"
)

type Router struct {
	controller productDomain.HttpController
}

func NewRouter(controller productDomain.HttpController) *Router {
	return &Router{
		controller: controller,
	}
}

func (r *Router) Register(e *echo.Group) {

	product := e.Group("/product")
	product.Use(r.JWTMiddleware())
	product.GET("", r.controller.GetProducts)
	product.GET("/:id", r.controller.GetProductByID)
}
