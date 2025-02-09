package whHttpController

import (
	productDomain "github.com/diki-haryadi/go-micro-template/internal/warehouse/domain"
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

	warehouse := e.Group("/warehouse")
	warehouse.Use(r.JWTMiddleware())
	warehouse.POST("", r.controller.CreateWarehouse)
	warehouse.PATCH("/:id/status", r.controller.UpdateWarehouseStatus)
	warehouse.POST("/transfer", r.controller.TransferStock)
	warehouse.GET("/:id/stock", r.controller.GetWarehouseStock)
}
