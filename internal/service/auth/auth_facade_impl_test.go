package auth

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
func setupMocksForFacadeImplTest() (*FacadeImpl, *ServiceMock, *logger.ServiceMock, *validation.UtilMock) {
	authService := new(ServiceMock)
	validationUtilMock := new(validation.UtilMock)
	logServiceMock := new(logger.ServiceMock)
	authFacade := &FacadeImpl{
		authService:    authService,
		logService:     logServiceMock,
		validationUtil: validationUtilMock,
	}
	return authFacade, authService, logServiceMock, validationUtilMock
}

// setupMocksForNewService returns mocked dependencies for NewService func
func setupMocksForNewFacade() (*logger.ServiceMock, *ServiceMock, *validation.UtilMock) {
	validationUtilMock := new(validation.UtilMock)
	logServiceMock := new(logger.ServiceMock)
	authServiceMock := new(ServiceMock)
	return logServiceMock, authServiceMock, validationUtilMock
}

func Test_NewFacade(t *testing.T) {
	// ARRANGE
	logServiceMock, serviceMock, validatonUtilMock := setupMocksForNewFacade()

	// ACT
	res := NewFacade(logServiceMock, serviceMock, validatonUtilMock)

	// ARRANGE
	resServiceImpl := res.(*FacadeImpl)

	assert.IsType(t, &FacadeImpl{}, res)
	assert.Equal(t, res, resServiceImpl)
}

func Test_Facade_Login_Invalid_Struct_In_Req_Bytes(t *testing.T) {
	// ARRANGE
	facade, _, _, _ := setupMocksForFacadeImplTest()

	ctx := context.Background()
	reqByte := []byte{}

	// ACT
	bytesRes, errRes := facade.Login(ctx, reqByte)

	// ARRANGE
	expectedErr := exception.NewInvalidReqFromBase(exception.Base{Message: errorcode.ReqDataMissing})

	assert.Equal(t, errRes, expectedErr)
	assert.Nil(t, bytesRes)
}

func Test_Facade_Login_Invalid_Struct_Data_In_Req_Bytes(t *testing.T) {
	// ARRANGE
	facade, _, _, validationUtilMock := setupMocksForFacadeImplTest()

	ctx := context.Background()
	loginApiReq := testutil.GenLoginApiReq(&model.LoginApiReq{Email: "invalidformat"})
	reqBytes, _ := json.Marshal(loginApiReq)
	validationErrDetails := map[string]string{}
	validationErrDetails["email"] = "invalid email"
	validationErr := validation.ValidationError{Details: validationErrDetails}

	validationUtilMock.On("ValidateStruct", loginApiReq).Return(validationErr)

	// ACT
	bytesRes, errRes := facade.Login(ctx, reqBytes)

	// ARRANGE
	expectedErr := exception.NewInvalidReqFromBase(exception.Base{Details: &validationErrDetails})

	assert.Equal(t, errRes, expectedErr)
	assert.Nil(t, bytesRes)
}

func Test_Facade_Login_Error_While_Validating_Req_Bytes(t *testing.T) {
	// ARRANGE
	facade, _, _, validationUtilMock := setupMocksForFacadeImplTest()

	ctx := context.Background()
	loginApiReq := testutil.GenLoginApiReq(&model.LoginApiReq{Email: "invalidformat"})
	reqBytes, _ := json.Marshal(loginApiReq)
	validationErr := fmt.Errorf("error from validator.ValidateStruct")

	validationUtilMock.On("ValidateStruct", loginApiReq).Return(validationErr)

	// ACT
	bytesRes, errRes := facade.Login(ctx, reqBytes)

	// ARRANGE
	expectedErr := validationErr

	assert.Equal(t, errRes, expectedErr)
	assert.Nil(t, bytesRes)
}

func Test_Facade_Login_Err_Logging_In(t *testing.T) {
	// ARRANGE
	facade, service, _, validationUtilMock := setupMocksForFacadeImplTest()

	ctx := context.Background()
	loginApiReq := testutil.GenLoginApiReq(nil)
	reqBytes, _ := json.Marshal(loginApiReq)
	loginErr := fmt.Errorf("error from Login")

	validationUtilMock.On("ValidateStruct", loginApiReq).Return(nil)
	service.On("Login", ctx, loginApiReq.Email, loginApiReq.Password).Return(model.User{}, "", loginErr)

	// ACT
	bytesRes, errRes := facade.Login(ctx, reqBytes)

	// ARRANGE
	expectedErr := loginErr

	assert.Equal(t, errRes, expectedErr)
	assert.Nil(t, bytesRes)
}

func Test_Facade_Login_Success_Res(t *testing.T) {
	// ARRANGE
	facade, service, _, validationUtilMock := setupMocksForFacadeImplTest()

	email := testutil.Fake.Internet().Email()
	ctx := context.Background()
	loginApiReq := testutil.GenLoginApiReq(&model.LoginApiReq{Email: email})
	reqBytes, _ := json.Marshal(loginApiReq)
	user := testutil.GenMockUser(&model.User{Email: email})
	jwtStr := testutil.Fake.RandomStringWithLength(255)

	validationUtilMock.On("ValidateStruct", loginApiReq).Return(nil)
	service.On("Login", ctx, loginApiReq.Email, loginApiReq.Password).Return(user, jwtStr, nil)

	// ACT
	bytesRes, errRes := facade.Login(ctx, reqBytes)

	// ARRANGE
	expectedUserRes := model.UserRes{
		Id:        user.Id,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
	expectedLoginApiRes := model.LoginApiRes{User: expectedUserRes, Jwt: jwtStr}
	expectedResByte, _ := json.Marshal(expectedLoginApiRes)

	assert.Equal(t, errRes, nil)
	assert.Equal(t, bytesRes, expectedResByte)
}
