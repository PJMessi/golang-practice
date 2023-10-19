package auth

import (
	"context"
	"errors"

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
	logService     logger.Service
}

func NewFacade(logService logger.Service, authService Service, validationUtil validation.Util) Facade {
	return &FacadeImpl{
		authService:    authService,
		validationUtil: validationUtil,
		logService:     logService,
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
		var valErr validation.ValidationError
		if errors.As(err, &valErr) {
			return nil, exception.NewInvalidReqFromBase(exception.Base{
				Details: &valErr.Details,
			})
		} else {
			return nil, err
		}
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
