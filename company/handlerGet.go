package company

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type HandlerGet struct {
	repository interface {
		FetchByID(ctx context.Context, id uuid.UUID) (*MongoModel, error)
	}

	validate interface {
		StructCtx(ctx context.Context, s interface{}) (err error)
	}
}

func NewHandlerGet(repository *Repository, validate *validator.Validate) *HandlerGet {
	return &HandlerGet{
		repository: repository,
		validate:   validate,
	}
}

type getParameters struct {
	CompanyID uuid.UUID `param:"companyId" validate:"required"`
}

type getOutput struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Employees   int         `json:"employees"`
	Registered  bool        `json:"registered"`
	Type        CompanyType `json:"type"`
}

func (g HandlerGet) Handle(c echo.Context) error {
	ctx := c.Request().Context()
	p, err := g.fetchParams(c, ctx)
	if err != nil {
		return err
	}
	mongoModel, err := g.repository.FetchByID(ctx, p.CompanyID)
	if err != nil {
		if errors.Is(err, ErrCompanyNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "no company found")
		}
		return err
	}
	_ = c.JSON(http.StatusOK, getOutput{
		ID:          mongoModel.ID,
		Name:        mongoModel.Name,
		Description: mongoModel.Description,
		Employees:   mongoModel.Employees,
		Registered:  mongoModel.Registered,
		Type:        mongoModel.Type,
	})
	return nil
}

func (g HandlerGet) fetchParams(c echo.Context, ctx context.Context) (getParameters, error) {
	p := getParameters{}
	err := c.Bind(&p)
	if err != nil {
		return getParameters{}, echo.NewHTTPError(http.StatusBadRequest, "invalid parameters")
	}

	err = g.validate.StructCtx(ctx, p)
	if err != nil {
		return getParameters{}, echo.NewHTTPError(http.StatusBadRequest, "invalid parameters")
	}
	return p, nil
}
