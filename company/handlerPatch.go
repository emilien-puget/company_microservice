package company

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type HandlerPatch struct {
	repository interface {
		FetchByID(ctx context.Context, id uuid.UUID) (*MongoModel, error)
		FetchByName(ctx context.Context, name string) (*MongoModel, error)
		Patch(ctx context.Context, id uuid.UUID, updatedCompany *MongoModel) error
	}

	validate interface {
		StructCtx(ctx context.Context, s interface{}) (err error)
	}
}

func NewHandlerPatch(repository *Repository, validate *validator.Validate) *HandlerPatch {
	return &HandlerPatch{
		repository: repository,
		validate:   validate,
	}
}

type patchParameters struct {
	CompanyID   uuid.UUID   `param:"companyId" validate:"required"`
	Name        string      `json:"name" validate:"required,max=15"`
	Description string      `json:"description,omitempty" validate:"max=3000"`
	Employees   int         `json:"employees" validate:"required"`
	Registered  bool        `json:"registered"`
	Type        CompanyType `json:"type" validate:"required"`
}

func (hp HandlerPatch) Handle(c echo.Context) error {
	p := patchParameters{}
	err := c.Bind(&p)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid parameters")
	}

	ctx := c.Request().Context()
	err = hp.validate.StructCtx(ctx, p)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid parameters")
	}

	if !p.Type.IsValid() {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid company type")
	}

	ogMongoModel, err := hp.repository.FetchByID(ctx, p.CompanyID)
	if err != nil {
		if errors.Is(err, ErrCompanyNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "no company found")
		}
		return err
	}

	name, err := hp.repository.FetchByName(ctx, p.Name)
	if err != nil {
		if !errors.Is(err, ErrCompanyNotFound) {
			return err
		}
	}
	if name != nil && name.ID != ogMongoModel.ID {
		return echo.NewHTTPError(http.StatusBadRequest, "a company by that name already exists")
	}

	err = hp.repository.Patch(ctx, p.CompanyID, &MongoModel{
		Name:        p.Name,
		Description: p.Description,
		Employees:   p.Employees,
		Registered:  p.Registered,
		Type:        p.Type,
	})
	if err != nil {
		return err
	}
	return nil
}
