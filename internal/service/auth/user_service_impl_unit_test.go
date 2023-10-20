package auth

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/pjmessi/golang-practice/internal/errorcode"
	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/internal/pkg/database"
	"github.com/pjmessi/golang-practice/internal/pkg/jwt"
	"github.com/pjmessi/golang-practice/internal/pkg/testutil"
	"github.com/pjmessi/golang-practice/pkg/exception"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// setupMocksForServiceImplTest creates ServiceImpl with mocked dependencies
func setupMocksForServiceImplTest() (*ServiceImpl, *database.DbMock, *jwt.HandlerMock, *logger.ServiceMock) {
	dbMock := new(database.DbMock)
	jwtHandlerMock := new(jwt.HandlerMock)
	logServiceMock := new(logger.ServiceMock)
	service := &ServiceImpl{
		db:         dbMock,
		jwtHandler: jwtHandlerMock,
		logService: logServiceMock,
	}
	return service, dbMock, jwtHandlerMock, logServiceMock
}

// setupMocksForNewService returns mocked dependencies for NewService func
func setupMocksForNewService() (*database.DbMock, *jwt.HandlerMock, *logger.ServiceMock) {
	dbMock := new(database.DbMock)
	jwtHandlerMock := new(jwt.HandlerMock)
	logServiceMock := new(logger.ServiceMock)
	return dbMock, jwtHandlerMock, logServiceMock
}

func Test_NewService(t *testing.T) {
	// ARRANGE
	dbMock, jwtHandlerMock, logServiceMock := setupMocksForNewService()

	// ACT
	res := NewService(logServiceMock, jwtHandlerMock, dbMock)

	// ARRANGE
	resServiceImpl := res.(*ServiceImpl)

	assert.IsType(t, &ServiceImpl{}, res)
	assert.Equal(t, res, resServiceImpl)
}

func Test_Service_Login_User_Doesnt_exist(t *testing.T) {
	// ARRANGE
	service, dbMock, _, logServiceMock := setupMocksForServiceImplTest()

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
	expectedErr := exception.NewUnauthenticatedFromBase(exception.Base{
		Message: "invalid credentials",
	})
	expectedLogStr := fmt.Sprintf("user with the email '%s' does not exist", email)

	assert.Equal(t, userRes, expectedUserRes)
	assert.Equal(t, jwtStrRes, expectedJwtStrRes)
	assert.Equal(t, errRes, expectedErr)
	logServiceMock.AssertCalled(t, "DebugCtx", ctx, expectedLogStr)
}

func Test_Service_Login_Err_Getting_User_By_Email(t *testing.T) {
	// ARRANGE
	service, dbMock, _, logServiceMock := setupMocksForServiceImplTest()

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
	service, dbMock, _, logServiceMock := setupMocksForServiceImplTest()

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
	service, dbMock, _, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := testutil.Fake.Internet().Password()
	user := testutil.GenMockUser(&model.User{Email: email})

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("GetUserByEmail", ctx, email).Return(true, user, nil)

	// ACT
	userRes, jwtStrRes, errRes := service.Login(ctx, email, password)

	// ASSERT
	expectedUserRes := model.User{}
	expectedJwtStrRes := ""
	expectedErr := exception.NewUnauthenticatedFromBase(exception.Base{
		Message: "invalid credentials",
	})

	assert.Equal(t, userRes, expectedUserRes)
	assert.Equal(t, jwtStrRes, expectedJwtStrRes)
	assert.Equal(t, errRes, expectedErr)
}

func Test_Service_Login_Success_Res(t *testing.T) {
	// ARRANGE
	service, dbMock, jwtHandlerMock, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := strings.ToUpper(testutil.Fake.Internet().Email())
	password := "Password123!"
	hashedPasswordByte, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedPassword := string(hashedPasswordByte)
	user := testutil.GenMockUser(&model.User{Email: email, Password: &hashedPassword})
	jwtStr := testutil.Fake.RandomStringWithLength(100)

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("GetUserByEmail", ctx, email).Return(true, user, nil)
	jwtHandlerMock.On("Generate", jwt.JwtPayload{UserId: user.Id, UserEmail: user.Email}).Return(jwtStr, nil)

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
	service, dbMock, jwtHandlerMock, logServiceMock := setupMocksForServiceImplTest()

	ctx := context.Background()
	email := testutil.Fake.Internet().Email()
	password := "Password123!"
	hashedPasswordByte, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedPassword := string(hashedPasswordByte)
	user := testutil.GenMockUser(&model.User{Email: email, Password: &hashedPassword})
	generateErr := fmt.Errorf("error from Generate")

	logServiceMock.On("DebugCtx", mock.Anything, mock.Anything)
	dbMock.On("GetUserByEmail", ctx, email).Return(true, user, nil)
	jwtHandlerMock.On("Generate", jwt.JwtPayload{UserId: user.Id, UserEmail: user.Email}).Return("", generateErr)

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
