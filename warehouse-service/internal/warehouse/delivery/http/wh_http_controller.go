package whHttpController

import (
	whDomain "github.com/diki-haryadi/go-micro-template/internal/warehouse/domain"
	whDto "github.com/diki-haryadi/go-micro-template/internal/warehouse/dto"
	"github.com/diki-haryadi/go-micro-template/pkg/response"
	"github.com/labstack/echo/v4"
)

type controller struct {
	useCase whDomain.UseCase
}

func NewController(uc whDomain.UseCase) whDomain.HttpController {
	return &controller{
		useCase: uc,
	}
}

func (c controller) CreateWarehouse(ctx echo.Context) error {
	res := response.NewJSONResponse()

	req := new(whDto.CreateWarehouseRequestDto)
	if err := ctx.Bind(req); err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	result, err := c.useCase.CreateWarehouse(ctx.Request().Context(), req)
	if err != nil {
		if err.Error() == "warehouse already exists" {
			res.SetError(response.ErrConflict).SetMessage(err.Error()).Send(ctx.Response().Writer)
			return nil
		}
		res.SetError(response.ErrInternalServerError).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	res.APIStatusSuccess().SetData(result).Send(ctx.Response().Writer)
	return nil
}

func (c controller) UpdateWarehouseStatus(ctx echo.Context) error {
	res := response.NewJSONResponse()

	id := ctx.Param("id")

	req := new(whDto.UpdateWarehouseStatusRequestDto)
	if err := ctx.Bind(req); err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}
	req.ID = id

	err := c.useCase.UpdateWarehouseStatus(ctx.Request().Context(), req)
	if err != nil {
		if err.Error() == "warehouse not found" {
			res.SetError(response.ErrNotFound).SetMessage(err.Error()).Send(ctx.Response().Writer)
			return nil
		}
		res.SetError(response.ErrInternalServerError).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	res.APIStatusSuccess().SetMessage("Status updated successfully").Send(ctx.Response().Writer)
	return nil
}

func (c controller) TransferStock(ctx echo.Context) error {
	res := response.NewJSONResponse()

	req := new(whDto.TransferStockRequestDto)
	if err := ctx.Bind(req); err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	err := c.useCase.TransferStock(ctx.Request().Context(), req)
	if err != nil {
		switch err.Error() {
		case "source warehouse not found", "destination warehouse not found":
			res.SetError(response.ErrNotFound).SetMessage(err.Error()).Send(ctx.Response().Writer)
			return nil
		case "source warehouse is not active", "destination warehouse is not active":
			res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
			return nil
		case "insufficient stock in source warehouse":
			res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(ctx.Response().Writer)
			return nil
		default:
			res.SetError(response.ErrInternalServerError).SetMessage(err.Error()).Send(ctx.Response().Writer)
			return nil
		}
	}

	res.APIStatusSuccess().SetMessage("Stock transferred successfully").Send(ctx.Response().Writer)
	return nil
}

func (c controller) GetWarehouseStock(ctx echo.Context) error {
	res := response.NewJSONResponse()

	id := ctx.Param("id")

	result, err := c.useCase.GetWarehouseStock(ctx.Request().Context(), id)
	if err != nil {
		if err.Error() == "warehouse not found" {
			res.SetError(response.ErrNotFound).SetMessage(err.Error()).Send(ctx.Response().Writer)
			return nil
		}
		res.SetError(response.ErrInternalServerError).SetMessage(err.Error()).Send(ctx.Response().Writer)
		return nil
	}

	res.APIStatusSuccess().SetData(result).Send(ctx.Response().Writer)
	return nil
}
