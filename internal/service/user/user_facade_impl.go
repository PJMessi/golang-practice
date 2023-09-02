package user

import (
	"github.com/pjmessi/go-database-usage/internal/dto"
	"github.com/pjmessi/go-database-usage/internal/model"
	"github.com/pjmessi/go-database-usage/pkg/exception"
	"github.com/pjmessi/go-database-usage/pkg/structutil"
	"github.com/pjmessi/go-database-usage/pkg/validation"
)

type FacadeImpl struct {
	Facade
	userService    Service
	validationUtil validation.Util
}

func NewFacade(userService Service, validationUtil validation.Util) Facade {
	return &FacadeImpl{
		userService:    userService,
		validationUtil: validationUtil,
	}
}

func (f *FacadeImpl) RegisterUser(reqBytes []byte) ([]byte, error) {
	var req model.UserRegApiReq

	err := structutil.ConvertFromBytes(reqBytes, &req)
	if err != nil {
		return nil, exception.NewInvalidReqFromBase(exception.Base{
			Message: "missing request data",
		})
	}

	err = f.validationUtil.ValidateStruct(req)
	if err != nil {
		details := f.validationUtil.FormatValidationError(err)
		return nil, exception.NewInvalidReqFromBase(exception.Base{
			Details: &details,
		})
	}

	user, err := f.userService.CreateUser(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	res := model.UserRegApiRes{
		User: dto.UserToUserRes(&user),
	}

	return structutil.ConvertToBytes(res)
}
