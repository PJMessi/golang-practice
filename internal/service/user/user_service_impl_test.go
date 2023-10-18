package user

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pjmessi/golang-practice/internal/errorcode"
	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/internal/pkg/database"
	"github.com/pjmessi/golang-practice/internal/pkg/password"
	"github.com/pjmessi/golang-practice/internal/pkg/testutil"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupMocksForServiceImplTest creates ServiceImpl with mocked dependencies
func setupMocksForServiceImplTest() (*ServiceImpl, *database.DbMock, *password.UtilMock, *logger.UtilMock, *uuid.UtilMock) {
	dbMock := new(database.DbMock)
	passwordUtil := new(password.UtilMock)
	loggerUtilMock := new(logger.UtilMock)
	uuidUtilMock := new(uuid.UtilMock)
	service := &ServiceImpl{
		db:           dbMock,
		passwordUtil: passwordUtil,
		loggerUtil:   loggerUtilMock,
		uuidUtil:     uuidUtilMock,
	}
	return service, dbMock, passwordUtil, loggerUtilMock, uuidUtilMock
}

// setupMocksForNewService returns mocked dependencies for NewService func
func setupMocksForNewService() (*database.DbMock, *password.UtilMock, *uuid.UtilMock, *logger.UtilMock) {
	dbMock := new(database.DbMock)
	loggerUtilMock := new(logger.UtilMock)
	passwordUtilMock := new(password.UtilMock)
	uuidUtilMock := new(uuid.UtilMock)
	return dbMock, passwordUtilMock, uuidUtilMock, loggerUtilMock
}

func Test_NewService(t *testing.T) {
	// ARRANGE
	dbMock, passwordUtilMock, uuidUtilMock, loggerUtilMock := setupMocksForNewService()

	// ACT
	res := NewService(loggerUtilMock, dbMock, passwordUtilMock, uuidUtilMock)

	// ARRANGE
	resServiceImpl := res.(*ServiceImpl)

	assert.IsType(t, &ServiceImpl{}, res)
	assert.Equal(t, res, resServiceImpl)
}

func Test_CreateUser_Email_Already_Taken(t *testing.T) {
	// ARRANGE
	service, dbMock, passwordUtilMock, loggerUtilMock, _ := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := testutil.Fake.Internet().Password()

	loggerUtilMock.On("DebugCtx", mock.Anything, mock.Anything)
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
	loggerUtilMock.AssertCalled(t, "DebugCtx", ctx, expectedLogStr)
}

func Test_CreateUser_Error_Checking_If_Email_Taken(t *testing.T) {
	// ARRANGE
	service, dbMock, passwordUtilMock, loggerUtilMock, _ := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := testutil.Fake.Internet().Password()
	errIsUserEmailTaken := fmt.Errorf("error from IsUserEmailTaken")

	loggerUtilMock.On("DebugCtx", mock.Anything, mock.Anything)
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
	service, dbMock, passwordUtilMock, loggerUtilMock, _ := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := "password"

	loggerUtilMock.On("DebugCtx", mock.Anything, mock.Anything)
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
	loggerUtilMock.AssertCalled(t, "DebugCtx", ctx, expectedLogStr)
}

func Test_CreateUser_Error_Hashing_Password(t *testing.T) {
	// ARRANGE
	service, dbMock, passwordUtilMock, loggerUtilMock, _ := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := "password"
	errHash := fmt.Errorf("error from Hash")

	loggerUtilMock.On("DebugCtx", mock.Anything, mock.Anything)
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

func Test_CreateUser_Error_Generating_UserId_Uuid(t *testing.T) {
	// ARRANGE
	service, dbMock, passwordUtilMock, loggerUtilMock, uuidUtilMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := "password"
	passwordHash := testutil.Fake.RandomStringWithLength(100)
	errGenUuidV4 := fmt.Errorf("error from GenUuidV4")

	loggerUtilMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("IsUserEmailTaken", ctx, email).Return(false, nil)
	passwordUtilMock.On("IsStrong", password).Return(true)
	passwordUtilMock.On("Hash", password).Return(passwordHash, nil)
	uuidUtilMock.On("GenUuidV4").Return("", errGenUuidV4)

	// ACT
	userRes, errRes := service.CreateUser(ctx, email, password)

	// ASSERT
	expectedUserRes := model.User{}
	expectedErrRes := errGenUuidV4

	assert.Equal(t, expectedUserRes, userRes)
	assert.Equal(t, expectedErrRes, errRes)
}

func Test_CreateUser_Should_Save_User_In_Db(t *testing.T) {
	// ARRANGE
	service, dbMock, passwordUtilMock, loggerUtilMock, uuidUtilMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := "password"
	passwordHash := testutil.Fake.RandomStringWithLength(100)
	uuidStr := testutil.Fake.UUID().V4()

	loggerUtilMock.On("DebugCtx", mock.Anything, mock.Anything)
	passwordUtilMock.On("IsStrong", password).Return(true)
	dbMock.On("IsUserEmailTaken", ctx, email).Return(false, nil)
	passwordUtilMock.On("Hash", password).Return(passwordHash, nil)
	uuidUtilMock.On("GenUuidV4").Return(uuidStr, nil)
	dbMock.On("SaveUser", ctx, mock.Anything).Return(nil)

	// ACT
	service.CreateUser(ctx, email, password)

	// ASSERT
	dbMock.AssertCalled(t, "SaveUser", ctx, mock.MatchedBy(func(user *model.User) bool {
		return user.Email == email &&
			*user.Password == passwordHash &&
			user.Id == uuidStr &&
			assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second) &&
			user.FirstName == nil &&
			user.LastName == nil &&
			user.UpdatedAt == nil
	}))
}

func Test_CreateUser_Error_Saving_User_In_Db(t *testing.T) {
	// ARRANGE
	service, dbMock, passwordUtilMock, loggerUtilMock, uuidUtilMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := "password"
	passwordHash := testutil.Fake.RandomStringWithLength(100)
	uuidStr := testutil.Fake.UUID().V4()
	errCreateUser := fmt.Errorf("error from CreateUser")

	loggerUtilMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("IsUserEmailTaken", ctx, email).Return(false, nil)
	passwordUtilMock.On("IsStrong", password).Return(true)
	passwordUtilMock.On("Hash", password).Return(passwordHash, nil)
	uuidUtilMock.On("GenUuidV4").Return(uuidStr, nil)
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
	service, dbMock, passwordUtilMock, loggerUtilMock, uuidUtilMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := "password"
	passwordHash := testutil.Fake.RandomStringWithLength(100)
	uuidStr := testutil.Fake.UUID().V4()

	loggerUtilMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("IsUserEmailTaken", ctx, email).Return(false, nil)
	passwordUtilMock.On("IsStrong", password).Return(true)
	passwordUtilMock.On("Hash", password).Return(passwordHash, nil)
	uuidUtilMock.On("GenUuidV4").Return(uuidStr, nil)
	dbMock.On("SaveUser", ctx, mock.Anything).Return(nil)

	// ACT
	userRes, errRes := service.CreateUser(ctx, email, password)

	// ASSERT
	assert.Equal(t, uuidStr, userRes.Id)
	assert.Equal(t, email, userRes.Email)
	assert.Equal(t, &passwordHash, userRes.Password)
	assert.WithinDuration(t, time.Now(), userRes.CreatedAt, time.Second)
	assert.Nil(t, userRes.UpdatedAt)
	assert.Nil(t, userRes.FirstName)
	assert.Nil(t, userRes.LastName)
	assert.Equal(t, nil, errRes)
}
