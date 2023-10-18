package user

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/pjmessi/golang-practice/internal/errorcode"
	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/internal/pkg/testutil"
	"github.com/pjmessi/golang-practice/pkg/event"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupMocksForFacadeImplTest creates ServiceImpl with mocked dependencies
func setupMocksForFacadeImplTest() (*FacadeImpl, *ServiceMock, *logger.UtilMock, *validation.UtilMock, *event.PubServiceMock) {
	userService := new(ServiceMock)
	validationUtilMock := new(validation.UtilMock)
	loggerUtilMock := new(logger.UtilMock)
	eventPubService := new(event.PubServiceMock)
	authFacade := &FacadeImpl{
		userService:     userService,
		loggerUtil:      loggerUtilMock,
		validationUtil:  validationUtilMock,
		eventPubService: eventPubService,
	}
	return authFacade, userService, loggerUtilMock, validationUtilMock, eventPubService
}

// setupMocksForNewService returns mocked dependencies for NewService func
func setupMocksForNewFacade() (*logger.UtilMock, *ServiceMock, *validation.UtilMock, *event.PubServiceMock) {
	validationUtilMock := new(validation.UtilMock)
	loggerUtilMock := new(logger.UtilMock)
	userService := new(ServiceMock)
	eventPubService := new(event.PubServiceMock)
	return loggerUtilMock, userService, validationUtilMock, eventPubService
}

func Test_NewFacade(t *testing.T) {
	// ARRANGE
	loggerUtilMock, serviceMock, validatonUtilMock, eventPubServiceMock := setupMocksForNewFacade()

	// ACT
	res := NewFacade(loggerUtilMock, serviceMock, validatonUtilMock, eventPubServiceMock)

	// ARRANGE
	resServiceImpl := res.(*FacadeImpl)

	assert.IsType(t, &FacadeImpl{}, res)
	assert.Equal(t, res, resServiceImpl)
}

func Test_Facade_RegisterUser_Invalid_Struct_In_Req_Bytes(t *testing.T) {
	// ARRANGE
	facade, _, _, _, _ := setupMocksForFacadeImplTest()

	ctx := context.Background()
	reqByte := []byte{}

	// ACT
	bytesRes, errRes := facade.RegisterUser(ctx, reqByte)

	// ARRANGE
	expectedErr := exception.NewInvalidReqFromBase(exception.Base{Message: errorcode.ReqDataMissing})

	assert.Equal(t, expectedErr, errRes)
	assert.Nil(t, bytesRes)
}

func Test_Facade_RegisterUser_Invalid_Struct_Data_In_Req_Bytes(t *testing.T) {
	// ARRANGE
	facade, _, _, validationUtilMock, _ := setupMocksForFacadeImplTest()

	ctx := context.Background()
	regUserApiReq := testutil.GenRegUserApiReq(&model.UserRegApiReq{Email: "invalidformat"})
	reqBytes, _ := json.Marshal(regUserApiReq)
	validationErrDetails := map[string]string{}
	validationErrDetails["email"] = "invalid email"
	validationErr := validation.ValidationError{Details: validationErrDetails}

	validationUtilMock.On("ValidateStruct", regUserApiReq).Return(validationErr)

	// ACT
	bytesRes, errRes := facade.RegisterUser(ctx, reqBytes)

	// ARRANGE
	expectedErr := exception.NewInvalidReqFromBase(exception.Base{Details: &validationErrDetails})

	assert.Equal(t, errRes, expectedErr)
	assert.Nil(t, bytesRes)
}

