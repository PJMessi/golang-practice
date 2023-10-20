package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestIntegrationLoginWithUnregisteredEmail(t *testing.T) {
	// ARRANGE
	setupIntegrationTest()
	defer teardownIntegrationTest()

	url := fmt.Sprintf("%s/auth/login", testServer.URL)
	email := testutil.Fake.Internet().Email()
	password := testutil.Fake.Internet().Password()

	// ACT
	reqBody := []byte(fmt.Sprintf(`{"email": "%s","password": "%s"}`, email, password))
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody))

	// ASSERT
	responseBody, _ := io.ReadAll(resp.Body)
	expectedResponseBody := `{"type":"UNAUTHENTICATED","message":"invalid credentials","details":null}`
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "should return 401 status code")
	assert.Equal(t, expectedResponseBody, string(responseBody), "should return error details in the response body")
}

func TestIntegrationLoginWithIncorrectPassword(t *testing.T) {
	// ARRANGE
	setupIntegrationTest()
	defer teardownIntegrationTest()

	user := testutil.GenMockUser(nil)
	url := fmt.Sprintf("%s/auth/login", testServer.URL)
	email := user.Email
	password := testutil.Fake.Internet().Password()

	// adding user with the email in the database with random password hash
	smt, _ := testDbCon.Prepare("INSERT INTO users (id, email, password, first_name, last_name, created_at, updated_at) VALUE (?, ?, ?, ?, ?, ?, ?)")
	smt.Exec(user.Id, user.Email, user.Password, user.FirstName, user.LastName, user.CreatedAt, user.UpdatedAt)
	smt.Close()

	// ACT
	reqBody := []byte(fmt.Sprintf(`{"email": "%s","password": "%s"}`, email, password))
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody))

	// ASSERT
	responseBody, _ := io.ReadAll(resp.Body)
	expectedResponseBody := `{"type":"UNAUTHENTICATED","message":"invalid credentials","details":null}`
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "should return 401 status code")
	assert.Equal(t, expectedResponseBody, string(responseBody), "should return error details in the response body")
}

func TestIntegrationLoginSuccessfulResponse(t *testing.T) {
	// ARRANGE
	setupIntegrationTest()
	defer teardownIntegrationTest()

	password := "Password123!"
	passwordHashBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	passwordHash := string(passwordHashBytes)
	user := testutil.GenMockUser(&model.User{Password: &passwordHash})
	url := fmt.Sprintf("%s/auth/login", testServer.URL)
	email := user.Email

	// adding user with the email in the database with random password hash
	smt, _ := testDbCon.Prepare("INSERT INTO users (id, email, password, first_name, last_name, created_at, updated_at) VALUE (?, ?, ?, ?, ?, ?, ?)")
	smt.Exec(user.Id, user.Email, user.Password, user.FirstName, user.LastName, user.CreatedAt, user.UpdatedAt)
	smt.Close()

	// ACT
	reqBody := []byte(fmt.Sprintf(`{"email": "%s","password": "%s"}`, email, password))
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody))

	// ASSERT
	responseBodyByte, _ := io.ReadAll(resp.Body)
	responseBody := model.LoginApiRes{}
	_ = json.Unmarshal(responseBodyByte, &responseBody)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "should return 200 status code")
	assert.Equal(t, email, responseBody.User.Email, "should return user email in the response body")
	assert.NotNil(t, responseBody.User.Id, "should return user id in the response body")
	assert.Equal(t, user.CreatedAt.Format(time.RFC3339), responseBody.User.CreatedAt, "should return user creation date in the response body")
	assert.NotNil(t, responseBody.Jwt, "should return jwt for the user in the response body")
}
