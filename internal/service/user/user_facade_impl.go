package user

import (
	"context"

	"github.com/pjmessi/go-database-usage/internal/dto"
	"github.com/pjmessi/go-database-usage/internal/model"
	"github.com/pjmessi/go-database-usage/pkg/exception"
	"github.com/pjmessi/go-database-usage/pkg/logger"
	"github.com/pjmessi/go-database-usage/pkg/structutil"
	"github.com/pjmessi/go-database-usage/pkg/validation"
)

type FacadeImpl struct {
	Facade
	userService    Service
	validationUtil validation.Util
	loggerUtil     logger.Util
}

func NewFacade(loggerUtil logger.Util, userService Service, validationUtil validation.Util) Facade {
	return &FacadeImpl{
		userService:    userService,
		validationUtil: validationUtil,
		loggerUtil:     loggerUtil,
	}
}

func (f *FacadeImpl) RegisterUser(ctx context.Context, reqBytes []byte) ([]byte, error) {
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

	user, err := f.userService.CreateUser(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	res := model.UserRegApiRes{
		User: dto.UserToUserRes(&user),
	}

	return structutil.ConvertToBytes(res)
}