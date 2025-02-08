package orderHttpController

import (
	"github.com/diki-haryadi/go-micro-template/pkg/response"
	"github.com/labstack/echo/v4"

	orderDomain "github.com/diki-haryadi/go-micro-template/internal/article/domain"
	orderDto "github.com/diki-haryadi/go-micro-template/internal/article/dto"
)

type controller struct {
	useCase orderDomain.UseCase
}

func NewController(uc orderDomain.UseCase) orderDomain.HttpController {
	return &controller{
		useCase: uc,
	}
}

func (c *controller) Checkout(ctx echo.Context) error {
	res := response.NewJSONResponse()
	req := new(orderDto.CheckoutRequestDto)

	if err := ctx.Bind(req); err != nil {
		res.SetError(response.ErrBadRequest).
			SetMessage(err.Error()).
			Send(ctx.Response().Writer)
		return nil
	}

	if err := req.Validate(); err != nil {
		res.SetError(response.ErrBadRequest).
			SetMessage(err.Error()).
			Send(ctx.Response().Writer)
		return nil
	}

	result, err := c.useCase.Checkout(ctx.Request().Context(), req)
	if err != nil {
		if err.Error() == "insufficient stock" {
			res.SetError(response.ErrBadRequest).
				SetMessage("insufficient stock").
				Send(ctx.Response().Writer)
			return nil
		}
		res.SetError(response.ErrInternalServerError).
			SetMessage(err.Error()).
			Send(ctx.Response().Writer)
		return nil
	}

	res.APIStatusSuccess().SetData(result).Send(ctx.Response().Writer)
	return nil
}
