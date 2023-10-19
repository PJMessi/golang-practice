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
	"github.com/pjmessi/golang-practice/internal/pkg/password"
	"github.com/pjmessi/golang-practice/internal/pkg/testutil"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupMocksForServiceImplTest creates ServiceImpl with mocked dependencies
func setupMocksForServiceImplTest() (*ServiceImpl, *database.DbMock, *password.UtilMock, *logger.ServiceMock) {
	dbMock := new(database.DbMock)
	passwordUtil := new(password.UtilMock)
	logServiceMock := new(logger.ServiceMock)
	service := &ServiceImpl{
		db:           dbMock,
		passwordUtil: passwordUtil,
		logService:   logServiceMock,
	}
	return service, dbMock, passwordUtil, logServiceMock
}

// setupMocksForNewService returns mocked dependencies for NewService func
func setupMocksForNewService() (*database.DbMock, *password.UtilMock, *logger.ServiceMock) {
	dbMock := new(database.DbMock)
	logServiceMock := new(logger.ServiceMock)
	passwordUtilMock := new(password.UtilMock)
	return dbMock, passwordUtilMock, logServiceMock
}

func Test_NewService(t *testing.T) {
	// ARRANGE
	dbMock, passwordUtilMock, logServiceMock := setupMocksForNewService()

	// ACT
	res := NewService(logServiceMock, dbMock, passwordUtilMock)

	// ARRANGE
	resServiceImpl := res.(*ServiceImpl)

	assert.IsType(t, &ServiceImpl{}, res)
	assert.Equal(t, res, resServiceImpl)
}

func Test_CreateUser_Email_Already_Taken(t *testing.T) {
	// ARRANGE
	service, dbMock, passwordUtilMock, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := testutil.Fake.Internet().Password()

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	passwordUtilMock.On("IsStrong", password).Return(true)
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
	service, dbMock, passwordUtilMock, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := testutil.Fake.Internet().Password()
	errIsUserEmailTaken := fmt.Errorf("error from IsUserEmailTaken")

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	passwordUtilMock.On("IsStrong", password).Return(true)
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
	service, dbMock, passwordUtilMock, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := "password"

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("IsUserEmailTaken", ctx, email).Return(false, nil)
	passwordUtilMock.On("IsStrong", password).Return(false)

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

func Test_CreateUser_Error_Hashing_Password(t *testing.T) {
	// ARRANGE
	service, dbMock, passwordUtilMock, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := "password"
	errHash := fmt.Errorf("error from Hash")

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("IsUserEmailTaken", ctx, email).Return(false, nil)
	passwordUtilMock.On("IsStrong", password).Return(true)
	passwordUtilMock.On("Hash", password).Return("", errHash)

	// ACT
	userRes, errRes := service.CreateUser(ctx, email, password)

	// ASSERT
	expectedUserRes := model.User{}
	expectedErrRes := errHash

	assert.Equal(t, expectedUserRes, userRes)
	assert.Equal(t, expectedErrRes, errRes)
}

func Test_CreateUser_Should_Save_User_In_Db(t *testing.T) {
	// ARRANGE
	service, dbMock, passwordUtilMock, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := "password"
	passwordHash := testutil.Fake.RandomStringWithLength(100)

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	passwordUtilMock.On("IsStrong", password).Return(true)
	dbMock.On("IsUserEmailTaken", ctx, email).Return(false, nil)
	passwordUtilMock.On("Hash", password).Return(passwordHash, nil)
	dbMock.On("SaveUser", ctx, mock.Anything).Return(nil)

	// ACT
	service.CreateUser(ctx, email, password)

	// ASSERT
	dbMock.AssertCalled(t, "SaveUser", ctx, mock.MatchedBy(func(user *model.User) bool {
		return user.Email == strings.ToLower(email) &&
			*user.Password == passwordHash &&
			assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second) &&
			user.FirstName == nil &&
			user.LastName == nil &&
			user.UpdatedAt == nil
	}))
}

func Test_CreateUser_Error_Saving_User_In_Db(t *testing.T) {
	// ARRANGE
	service, dbMock, passwordUtilMock, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := "password"
	passwordHash := testutil.Fake.RandomStringWithLength(100)
	errCreateUser := fmt.Errorf("error from CreateUser")

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("IsUserEmailTaken", ctx, email).Return(false, nil)
	passwordUtilMock.On("IsStrong", password).Return(true)
	passwordUtilMock.On("Hash", password).Return(passwordHash, nil)
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
	service, dbMock, passwordUtilMock, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := "password"
	passwordHash := testutil.Fake.RandomStringWithLength(100)

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("IsUserEmailTaken", ctx, email).Return(false, nil)
	passwordUtilMock.On("IsStrong", password).Return(true)
	passwordUtilMock.On("Hash", password).Return(passwordHash, nil)
	dbMock.On("SaveUser", ctx, mock.Anything).Return(nil)

	// ACT
	userRes, errRes := service.CreateUser(ctx, email, password)

	// ASSERT
	assert.Equal(t, email, userRes.Email)
	assert.Equal(t, &passwordHash, userRes.Password)
	assert.WithinDuration(t, time.Now(), userRes.CreatedAt, time.Second)
	assert.Nil(t, userRes.UpdatedAt)
	assert.Nil(t, userRes.FirstName)
	assert.Nil(t, userRes.LastName)
	assert.Equal(t, nil, errRes)
}
