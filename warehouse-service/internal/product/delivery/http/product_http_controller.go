package productHttpController

import (
	"github.com/diki-haryadi/go-micro-template/pkg/response"
	"github.com/labstack/echo/v4"

	productDomain "github.com/diki-haryadi/go-micro-template/internal/product/domain"
	productDto "github.com/diki-haryadi/go-micro-template/internal/product/dto"
)

type controller struct {
	useCase productDomain.UseCase
}

func NewController(uc productDomain.UseCase) productDomain.HttpController {
	return &controller{
		useCase: uc,
	}
}

func (c controller) GetProducts(ctx echo.Context) error {
	res := response.NewJSONResponse()

	filter := new(productDto.GetProductsRequestDto)
	if err := ctx.Bind(filter); err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	products, err := c.useCase.GetProducts(ctx.Request().Context(), filter)
	if err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	res.APIStatusSuccess().SetData(products).Send(ctx.Response().Writer)
	return nil
}

func (c controller) GetProductByID(ctx echo.Context) error {
	res := response.NewJSONResponse()

	param := new(productDto.GetProductByIDRequestDto)
	if err := ctx.Bind(param); err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	product, err := c.useCase.GetProductByID(ctx.Request().Context(), param.ID)
	if err != nil {
		if err.Error() == "product not found" {
			res.SetError(response.ErrNotFound).SetMessage(err.Error()).Send(ctx.Response().Writer)
			return nil
		}
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	res.APIStatusSuccess().SetData(product).Send(ctx.Response().Writer)
	return nil
}
