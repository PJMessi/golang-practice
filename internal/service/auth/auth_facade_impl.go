package auth

import (
	"context"

	"github.com/pjmessi/golang-practice/internal/dto"
	"github.com/pjmessi/golang-practice/internal/errorcode"
	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/structutil"
	"github.com/pjmessi/golang-practice/pkg/validation"
)

type FacadeImpl struct {
	authService    Service
	validationUtil validation.Util
	loggerUtil     logger.Util
}

func NewFacade(loggerUtil logger.Util, authService Service, validationUtil validation.Util) Facade {
	return &FacadeImpl{
		authService:    authService,
		validationUtil: validationUtil,
		loggerUtil:     loggerUtil,
	}
}

func (f *FacadeImpl) Login(ctx context.Context, reqBytes []byte) ([]byte, error) {
	var req model.LoginApiReq

	err := structutil.ConvertFromBytes(reqBytes, &req)
	if err != nil {
		return nil, exception.NewInvalidReqFromBase(exception.Base{Message: errorcode.ReqDataMissing})
	}

	err = f.validationUtil.ValidateStruct(req)
	if err != nil {
		details := f.validationUtil.FormatValidationError(err)
		return nil, exception.NewInvalidReqFromBase(exception.Base{
			Details: &details,
		})
	}

	user, jwtStr, err := f.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	res := model.LoginApiRes{
		User: dto.UserToUserRes(&user),
		Jwt:  jwtStr,
	}

	return structutil.ConvertToBytes(res)
}
