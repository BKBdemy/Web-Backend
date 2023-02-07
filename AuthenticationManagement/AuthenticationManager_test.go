package AuthenticationManagement_test

import (
	"EntitlementServer/AuthenticationManagement"
	"EntitlementServer/DatabaseAbstraction"
	"EntitlementServer/mocks"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestAuthenticationService_AuthenticationMiddlewareValidToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header), // if you need to test headers
	}
	// example: req.Header.Add("Accept", "application/json")

	req.Header.Add("Authorization", "Bearer 1234567890")

	// finally set the request to the gin context
	c.Request = req

	// Get fake DB
	mockDB := &mocks.DBOrm{}

	// Get fake token
	mockDB.On("GetTokenByHash", "1234567890").Return(DatabaseAbstraction.Token{
		IndexID: 1,
		UserID:  1,
		Token:   "1234567890",
		Expiry:  time.Now().Add(time.Hour * 24 * 7),
	}, nil)

	// Get fake user
	mockDB.On("GetUserByIndexID", 1).Return(DatabaseAbstraction.User{
		IndexID:   1,
		Username:  "testuser",
		Password:  "testpassword",
		Balance:   5325325,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}, nil)

	// Run the middleware
	authSvc := AuthenticationManagement.AuthenticationService{DB: mockDB}

	authSvc.AuthenticationMiddleware(c)

	// Check if the user was set
	if c.Keys["user"] == nil {
		t.Error("User was not set")
	}

	// Check if the user is the correct user
	user := c.Keys["user"].(DatabaseAbstraction.User)
	if user.Username != "testuser" {
		t.Error("User was not set correctly")
	}

	// Check if the status code is 200
	if w.Code != 200 {
		t.Error("Status code was not 200")
	}
}

func TestAuthenticationService_AuthenticationMiddlewareInvalidToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header), // if you need to test headers
	}
	// example: req.Header.Add("Accept", "application/json")

	req.Header.Add("Authorization", "Bearer 1234567890")

	// finally set the request to the gin context
	c.Request = req

	// Get fake DB
	mockDB := &mocks.DBOrm{}

	// Get no token
	mockDB.On("GetTokenByHash", "1234567890").Return(DatabaseAbstraction.Token{}, errors.New("error"))

	// Get fake user
	mockDB.On("GetUserByIndexID", 1).Return(DatabaseAbstraction.User{
		IndexID:   1,
		Username:  "testuser",
		Password:  "testpassword",
		Balance:   5325325,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}, nil)

	// Run the middleware
	authSvc := AuthenticationManagement.AuthenticationService{DB: mockDB}
	authSvc.AuthenticationMiddleware(c)

	// Check if the user was set
	if c.Keys["user"] != nil {
		t.Error("User was set")
	}

	// Check if the status code is 401
	if w.Code != 401 {
		t.Error("Status code was not 401")
	}

}
