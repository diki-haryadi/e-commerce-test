package shopHttpController

import (
	shopDomain "github.com/diki-haryadi/go-micro-template/internal/shop/domain"
	"github.com/labstack/echo/v4"
)

type Router struct {
	controller shopDomain.HttpController
}

func NewRouter(controller shopDomain.HttpController) *Router {
	return &Router{
		controller: controller,
	}
}

func (r *Router) Register(e *echo.Group) {
	shop := e.Group("/shop")
	shop.Use(r.JWTMiddleware())

	// Register shop endpoints
	shop.POST("", r.controller.CreateShop)
	shop.GET("/:id/warehouses", r.controller.GetShopWarehouses)
}
