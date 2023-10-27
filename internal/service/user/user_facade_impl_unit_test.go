package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pjmessi/golang-practice/config"
	"testing"
	"time"

	"github.com/pjmessi/golang-practice/internal/errorcode"
	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/internal/pkg/testutil"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/nats"
	"github.com/pjmessi/golang-practice/pkg/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupMocksForFacadeImplTest creates ServiceImpl with mocked dependencies
func setupMocksForFacadeImplTest() (*FacadeImpl, *ServiceMock, *logger.ServiceMock, *validation.HandlerMock, *nats.PubServiceMock, string) {
	userService := new(ServiceMock)
	validationUtilMock := new(validation.HandlerMock)
	logServiceMock := new(logger.ServiceMock)
	natsService := new(nats.PubServiceMock)
	userRegEvent := "EVENT.USER.NEW"
	authFacade := &FacadeImpl{
		userService:       userService,
		logService:        logServiceMock,
		validationHandler: validationUtilMock,
		natsService:       natsService,
		userRegEvent:      userRegEvent,
	}
	return authFacade, userService, logServiceMock, validationUtilMock, natsService, userRegEvent
}

// setupMocksForNewService returns mocked dependencies for NewService func
func setupMocksForNewFacade() (config.AppConfig, *logger.ServiceMock, *ServiceMock, *validation.HandlerMock, *nats.PubServiceMock) {
	validationUtilMock := new(validation.HandlerMock)
	logServiceMock := new(logger.ServiceMock)
	userServiceMock := new(ServiceMock)
	natsServiceMock := new(nats.PubServiceMock)
	appConfigMock := testutil.GetMockAppConfig(nil)
	return appConfigMock, logServiceMock, userServiceMock, validationUtilMock, natsServiceMock
}

func Test_NewFacade(t *testing.T) {
	// ARRANGE
	appConfig, logServiceMock, serviceMock, validationUtilMock, natsServiceMock := setupMocksForNewFacade()

	// ACT
	res := NewFacade(&appConfig, logServiceMock, serviceMock, validationUtilMock, natsServiceMock)

	// ASSERT
	resServiceImpl := res.(*FacadeImpl)

	assert.IsType(t, &FacadeImpl{}, res)
	assert.Equal(t, resServiceImpl, res)
}

func Test_Facade_RegisterUser_Invalid_Struct_In_Req_Bytes(t *testing.T) {
	// ARRANGE
	facade, _, _, _, _, _ := setupMocksForFacadeImplTest()

	ctx := context.Background()
	var reqByte []byte

	// ACT
	bytesRes, errRes := facade.RegisterUser(ctx, reqByte)

	// ASSERT
	expectedErr := exception.NewInvalidReqFromBase(exception.Base{Message: errorcode.ReqDataMissing})

	assert.Equal(t, expectedErr, errRes)
	assert.Nil(t, bytesRes)
}

func Test_Facade_RegisterUser_Invalid_Struct_Data_In_Req_Bytes(t *testing.T) {
	// ARRANGE
	facade, _, _, validationUtilMock, _, _ := setupMocksForFacadeImplTest()

	ctx := context.Background()
	regUserApiReq := testutil.GenMockRegUserApiReq(&model.UserRegApiReq{Email: "invalid_format"})
	reqBytes, _ := json.Marshal(regUserApiReq)
	validationErrDetails := map[string]string{}
	validationErrDetails["email"] = "invalid email"
	validationErr := validation.ValidationError{Details: validationErrDetails}

	validationUtilMock.On("ValidateStruct", regUserApiReq).Return(validationErr)

	// ACT
	bytesRes, errRes := facade.RegisterUser(ctx, reqBytes)

	// ASSERT
	expectedErr := exception.NewInvalidReqFromBase(exception.Base{Details: &validationErrDetails})

	assert.Equal(t, errRes, expectedErr)
	assert.Nil(t, bytesRes)
}

