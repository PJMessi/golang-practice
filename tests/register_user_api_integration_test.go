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
)

func TestIntegrationRegisterUserWithWeakPassword(t *testing.T) {
	// ARRANGE
	setupIntegrationTest()
	defer teardownIntegrationTest()

	url := fmt.Sprintf("%s/users/registration", testServer.URL)
	email := testutil.Fake.Internet().Email()
	password := "password"

	// ACT
	reqBody := []byte(fmt.Sprintf(`{"email": "%s","password": "%s"}`, email, password))
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody))

	// ASSERT
	responseBody, _ := io.ReadAll(resp.Body)
	expectedResponseBody := `{"type":"REQUEST_DATA.INVALID","message":"invalid request data","details":{"password":"password not strong"}}`
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode, "should return 422 status code")
	assert.Equal(t, expectedResponseBody, string(responseBody), "should return error details in the response body")
}

func TestIntegrationRegisterUserWithInvalidRequestBody(t *testing.T) {
	// ARRANGE
	setupIntegrationTest()
	defer teardownIntegrationTest()

	url := fmt.Sprintf("%s/users/registration", testServer.URL)
	email := ""
	password := ""

	// ACT
	reqBody := []byte(fmt.Sprintf(`{"email": "%s","password": "%s"}`, email, password))
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody))

	// ASSERT
	responseBody, _ := io.ReadAll(resp.Body)
	expectedResponseBody := `{"type":"REQUEST_DATA.INVALID","message":"invalid request data","details":{"email":"validation failed for tag: 'required'","password":"validation failed for tag: 'required'"}}`
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode, "should return 422 status code")
	assert.Equal(t, expectedResponseBody, string(responseBody), "should return the error details in the response body")
}

func TestIntegrationRegisterUserWithUsedEmail(t *testing.T) {
	// ARRANGE
	setupIntegrationTest()
	defer teardownIntegrationTest()

	user := testutil.GenMockUser(nil)
	url := fmt.Sprintf("%s/users/registration", testServer.URL)
	email := user.Email
	password := "Password123!"

	// adding user with the email in the database
	smt, _ := testDbCon.Prepare("INSERT INTO users (id, email, password, first_name, last_name, created_at, updated_at) VALUE (?, ?, ?, ?, ?, ?, ?)")
	smt.Exec(user.Id, user.Email, user.Password, user.FirstName, user.LastName, user.CreatedAt, user.UpdatedAt)
	smt.Close()

	// ACT
	reqBody := []byte(fmt.Sprintf(`{"email": "%s","password": "%s"}`, email, password))
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody))

	// ASSERT
	responseBody, _ := io.ReadAll(resp.Body)
	expectedResponseBody := fmt.Sprintf(`{"type":"USER.ALREADY_EXISTS","message":"user with the email '%s' already exists","details":null}`, email)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "should return 400 status code")
	assert.Equal(t, expectedResponseBody, string(responseBody), "should return the error details in the response body")
}

func TestIntegrationRegisterUserSuccessfulRegistration(t *testing.T) {
	// ARRANGE
	setupIntegrationTest()
	defer teardownIntegrationTest()

	user := testutil.GenMockUser(nil)
	url := fmt.Sprintf("%s/users/registration", testServer.URL)
	email := user.Email
	password := "Password123!"

	// ACT
	reqBody := []byte(fmt.Sprintf(`{"email": "%s","password": "%s"}`, email, password))
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody))

	// ASSERT
	responseBodyByte, _ := io.ReadAll(resp.Body)
	responseBody := model.UserRegApiRes{}
	_ = json.Unmarshal(responseBodyByte, &responseBody)
	userCreatedAtRes, _ := time.Parse(time.RFC3339, responseBody.User.CreatedAt)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "should return 200 status code")
	assert.Equal(t, email, responseBody.User.Email, "should return user email in the response body")
	assert.NotNil(t, responseBody.User.Id, "should return user id in the response body")
	assert.WithinDuration(t, time.Now(), userCreatedAtRes, time.Second, "should return user creation date in the response body")

	count := 0
	res, _ := testDbCon.Query(fmt.Sprintf("SELECT COUNT(*) FROM users WHERE email=\"%s\"", email))
	if res.Next() {
		res.Scan(&count)
	}
	assert.Equal(t, 1, count, "there should be a user with the email in the database")
}
