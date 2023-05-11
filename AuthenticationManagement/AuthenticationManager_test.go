package AuthenticationManagement

import (
	"EntitlementServer/DatabaseAbstraction"
	"EntitlementServer/DatabaseAbstraction/mocks"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) CreateToken(userid int) (string, error) {
	args := m.Called(userid)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) ValidateToken(token string) (bool, DatabaseAbstraction.User, error) {
	args := m.Called(token)
	return args.Bool(0), args.Get(1).(DatabaseAbstraction.User), args.Error(2)
}

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
	authSvc := AuthenticationService{DB: mockDB}

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
	authSvc := AuthenticationService{DB: mockDB}
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

func TestAuthenticationMiddleware(t *testing.T) {
	mockDB := new(mocks.DBOrm)
	mockTokenService := new(MockTokenService)
	am := AuthenticationService{
		DB: mockDB,
	}

	// Test valid token in header
	t.Run("valid token in header", func(t *testing.T) {
		user := DatabaseAbstraction.User{IndexID: 1, Username: "testuser", Password: "hashed_password"}

		mockDB.On("GetTokenByHash", "valid_token").Return(DatabaseAbstraction.Token{UserID: 1}, nil)
		mockDB.On("GetUserByIndexID", 1).Return(user, nil)

		mockTokenService.On("ValidateToken", "valid_token").Return(true, user, nil)

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		router.GET("/test", am.AuthenticationMiddleware, func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	// Test valid token in cookie
	t.Run("valid token in cookie", func(t *testing.T) {
		user := DatabaseAbstraction.User{IndexID: 1, Username: "testuser", Password: "hashed_password"}

		mockTokenService.On("ValidateToken", "valid_token").Return(true, user, nil)

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		router.GET("/test", am.AuthenticationMiddleware, func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.AddCookie(&http.Cookie{Name: "authtoken", Value: "valid_token"})
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	// Test invalid token
	t.Run("invalid token", func(t *testing.T) {
		mockDB.On("GetTokenByHash", "invalid_token").Return(DatabaseAbstraction.Token{}, errors.New("no result"))
		mockTokenService.On("ValidateToken", "invalid_token").Return(false, DatabaseAbstraction.User{}, errors.New("invalid token"))

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		router.GET("/test", am.AuthenticationMiddleware, func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	// Test no token
	t.Run("no token", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		router.GET("/test", am.AuthenticationMiddleware, func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

}

func TestLogin(t *testing.T) {
	mockDB := new(mocks.DBOrm)
	mockTokenService := new(MockTokenService)
	am := AuthenticationService{
		DB: mockDB,
	}

	user := DatabaseAbstraction.User{IndexID: 1, Username: "testuser", Password: "$argon2id$v=19$m=256000,t=6,p=1$dGVzdHRlc3Q$MMMzLViNOBi+zmhnFWj4y1y6TqYfRvmUAI6BiH30mIk"}

	t.Run("successful login", func(t *testing.T) {
		mockDB.On("GetUserByUsername", "testuser").Return(user, nil)
		mockDB.On("AuthenticateUser", "testuser", "admin").Return(true, nil)
		mockTokenService.On("CreateToken", user.IndexID).Return("valid_token", nil)
		mockDB.On("GetUserByIndexID", 1).Return(user, nil)
		mockDB.On("AddToken", 1, mock.Anything, mock.Anything).Return(nil)

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		router.POST("/api/auth/login", am.Login)

		reqBody := `{"username":"testuser","password":"admin"}`
		req, _ := http.NewRequest("POST", "/api/auth/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.NotContains(t, resp.Body.String(), "token:\":\"\"")
	})

	t.Run("invalid login", func(t *testing.T) {
		mockDB.On("AuthenticateUser", "testuser", "wrongpassword").Return(false, nil)

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		router.POST("/api/auth/login", am.Login)

		reqBody := `{"username":"testuser","password":"wrongpassword"}`
		req, _ := http.NewRequest("POST", "/api/auth/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("invalid request", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		router.POST("/api/auth/login", am.Login)

		reqBody := `{"username":"testuser"}`
		req, _ := http.NewRequest("POST", "/api/auth/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

}

func TestAuthenticationService_RegisterUserHandler(t *testing.T) {
	mockDB := new(mocks.DBOrm)
	authSvc := AuthenticationService{DB: mockDB}
	router := gin.Default()

	mockDB.On("GetUserByUsername", "testuser").Return(DatabaseAbstraction.User{}, errors.New("User not found")).Once()
	mockDB.On("AddUser", "testuser", mock.AnythingOfType("string")).Return(nil)
	mockDB.On("GetUserByUsername", "testuser").Return(DatabaseAbstraction.User{
		IndexID:  1,
		Username: "testuser",
	}, nil)
	mockDB.On("GetUserByIndexID", 1).Return(DatabaseAbstraction.User{
		IndexID:  1,
		Username: "testuser",
	}, nil)
	mockDB.On("AddToken", 1, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)

	authSvc.RegisterHandlers(router)

	// Test user registration
	data := loginRequest{
		Username: "testuser",
		Password: "testpassword",
	}

	encoded, err := json.Marshal(data)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(string(encoded)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.Token)
	assert.Empty(t, response.Error)
}

func TestAuthenticationService_RegisterUserHandlerInvalidInput(t *testing.T) {
	mockDB := new(mocks.DBOrm)
	authSvc := AuthenticationService{DB: mockDB}
	router := gin.Default()

	authSvc.RegisterHandlers(router)

	// Test user registration with missing password
	data := loginRequest{
		Username: "testuser",
	}

	encoded, err := json.Marshal(data)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(string(encoded)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response loginResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Empty(t, response.Token)
	assert.NotEmpty(t, response.Error)
}
