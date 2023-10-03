package auth

import (
	"context"

	"github.com/pjmessi/go-database-usage/internal/dto"
	"github.com/pjmessi/go-database-usage/internal/model"
	"github.com/pjmessi/go-database-usage/pkg/exception"
	"github.com/pjmessi/go-database-usage/pkg/structutil"
	"github.com/pjmessi/go-database-usage/pkg/validation"
)

type FacadeImpl struct {
	authService    Service
	validationUtil validation.Util
}

func NewFacade(authService Service, validationUtil validation.Util) Facade {
	return &FacadeImpl{
		authService:    authService,
		validationUtil: validationUtil,
	}
}

func (f *FacadeImpl) Login(ctx context.Context, reqBytes []byte) ([]byte, error) {
	var req model.LoginApiReq

	err := structutil.ConvertFromBytes(reqBytes, &req)
	if err != nil {
		return nil, exception.NewInvalidReqFromBase(exception.Base{Message: "missing request data"})
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
