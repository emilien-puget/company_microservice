package company

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type HandlerDelete struct {
	repository interface {
		FetchByID(ctx context.Context, id uuid.UUID) (*MongoModel, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	validate interface {
		StructCtx(ctx context.Context, s interface{}) (err error)
	}
}

func NewHandlerDelete(repository *Repository, validate *validator.Validate) *HandlerDelete {
	return &HandlerDelete{
		repository: repository,
		validate:   validate,
	}
}

type deleteInput struct {
	CompanyID uuid.UUID `param:"companyId" validate:"required"`
}

func (hd HandlerDelete) Handle(c echo.Context) error {
	p := deleteInput{}
	err := c.Bind(&p)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid parameters")
	}

	ctx := c.Request().Context()
	err = hd.validate.StructCtx(ctx, p)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid parameters")
	}
	_, err = hd.repository.FetchByID(ctx, p.CompanyID)
	if err != nil {
		if errors.Is(err, ErrCompanyNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "no company found")
		}
		return err
	}

	err = hd.repository.Delete(ctx, p.CompanyID)
	if err != nil {
		return err
	}
	return nil
}
