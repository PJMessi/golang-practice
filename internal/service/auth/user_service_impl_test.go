package auth

import (
	"context"
	"fmt"
	"testing"

	"github.com/pjmessi/golang-practice/internal/errorcode"
	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/internal/pkg/database"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/hash"
	"github.com/pjmessi/golang-practice/pkg/jwt"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// service and dependencies
var (
	authService    Service
	dbMock         database.DbMockImpl
	hashUtilMock   hash.UtilMockImpl
	jwtUtilMock    jwt.UtilMockImpl
	loggerUtilMock logger.UtilMockImpl
)

// setupMocks reset dependencies mocks
func setupMocks() {
	dbMock = database.DbMockImpl{}
	hashUtilMock = hash.UtilMockImpl{}
	jwtUtilMock = jwt.UtilMockImpl{}
	loggerUtilMock = logger.UtilMockImpl{}
	authService = NewService(&loggerUtilMock, &jwtUtilMock, &dbMock, &hashUtilMock)
}

var (
	ctx      context.Context
	email    string
	password string

	user   model.User
	jwtStr string
)

func setupVars() {
	email = "prajwalshrestha@test.com"
	password = "i_love_golang_3000"
	ctx = context.Background()

	hashedPw := fmt.Sprintf("hashed_%s", password)
	user = model.User{Email: email, Password: &hashedPw}
	jwtStr = "jwt_for_prajwalshrestha@test.com"
}

func Test_Login_User_Doesnt_exist(t *testing.T) {
	// ARRANGE
	setupMocks()
	setupVars()

	loggerUtilMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("GetUserByEmail", ctx, email).Return(false, model.User{}, nil)

	// ACT
	_, _, errRes := authService.Login(ctx, email, password)

	// ASSERT
	expectedErrStr := exception.NewUnauthenticated().Error()
	expectedLogStr := fmt.Sprintf("user with the email '%s' does not exist", email)
	assert.EqualError(t, errRes, expectedErrStr)
	loggerUtilMock.AssertCalled(t, "DebugCtx", ctx, expectedLogStr)
}

func Test_Login_Err_Getting_User_By_Email(t *testing.T) {
	// ARRANGE
	setupMocks()
	setupVars()

	getUserByEmailErr := fmt.Errorf("error from GetUserByEmail")
	loggerUtilMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("GetUserByEmail", ctx, email).Return(false, model.User{}, getUserByEmailErr)

	// ACT
	_, _, errRes := authService.Login(ctx, email, password)

	// ASSERT
	assert.EqualError(t, errRes, getUserByEmailErr.Error())
}

func Test_Login_User_Hasnt_Setup_Pw(t *testing.T) {
	// ARRANGE
	setupMocks()
	setupVars()

	user.Password = nil
	loggerUtilMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("GetUserByEmail", ctx, email).Return(true, user, nil)

	// ACT
	_, _, errRes := authService.Login(ctx, email, password)

	// ASSERT
	expectedErrStr := exception.NewUnauthenticatedFromBase(exception.Base{Type: errorcode.UserPwNotSet}).Error()
	expectedLogStr := fmt.Sprintf("user with the email '%s' hasn't setup his password", email)
	assert.EqualError(t, errRes, expectedErrStr)
	loggerUtilMock.AssertCalled(t, "DebugCtx", ctx, expectedLogStr)
}

func Test_Login_Incorrect_Pw(t *testing.T) {
	// ARRANGE
	setupMocks()
	setupVars()

	loggerUtilMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("GetUserByEmail", ctx, email).Return(true, user, nil)
	hashUtilMock.On("VerifyHash", *user.Password, password).Return(false)

	// ACT
	_, _, errRes := authService.Login(ctx, email, password)

	// ASSERT
	expectedErrStr := exception.NewUnauthenticated().Error()
	expectedLogStr := fmt.Sprintf("user with the email '%s' did not provide correct password", email)
	assert.EqualError(t, errRes, expectedErrStr)
	loggerUtilMock.AssertCalled(t, "DebugCtx", ctx, expectedLogStr)
}

func Test_Login_Success_Res(t *testing.T) {
	// ARRANGE
	setupMocks()
	setupVars()

	loggerUtilMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("GetUserByEmail", ctx, email).Return(true, user, nil)
	hashUtilMock.On("VerifyHash", *user.Password, password).Return(true)
	jwtUtilMock.On("Generate", user.Id, user.Email).Return(jwtStr, nil)

	// ACT
	userRes, jwtStrRes, errRes := authService.Login(ctx, email, password)

	// ASSERT
	assert.Equal(t, errRes, nil)
	assert.Equal(t, jwtStrRes, jwtStr)
	assert.Equal(t, userRes, user)
}

func Test_Login_Err_Generating_Jwt(t *testing.T) {
	// ARRANGE
	setupMocks()
	setupVars()

	generateErr := fmt.Errorf("error from Generate")

	loggerUtilMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("GetUserByEmail", ctx, email).Return(true, user, nil)
	hashUtilMock.On("VerifyHash", *user.Password, password).Return(true)
	jwtUtilMock.On("Generate", user.Id, user.Email).Return("", generateErr)

	// ACT
	_, _, errRes := authService.Login(ctx, email, password)

	// ASSERT
	assert.EqualError(t, errRes, generateErr.Error())
}
