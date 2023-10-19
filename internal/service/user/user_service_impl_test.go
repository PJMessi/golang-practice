package user

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/pjmessi/golang-practice/internal/errorcode"
	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/internal/pkg/database"
	"github.com/pjmessi/golang-practice/internal/pkg/testutil"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupMocksForServiceImplTest creates ServiceImpl with mocked dependencies
func setupMocksForServiceImplTest() (*ServiceImpl, *database.DbMock, *logger.ServiceMock) {
	dbMock := new(database.DbMock)
	logServiceMock := new(logger.ServiceMock)
	service := &ServiceImpl{
		db:         dbMock,
		logService: logServiceMock,
	}
	return service, dbMock, logServiceMock
}

// setupMocksForNewService returns mocked dependencies for NewService func
func setupMocksForNewService() (*database.DbMock, *logger.ServiceMock) {
	dbMock := new(database.DbMock)
	logServiceMock := new(logger.ServiceMock)
	return dbMock, logServiceMock
}

func Test_NewService(t *testing.T) {
	// ARRANGE
	dbMock, logServiceMock := setupMocksForNewService()

	// ACT
	res := NewService(logServiceMock, dbMock)

	// ARRANGE
	resServiceImpl := res.(*ServiceImpl)

	assert.IsType(t, &ServiceImpl{}, res)
	assert.Equal(t, res, resServiceImpl)
}

func Test_CreateUser_Email_Already_Taken(t *testing.T) {
	// ARRANGE
	service, dbMock, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := "Password123!"

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("IsUserEmailTaken", ctx, email).Return(true, nil)

	// ACT
	userRes, errRes := service.CreateUser(ctx, email, password)

	// ASSERT
	expectedLogStr := fmt.Sprintf("user with the email '%s' already exists", email)
	expectedUserRes := model.User{}
	expectedErrRes := exception.NewAlreadyExistsFromBase(exception.Base{
		Message: fmt.Sprintf("user with the email '%s' already exists", email),
		Type:    errorcode.UserAlreadyExist,
	})

	assert.Equal(t, expectedUserRes, userRes)
	assert.Equal(t, expectedErrRes, errRes)
	logServiceMock.AssertCalled(t, "DebugCtx", ctx, expectedLogStr)
}

func Test_CreateUser_Error_Checking_If_Email_Taken(t *testing.T) {
	// ARRANGE
	service, dbMock, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := "Password123!"
	errIsUserEmailTaken := fmt.Errorf("error from IsUserEmailTaken")

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("IsUserEmailTaken", ctx, email).Return(false, errIsUserEmailTaken)

	// ACT
	userRes, errRes := service.CreateUser(ctx, email, password)

	// ASSERT
	expectedUserRes := model.User{}
	expectedErrRes := errIsUserEmailTaken

	assert.Equal(t, expectedUserRes, userRes)
	assert.Equal(t, expectedErrRes, errRes)
}

func Test_CreateUser_Weak_Password(t *testing.T) {
	// ARRANGE
	service, dbMock, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := "weakpw"

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("IsUserEmailTaken", ctx, email).Return(false, nil)

	// ACT
	userRes, errRes := service.CreateUser(ctx, email, password)

	// ASSERT
	expectedLogStr := "user did not provide strong password"
	expectedUserRes := model.User{}
	expectedErrRes := exception.NewInvalidReqFromBase(exception.Base{
		Message: "password is not strong enough",
	})

	assert.Equal(t, expectedUserRes, userRes)
	assert.Equal(t, expectedErrRes, errRes)
	logServiceMock.AssertCalled(t, "DebugCtx", ctx, expectedLogStr)
}

func Test_CreateUser_Should_Save_User_In_Db(t *testing.T) {
	// ARRANGE
	service, dbMock, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := "Password123!"

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("IsUserEmailTaken", ctx, email).Return(false, nil)
	dbMock.On("SaveUser", ctx, mock.Anything).Return(nil)

	// ACT
	service.CreateUser(ctx, email, password)

	// ASSERT
	dbMock.AssertCalled(t, "SaveUser", ctx, mock.MatchedBy(func(user *model.User) bool {
		return user.Email == strings.ToLower(email) &&
			*user.Password != password && // should be hashed password
			assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second) &&
			user.FirstName == nil &&
			user.LastName == nil &&
			user.UpdatedAt == nil
	}))
}

func Test_CreateUser_Error_Saving_User_In_Db(t *testing.T) {
	// ARRANGE
	service, dbMock, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := "Password123!"
	errCreateUser := fmt.Errorf("error from CreateUser")

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("IsUserEmailTaken", ctx, email).Return(false, nil)
	dbMock.On("SaveUser", ctx, mock.Anything).Return(errCreateUser)

	// ACT
	userRes, errRes := service.CreateUser(ctx, email, password)

	// ASSERT
	expectedUserRes := model.User{}
	expectedErrRes := errCreateUser

	assert.Equal(t, expectedUserRes, userRes)
	assert.Equal(t, expectedErrRes, errRes)
}

func Test_CreateUser_Success_Res(t *testing.T) {
	// ARRANGE
	service, dbMock, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := "Password123!"

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("IsUserEmailTaken", ctx, email).Return(false, nil)
	dbMock.On("SaveUser", ctx, mock.Anything).Return(nil)

	// ACT
	userRes, errRes := service.CreateUser(ctx, email, password)

	// ASSERT
	assert.Equal(t, email, userRes.Email)
	assert.NotEqual(t, &password, userRes.Password) // password should be hashed
	assert.WithinDuration(t, time.Now(), userRes.CreatedAt, time.Second)
	assert.Nil(t, userRes.UpdatedAt)
	assert.Nil(t, userRes.FirstName)
	assert.Nil(t, userRes.LastName)
	assert.Equal(t, nil, errRes)
}
