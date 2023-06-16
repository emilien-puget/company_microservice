package company

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type HandlerPost struct {
	repository interface {
		Create(ctx context.Context, company *MongoModel) error
		FetchByName(ctx context.Context, name string) (*MongoModel, error)
	}

	validate interface {
		StructCtx(ctx context.Context, s interface{}) (err error)
	}
}

func NewHandlerPost(repository *Repository, validate *validator.Validate) *HandlerPost {
	return &HandlerPost{
		repository: repository,
		validate:   validate,
	}
}

type postInput struct {
	Name        string      `json:"name" validate:"required,max=15"`
	Description string      `json:"description,omitempty" validate:"max=3000"`
	Employees   int         `json:"employees" validate:"required"`
	Registered  bool        `json:"registered"`
	Type        CompanyType `json:"type" validate:"required"`
}

type postOutput struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Employees   int         `json:"employees"`
	Registered  bool        `json:"registered"`
	Type        CompanyType `json:"type"`
}

func (hp HandlerPost) Handle(c echo.Context) error {
	companyInput := postInput{}
	err := c.Bind(&companyInput)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid parameters")
	}
	ctx := c.Request().Context()

	err = hp.validate.StructCtx(ctx, companyInput)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if !companyInput.Type.IsValid() {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid company type")
	}

	existingName, err := hp.repository.FetchByName(ctx, companyInput.Name)
	if err != nil {
		if !errors.Is(err, ErrCompanyNotFound) {
			return err
		}
	}
	if existingName != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "company already exists by name")
	}
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	err = hp.repository.Create(ctx, &MongoModel{
		ID:          newUUID.String(),
		Name:        companyInput.Name,
		Description: companyInput.Description,
		Employees:   companyInput.Employees,
		Registered:  companyInput.Registered,
		Type:        companyInput.Type,
	})
	if err != nil {
		return err
	}
	_ = c.JSON(http.StatusOK, postOutput{
		ID:          newUUID,
		Name:        companyInput.Name,
		Description: companyInput.Description,
		Employees:   companyInput.Employees,
		Registered:  companyInput.Registered,
		Type:        companyInput.Type,
	})
	return nil
}
