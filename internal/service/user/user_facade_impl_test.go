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
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/validation"
	"github.com/stretchr/testify/assert"
)

// setupMocksForFacadeImplTest creates ServiceImpl with mocked dependencies
func setupMocksForFacadeImplTest() (*FacadeImpl, *ServiceMock, *logger.UtilMock, *validation.UtilMock) {
	userService := new(ServiceMock)
	validationUtilMock := new(validation.UtilMock)
	loggerUtilMock := new(logger.UtilMock)
	authFacade := &FacadeImpl{
		userService:    userService,
		loggerUtil:     loggerUtilMock,
		validationUtil: validationUtilMock,
	}
	return authFacade, userService, loggerUtilMock, validationUtilMock
}

// setupMocksForNewService returns mocked dependencies for NewService func
func setupMocksForNewFacade() (*logger.UtilMock, *ServiceMock, *validation.UtilMock) {
	validationUtilMock := new(validation.UtilMock)
	loggerUtilMock := new(logger.UtilMock)
	userService := new(ServiceMock)
	return loggerUtilMock, userService, validationUtilMock
}

func Test_NewFacade(t *testing.T) {
	// ARRANGE
	loggerUtilMock, serviceMock, validatonUtilMock := setupMocksForNewFacade()

	// ACT
	res := NewFacade(loggerUtilMock, serviceMock, validatonUtilMock)

	// ARRANGE
	resServiceImpl := res.(*FacadeImpl)

	assert.IsType(t, &FacadeImpl{}, res)
	assert.Equal(t, res, resServiceImpl)
}

func Test_Facade_RegisterUser_Invalid_Struct_In_Req_Bytes(t *testing.T) {
	// ARRANGE
	facade, _, _, _ := setupMocksForFacadeImplTest()

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
	facade, _, _, validationUtilMock := setupMocksForFacadeImplTest()

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
	facade, _, _, validationUtilMock := setupMocksForFacadeImplTest()

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

func Test_Facade_RegisterUser_Err_Registering_User(t *testing.T) {
	// ARRANGE
	facade, service, _, validationUtilMock := setupMocksForFacadeImplTest()

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

func Test_Facade_RegisterUser_Success_Res(t *testing.T) {
	// ARRANGE
	facade, service, _, validationUtilMock := setupMocksForFacadeImplTest()

	email := testutil.Fake.Internet().Email()
	ctx := context.Background()
	regUserApiReq := testutil.GenRegUserApiReq(&model.UserRegApiReq{Email: email})
	reqBytes, _ := json.Marshal(regUserApiReq)
	user := testutil.GenMockUser(&model.User{Email: email})

	validationUtilMock.On("ValidateStruct", regUserApiReq).Return(nil)
	service.On("CreateUser", ctx, regUserApiReq.Email, regUserApiReq.Password).Return(user, nil)

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