func Test_Facade_RegisterUser_Error_While_Validating_Req_Bytes(t *testing.T) {
	// ARRANGE
	facade, _, _, validationUtilMock, _ := setupMocksForFacadeImplTest()

	ctx := context.Background()
	regUserApiReq := testutil.GenRegUserApiReq(&model.UserRegApiReq{Email: "invalidformat"})
	reqBytes, _ := json.Marshal(regUserApiReq)
	validationErr := fmt.Errorf("error from validator.ValidateStruct")

	validationUtilMock.On("ValidateStruct", regUserApiReq).Return(validationErr)

	// ACT
	bytesRes, errRes := facade.RegisterUser(ctx, reqBytes)

	// ARRANGE
	expectedErr := validationErr

	assert.Equal(t, errRes, expectedErr)
	assert.Nil(t, bytesRes)
}

func Test_Facade_RegisterUser_Err_Creating_User(t *testing.T) {
	// ARRANGE
	facade, service, _, validationUtilMock, _ := setupMocksForFacadeImplTest()

	ctx := context.Background()
	regUserApiReq := testutil.GenRegUserApiReq(nil)
	reqBytes, _ := json.Marshal(regUserApiReq)
	createUserErr := fmt.Errorf("error from CreateUser")

	validationUtilMock.On("ValidateStruct", regUserApiReq).Return(nil)
	service.On("CreateUser", ctx, regUserApiReq.Email, regUserApiReq.Password).Return(model.User{}, createUserErr)

	// ACT
	bytesRes, errRes := facade.RegisterUser(ctx, reqBytes)

	// ARRANGE
	expectedErr := createUserErr

	assert.Equal(t, errRes, expectedErr)
	assert.Nil(t, bytesRes)
}

func Test_Facade_RegisterUser_Err_Publishing_Event(t *testing.T) {
	// ARRANGE
	facade, service, loggerUtilMock, validationUtilMock, eventPubServiceMock := setupMocksForFacadeImplTest()

	email := testutil.Fake.Internet().Email()
	ctx := context.Background()
	regUserApiReq := testutil.GenRegUserApiReq(&model.UserRegApiReq{Email: email})
	reqBytes, _ := json.Marshal(regUserApiReq)
	user := testutil.GenMockUser(&model.User{Email: email})
	publishErr := fmt.Errorf("Error from Publish")

	loggerUtilMock.On("ErrorCtx", mock.Anything, mock.Anything)
	validationUtilMock.On("ValidateStruct", regUserApiReq).Return(nil)
	service.On("CreateUser", ctx, regUserApiReq.Email, regUserApiReq.Password).Return(user, nil)
	eventPubServiceMock.On("Publish", "event.user.new_registration", mock.Anything).Return(publishErr)

	// ACT
	_, errRes := facade.RegisterUser(ctx, reqBytes)

	// ARRANGE
	// should log instead of returning error
	expectedLogStr := fmt.Sprintf("error publishing 'event.user.new_registration' event for userId '%s' and email '%s': %s", user.Id, user.Email, publishErr)

	assert.Equal(t, errRes, nil)
	loggerUtilMock.AssertCalled(t, "ErrorCtx", ctx, expectedLogStr)
}

func Test_Facade_RegisterUser_Success_Res(t *testing.T) {
	// ARRANGE
	facade, service, loggerUtilMock, validationUtilMock, eventPubServiceMock := setupMocksForFacadeImplTest()

	email := testutil.Fake.Internet().Email()
	ctx := context.Background()
	regUserApiReq := testutil.GenRegUserApiReq(&model.UserRegApiReq{Email: email})
	reqBytes, _ := json.Marshal(regUserApiReq)
	user := testutil.GenMockUser(&model.User{Email: email})

	loggerUtilMock.On("DebugCtx", mock.Anything, mock.Anything)
	validationUtilMock.On("ValidateStruct", regUserApiReq).Return(nil)
	service.On("CreateUser", ctx, regUserApiReq.Email, regUserApiReq.Password).Return(user, nil)
	eventPubServiceMock.On("Publish", "event.user.new_registration", mock.Anything).Return(nil)

	// ACT
	bytesRes, errRes := facade.RegisterUser(ctx, reqBytes)

	// ARRANGE
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