func Test_Facade_RegisterUser_Error_While_Validating_Req_Bytes(t *testing.T) {
	// ARRANGE
	facade, _, _, validationUtilMock, _, _ := setupMocksForFacadeImplTest()

	ctx := context.Background()
	regUserApiReq := testutil.GenMockRegUserApiReq(&model.UserRegApiReq{Email: "invalid_format"})
	reqBytes, _ := json.Marshal(regUserApiReq)
	validationErr := fmt.Errorf("error from validator.ValidateStruct")

	validationUtilMock.On("ValidateStruct", regUserApiReq).Return(validationErr)

	// ACT
	bytesRes, errRes := facade.RegisterUser(ctx, reqBytes)

	// ASSERT
	expectedErr := validationErr

	assert.Equal(t, errRes, expectedErr)
	assert.Nil(t, bytesRes)
}

func Test_Facade_RegisterUser_Err_Creating_User(t *testing.T) {
	// ARRANGE
	facade, service, _, validationUtilMock, _, _ := setupMocksForFacadeImplTest()

	ctx := context.Background()
	regUserApiReq := testutil.GenMockRegUserApiReq(nil)
	reqBytes, _ := json.Marshal(regUserApiReq)
	createUserErr := fmt.Errorf("error from CreateUser")

	validationUtilMock.On("ValidateStruct", regUserApiReq).Return(nil)
	service.On("CreateUser", ctx, regUserApiReq.Email, regUserApiReq.Password).Return(model.User{}, createUserErr)

	// ACT
	bytesRes, errRes := facade.RegisterUser(ctx, reqBytes)

	// ASSERT
	expectedErr := createUserErr

	assert.Equal(t, errRes, expectedErr)
	assert.Nil(t, bytesRes)
}

func Test_Facade_RegisterUser_Err_Publishing_Event(t *testing.T) {
	// ARRANGE
	facade, service, logServiceMock, validationUtilMock, natsServiceMock, userRegistrationEvent := setupMocksForFacadeImplTest()

	email := testutil.Fake.Internet().Email()
	ctx := context.Background()
	regUserApiReq := testutil.GenMockRegUserApiReq(&model.UserRegApiReq{Email: email})
	reqBytes, _ := json.Marshal(regUserApiReq)
	user := testutil.GenMockUser(&model.User{Email: email})
	publishErr := fmt.Errorf("error from Publish")

	logServiceMock.On("ErrorCtx", mock.Anything, mock.Anything)
	validationUtilMock.On("ValidateStruct", regUserApiReq).Return(nil)
	service.On("CreateUser", ctx, regUserApiReq.Email, regUserApiReq.Password).Return(user, nil)
	natsServiceMock.On("Publish", userRegistrationEvent, mock.Anything).Return(publishErr)

	// ACT
	_, errRes := facade.RegisterUser(ctx, reqBytes)

	// ASSERT
	// should log instead of returning error
	expectedLogStr := fmt.Sprintf("error publishing 'nats.user.new_registration' nats for userId '%s' and email '%s': %s", user.Id, user.Email, publishErr)

	assert.Equal(t, errRes, nil)
	logServiceMock.AssertCalled(t, "ErrorCtx", ctx, expectedLogStr)
}

func Test_Facade_RegisterUser_Success_Res(t *testing.T) {
	// ARRANGE
	facade, service, logServiceMock, validationUtilMock, natsServiceMock, userRegistrationEvent := setupMocksForFacadeImplTest()

	email := testutil.Fake.Internet().Email()
	ctx := context.Background()
	regUserApiReq := testutil.GenMockRegUserApiReq(&model.UserRegApiReq{Email: email})
	reqBytes, _ := json.Marshal(regUserApiReq)
	user := testutil.GenMockUser(&model.User{Email: email})

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	validationUtilMock.On("ValidateStruct", regUserApiReq).Return(nil)
	service.On("CreateUser", ctx, regUserApiReq.Email, regUserApiReq.Password).Return(user, nil)
	natsServiceMock.On("Publish", userRegistrationEvent, mock.Anything).Return(nil)

	// ACT
	bytesRes, errRes := facade.RegisterUser(ctx, reqBytes)

	// ASSERT
	expectedUserRes := model.UserRes{
		Id:        user.Id,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
	expectedRegUserApiRes := model.UserRegApiRes{User: expectedUserRes}
	expectedResByte, _ := json.Marshal(expectedRegUserApiRes)

	assert.Equal(t, errRes, nil)
	assert.Equal(t, bytesRes, expectedResByte)
}
