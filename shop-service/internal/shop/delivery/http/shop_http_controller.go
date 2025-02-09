package shopHttpController

import (
	"github.com/diki-haryadi/go-micro-template/pkg/response"
	"github.com/labstack/echo/v4"
	"strconv"

	shopDomain "github.com/diki-haryadi/go-micro-template/internal/shop/domain"
	shopDto "github.com/diki-haryadi/go-micro-template/internal/shop/dto"
)

type controller struct {
	useCase shopDomain.UseCase
}

func NewController(uc shopDomain.UseCase) shopDomain.HttpController {
	return &controller{
		useCase: uc,
	}
}

func (c *controller) CreateShop(ctx echo.Context) error {
	res := response.NewJSONResponse()

	req := new(shopDto.CreateShopRequestDto)
	if err := ctx.Bind(req); err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	shop, err := c.useCase.CreateShop(ctx.Request().Context(), req)
	if err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	res.APIStatusSuccess().SetData(shop).Send(ctx.Response().Writer)
	return nil
}

func (c *controller) GetShopWarehouses(ctx echo.Context) error {
	res := response.NewJSONResponse()

	shopID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		res.SetError(response.ErrBadRequest).SetMessage("invalid shop id").Send(ctx.Response().Writer)
		return nil
	}

	warehouses, err := c.useCase.GetShopWarehouses(ctx.Request().Context(), shopID)
	if err != nil {
		if err.Error() == "shop not found" {
			res.SetError(response.ErrNotFound).SetMessage(err.Error()).Send(ctx.Response().Writer)
			return nil
		}
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	res.APIStatusSuccess().SetData(warehouses).Send(ctx.Response().Writer)
	return nil
}
