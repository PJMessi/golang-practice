package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/pjmessi/golang-practice/internal/dto"
	"github.com/pjmessi/golang-practice/internal/errorcode"
	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/pkg/event"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/structutil"
	"github.com/pjmessi/golang-practice/pkg/validation"
)

type FacadeImpl struct {
	userService     Service
	validationUtil  validation.Util
	loggerUtil      logger.Util
	eventPubService event.PubService
}

func NewFacade(loggerUtil logger.Util, userService Service, validationUtil validation.Util, eventPubService event.PubService) Facade {
	return &FacadeImpl{
		userService:     userService,
		validationUtil:  validationUtil,
		loggerUtil:      loggerUtil,
		eventPubService: eventPubService,
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
		f.loggerUtil.ErrorCtx(ctx, fmt.Sprintf("error generating payload for 'event.user.new_registration' event for userId '%s' and email '%s': %s", userId, email, err))
		return nil
	}

	return eventPayloadBytes
}

func (f *FacadeImpl) publishNewRegEventPayload(ctx context.Context, payload []byte, userId string, email string) {
	err := f.eventPubService.Publish("event.user.new_registration", payload)
	if err != nil {
		f.loggerUtil.ErrorCtx(ctx, fmt.Sprintf("error publishing 'event.user.new_registration' event for userId '%s' and email '%s': %s", userId, email, err))
	} else {
		f.loggerUtil.DebugCtx(ctx, fmt.Sprintf("published 'event.user.new_registration' event for userId '%s' and email '%s'", userId, email))
	}
}
