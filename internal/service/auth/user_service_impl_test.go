package auth

import (
	"context"
	"fmt"
	"testing"

	"github.com/pjmessi/golang-practice/internal/errorcode"
	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/internal/pkg/database"
	"github.com/pjmessi/golang-practice/internal/pkg/password"
	"github.com/pjmessi/golang-practice/internal/pkg/testutil"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/jwt"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupMocksForServiceImplTest creates ServiceImpl with mocked dependencies
func setupMocksForServiceImplTest() (*ServiceImpl, *database.DbMock, *password.UtilMock, *jwt.UtilMock, *logger.ServiceMock) {
	dbMock := new(database.DbMock)
	passwordUtil := new(password.UtilMock)
	jwtUtilMock := new(jwt.UtilMock)
	logServiceMock := new(logger.ServiceMock)
	service := &ServiceImpl{
		db:           dbMock,
		passwordUtil: passwordUtil,
		jwtHandler:   jwtUtilMock,
		logService:   logServiceMock,
	}
	return service, dbMock, passwordUtil, jwtUtilMock, logServiceMock
}

// setupMocksForNewService returns mocked dependencies for NewService func
func setupMocksForNewService() (*database.DbMock, *password.UtilMock, *jwt.UtilMock, *logger.ServiceMock) {
	dbMock := new(database.DbMock)
	passwordUtilMock := new(password.UtilMock)
	jwtUtilMock := new(jwt.UtilMock)
	logServiceMock := new(logger.ServiceMock)
	return dbMock, passwordUtilMock, jwtUtilMock, logServiceMock
}

func Test_NewService(t *testing.T) {
	// ARRANGE
	dbMock, passwordUtilMock, jwtUtilMock, logServiceMock := setupMocksForNewService()

	// ACT
	res := NewService(logServiceMock, jwtUtilMock, dbMock, passwordUtilMock)

	// ARRANGE
	resServiceImpl := res.(*ServiceImpl)

	assert.IsType(t, &ServiceImpl{}, res)
	assert.Equal(t, res, resServiceImpl)
}

func Test_Service_Login_User_Doesnt_exist(t *testing.T) {
	// ARRANGE
	service, dbMock, _, _, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := testutil.Fake.Internet().Password()

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("GetUserByEmail", ctx, email).Return(false, model.User{}, nil)

	// ACT
	userRes, jwtStrRes, errRes := service.Login(ctx, email, password)

	// ASSERT
	expectedUserRes := model.User{}
	expectedJwtStrRes := ""
	expectedErr := exception.NewUnauthenticated()
	expectedLogStr := fmt.Sprintf("user with the email '%s' does not exist", email)

	assert.Equal(t, userRes, expectedUserRes)
	assert.Equal(t, jwtStrRes, expectedJwtStrRes)
	assert.Equal(t, errRes, expectedErr)
	logServiceMock.AssertCalled(t, "DebugCtx", ctx, expectedLogStr)
}

func Test_Service_Login_Err_Getting_User_By_Email(t *testing.T) {
	// ARRANGE
	service, dbMock, _, _, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := testutil.Fake.Internet().Password()

	getUserByEmailErr := fmt.Errorf("error from GetUserByEmail")
	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("GetUserByEmail", ctx, email).Return(false, model.User{}, getUserByEmailErr)

	// ACT
	userRes, jwtStrRes, errRes := service.Login(ctx, email, password)

	// ASSERT
	expectedUserRes := model.User{}
	expectedJwtStrRes := ""
	expectedErr := getUserByEmailErr

	assert.Equal(t, userRes, expectedUserRes)
	assert.Equal(t, jwtStrRes, expectedJwtStrRes)
	assert.Equal(t, errRes, expectedErr)
}

func Test_Service_Login_User_Hasnt_Setup_Pw(t *testing.T) {
	// ARRANGE
	service, dbMock, _, _, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := testutil.Fake.Internet().Password()
	user := testutil.GenMockUser(&model.User{Email: email})

	user.Password = nil
	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("GetUserByEmail", ctx, email).Return(true, user, nil)

	// ACT
	userRes, jwtStrRes, errRes := service.Login(ctx, email, password)

	// ASSERT
	expectedUserRes := model.User{}
	expectedJwtStrRes := ""
	expectedErr := exception.NewUnauthenticatedFromBase(exception.Base{Type: errorcode.UserPwNotSet})

	assert.Equal(t, userRes, expectedUserRes)
	assert.Equal(t, jwtStrRes, expectedJwtStrRes)
	assert.Equal(t, errRes, expectedErr)
}

func Test_Service_Login_Incorrect_Pw(t *testing.T) {
	// ARRANGE
	service, dbMock, passwordUtilMock, _, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := testutil.Fake.Internet().Password()
	user := testutil.GenMockUser(&model.User{Email: email})

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("GetUserByEmail", ctx, email).Return(true, user, nil)
	passwordUtilMock.On("IsHashCorrect", *user.Password, password).Return(false)

	// ACT
	userRes, jwtStrRes, errRes := service.Login(ctx, email, password)

	// ASSERT
	expectedUserRes := model.User{}
	expectedJwtStrRes := ""
	expectedErr := exception.NewUnauthenticated()

	assert.Equal(t, userRes, expectedUserRes)
	assert.Equal(t, jwtStrRes, expectedJwtStrRes)
	assert.Equal(t, errRes, expectedErr)
}

func Test_Service_Login_Success_Res(t *testing.T) {
	// ARRANGE
	service, dbMock, passwordUtilMock, jwtUtilMock, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := "PJMESSI25@icloud.com" // testutil.Fake.Internet().Email()
	password := testutil.Fake.Internet().Password()
	user := testutil.GenMockUser(&model.User{Email: email})
	jwtStr := testutil.Fake.RandomStringWithLength(100)

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("GetUserByEmail", ctx, email).Return(true, user, nil)
	passwordUtilMock.On("IsHashCorrect", *user.Password, password).Return(true)
	jwtUtilMock.On("Generate", user.Id, user.Email).Return(jwtStr, nil)

	// ACT
	userRes, jwtStrRes, errRes := service.Login(ctx, email, password)

	// ASSERT
	expectedUserRes := user
	expectedJwtStrRes := jwtStr

	assert.Equal(t, userRes, expectedUserRes)
	assert.Equal(t, jwtStrRes, expectedJwtStrRes)
	assert.Equal(t, errRes, nil)
}

func Test_Service_Login_Err_Generating_Jwt(t *testing.T) {
	// ARRANGE
	service, dbMock, passwordUtilMock, jwtUtilMock, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := testutil.Fake.Internet().Password()
	user := testutil.GenMockUser(&model.User{Email: email})
	generateErr := fmt.Errorf("error from Generate")

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("GetUserByEmail", ctx, email).Return(true, user, nil)
	passwordUtilMock.On("IsHashCorrect", *user.Password, password).Return(true)
	jwtUtilMock.On("Generate", user.Id, user.Email).Return("", generateErr)

	// ACT
	userRes, jwtStrRes, errRes := service.Login(ctx, email, password)

	// ASSERT
	expectedUserRes := model.User{}
	expectedJwtStrRes := ""
	expectedErr := generateErr

	assert.Equal(t, userRes, expectedUserRes)
	assert.Equal(t, jwtStrRes, expectedJwtStrRes)
	assert.Equal(t, errRes, expectedErr)
}
