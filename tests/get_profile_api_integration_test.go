package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/pjmessi/golang-practice/internal/model"
	"github.com/pjmessi/golang-practice/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationGetProfileWithoutToken(t *testing.T) {
	// ARRANGE
	setupIntegrationTest()
	defer teardownIntegrationTest()

	url := fmt.Sprintf("%s/users/profile", testServer.URL)

	// ACT
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := http.DefaultClient.Do(req)

	// ASSERT
	responseBody, _ := io.ReadAll(resp.Body)
	expectedResponseBody := `{"type":"UNAUTHENTICATED","message":"user not authenticated","details":null}`
	assert.Equal(t, expectedResponseBody, string(responseBody), "should return error details in the response body")
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "should return 401 status code")
}

func TestIntegrationGetProfileWithInvalidToken(t *testing.T) {
	// ARRANGE
	setupIntegrationTest()
	defer teardownIntegrationTest()

	url := fmt.Sprintf("%s/users/profile", testServer.URL)

	// ACT
	req, _ := http.NewRequest("GET", url, nil)
	headers := http.Header{}
	headers.Add("Authorization", fmt.Sprintf("Bearer %s", "invalidjwt"))
	req.Header = headers
	resp, _ := http.DefaultClient.Do(req)

	// ASSERT
	responseBody, _ := io.ReadAll(resp.Body)
	expectedResponseBody := `{"type":"UNAUTHENTICATED","message":"user not authenticated","details":null}`
	assert.Equal(t, expectedResponseBody, string(responseBody), "should return error details in the response body")
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "should return 401 status code")
}

func TestIntegrationGetProfileSuccessfulResponse(t *testing.T) {
	// ARRANGE
	setupIntegrationTest()
	defer teardownIntegrationTest()

	loginRes := testutil.SetupTestUser(testServer.URL)

	url := fmt.Sprintf("%s/users/profile", testServer.URL)

	// ACT
	req, _ := http.NewRequest("GET", url, nil)
	headers := http.Header{}
	headers.Add("Authorization", fmt.Sprintf("Bearer %s", loginRes.Jwt))
	req.Header = headers
	resp, _ := http.DefaultClient.Do(req)

	// ASSERT
	responseBodyByte, _ := io.ReadAll(resp.Body)
	responseBody := model.GetProfileApiRes{}
	_ = json.Unmarshal(responseBodyByte, &responseBody)
	assert.Equal(t, loginRes.User.Email, responseBody.User.Email, "should return user email in the response body")
	assert.Equal(t, loginRes.User.Id, responseBody.User.Id, "should return user id in the response body")
	assert.Equal(t, loginRes.User.CreatedAt, responseBody.User.CreatedAt, "should return user creation time in the response body")
	assert.Equal(t, loginRes.User.FirstName, responseBody.User.FirstName, "should return user first name in the response body")
	assert.Equal(t, loginRes.User.LastName, responseBody.User.LastName, "should return user last name in the response body")
}
