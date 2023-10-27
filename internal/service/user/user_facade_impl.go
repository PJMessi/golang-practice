package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/pjmessi/golang-practice/config"

	"github.com/pjmessi/golang-practice/internal/dto"
	"github.com/pjmessi/golang-practice/internal/errorcode"
	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/internal/pkg/jwt"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/nats"
	"github.com/pjmessi/golang-practice/pkg/structutil"
	"github.com/pjmessi/golang-practice/pkg/validation"
)

type FacadeImpl struct {
	userService       Service
	validationHandler validation.Handler
	logService        logger.Service
	natsService       nats.Service
	userRegEvent      string
}

func NewFacade(appConfig *config.AppConfig, logService logger.Service, userService Service, validationHandler validation.Handler, natsService nats.Service) Facade {
	return &FacadeImpl{
		userService:       userService,
		validationHandler: validationHandler,
		logService:        logService,
		natsService:       natsService,
		userRegEvent:      appConfig.NATS_EVENT_USER_REGISTRATION,
	}
}

func (f *FacadeImpl) RegisterUser(ctx context.Context, reqBytes []byte) ([]byte, error) {
	var req model.UserRegApiReq

	err := structutil.ConvertFromBytes(reqBytes, &req)
	if err != nil {
		return nil, exception.NewInvalidReqFromBase(exception.Base{
			Message: errorcode.ReqDataMissing,
		})
	}

	err = f.validationHandler.ValidateStruct(req)
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

	eventPayload := f.genNewRegEventPayload(ctx, user.Id, user.Email)
	if eventPayload != nil {
		f.publishNewRegEventPayload(ctx, eventPayload, user.Id, user.Email)
	}

	res := model.UserRegApiRes{
		User: dto.UserToUserRes(&user),
	}

	return structutil.ConvertToBytes(res)
}

func (f *FacadeImpl) genNewRegEventPayload(ctx context.Context, userId string, email string) []byte {
	eventPayload := map[string]string{
		"email": email,
		"id":    userId,
	}

	eventPayloadBytes, err := structutil.ConvertToBytes(eventPayload)
	if err != nil {
		f.logService.ErrorCtx(ctx, fmt.Sprintf("error generating payload for 'nats.user.new_registration' nats for userId '%s' and email '%s': %s", userId, email, err))
		return nil
	}

	return eventPayloadBytes
}

func (f *FacadeImpl) publishNewRegEventPayload(ctx context.Context, payload []byte, userId string, email string) {
	err := f.natsService.Publish(f.userRegEvent, payload)
	if err != nil {
		f.logService.ErrorCtx(ctx, fmt.Sprintf("error publishing 'nats.user.new_registration' nats for userId '%s' and email '%s': %s", userId, email, err))
	} else {
		f.logService.DebugCtx(ctx, fmt.Sprintf("published 'nats.user.new_registration' nats for userId '%s' and email '%s'", userId, email))
	}
}

func (f *FacadeImpl) GetProfile(ctx context.Context, reqBytes []byte, jwtPayload jwt.JwtPayload) ([]byte, error) {
	user, err := f.userService.GetProfile(ctx, jwtPayload.UserId)
	if err != nil {
		return nil, err
	}

	getProfileRes := model.GetProfileApiRes{}
	getProfileRes.User = dto.UserToUserRes(&user)

	getProfileResBytes, err := structutil.ConvertToBytes(getProfileRes)
	if err != nil {
		return nil, err
	}

	return getProfileResBytes, nil
}
