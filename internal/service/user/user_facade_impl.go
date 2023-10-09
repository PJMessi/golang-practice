package user

import (
	"context"
	"errors"

	"github.com/pjmessi/golang-practice/internal/dto"
	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/structutil"
	"github.com/pjmessi/golang-practice/pkg/validation"
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
		var valErr validation.ValidationError
		if errors.As(err, &valErr) {
			return nil, exception.NewInvalidReqFromBase(exception.Base{
				Details: &valErr.Details,
			})
		} else {
			return nil, err
		}
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
