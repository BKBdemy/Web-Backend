package AuthenticationManagement

import (
	"EntitlementServer/DatabaseAbstraction"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
)

type AuthenticationManager interface {
	AuthenticateUser(username string, password string) (bool, error)
	CreateToken(userid int) (string, error)
	ValidateToken(token string) (bool, DatabaseAbstraction.User, error)
	CreateUser(username string, password string) error
	HashPassword(password string) (string, error)
	ComparePasswords(hashedPassword string, password string) (bool, error)
	AuthenticationMiddleware(c *gin.Context)
}

type AuthenticationService struct {
	DB DatabaseAbstraction.DBOrm
}

func (am AuthenticationService) AuthenticationMiddleware(c *gin.Context) {
	// Check token in authorization header
	// Populate user in context

	// Get the token from the header
	token := c.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	// Validate the token
	valid, user, err := am.ValidateToken(token)

	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	if !valid {
		c.JSON(401, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	// Set the user in the context
	c.Set("user", user)
	c.Next()
}

func (am AuthenticationService) RegisterHandlers(r *gin.Engine, _ ...gin.HandlerFunc) {
	r.POST("/api/auth/login", am.Login)
	r.GET("/api/auth/me", am.AuthenticationMiddleware, func(c *gin.Context) {
		user, _ := c.Get("user")
		c.JSON(200, gin.H{"id": user.(DatabaseAbstraction.User).IndexID, "username": user.(DatabaseAbstraction.User).Username})
	})
	r.POST("/api/logout", am.Logout)
	r.POST("/api/register", am.Register)
}

func (am AuthenticationService) GetLabel() string {
	return "Authentication Service"
}

func (am AuthenticationService) Login(ctx *gin.Context) {
	type loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var request loginRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Bad request"})
		return
	}

	if request.Username == "" || request.Password == "" {
		ctx.JSON(400, gin.H{"error": "Empty username or password"})
		return
	}

	valid, err := am.AuthenticateUser(request.Username, request.Password)
	if err != nil {
		ctx.JSON(401, gin.H{"error": "User not found"})
		return
	}

	if !valid {
		ctx.JSON(401, gin.H{"error": "Invalid username or password"})
		return
	}

	// Now we need to create a token for the user

	// Get the user from the database
	user, err := am.DB.GetUserByUsername(request.Username)

	userToken, err := am.CreateToken(user.IndexID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"token": userToken})
}

// @description Invalidate the current token
// @tags Authentication
// @accept json
// @produce json
// @success 200 {object} string "OK"
func (am AuthenticationService) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		ctx.JSON(400, gin.H{"error": "No token provided"})
		// The only reason this would happen is if the user is not logged in,
		// which should be literally impossible because of the authentication middleware
		logrus.Error("No token available for logout")
		return
	}

	// Invalidate the token
	err := am.DB.DeleteTokenByHash(token)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "OK"})
}

// @description Register an account
// @tags Authentication
// @accept json
// @produce json
// @param username body string true "Username"
// @param password body string true "Password"
// @success 200 {object} string "token"
func (am AuthenticationService) Register(ctx *gin.Context) {
	type registerRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var request registerRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Bad request"})
		return
	}

	if request.Username == "" || request.Password == "" {
		ctx.JSON(400, gin.H{"error": "Empty username or password"})
		return
	}

	// Check if the user already exists
	user, err := am.DB.GetUserByUsername(request.Username)
	if err == nil {
		ctx.JSON(400, gin.H{"error": "User already exists"})
		return
	}

	// Create the user
	err = am.CreateUser(request.Username, request.Password)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Get the user from the database
	user, err = am.DB.GetUserByUsername(request.Username)

	userToken, err := am.CreateToken(user.IndexID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"token": userToken})
}
