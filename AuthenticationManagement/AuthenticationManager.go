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

type NotSignedInResponse struct {
	error string
}

func (am AuthenticationService) AuthenticationMiddleware(c *gin.Context) {
	// Check token in autorization header
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
	r.GET("/api/auth/me", am.AuthenticationMiddleware, am.GetUserHandler)
	/*r.POST("/api/logout", am.Logout)*/
	r.POST("/api/auth/register", am.RegisterUserHandler)
}

func (am AuthenticationService) GetLabel() string {
	return "Authentication Service"
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type loginResponse struct {
	Token string `json:"token"`
	Error string `json:"error"`
}

// Login godoc
// @Summary Login to the application and get a token
// @Description Login to the application and get a token, token is valid for 7 days
// @Description error is empty if login was successful
// @Tags Authentication
// @Accept json
// @Produce json
// @Param loginRequest body loginRequest true "Login request"
// @Success 200 {object} loginResponse
// @Failure 400 {object} loginResponse
// @Failure 401 {object} loginResponse
// @Failure 500 {object} loginResponse
// @Router /api/auth/login [post]
func (am AuthenticationService) Login(ctx *gin.Context) {

	var request loginRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(400, loginResponse{
			Token: "",
			Error: "Invalid request",
		})
		return
	}

	if request.Username == "" || request.Password == "" {
		ctx.JSON(400, loginResponse{
			Token: "",
			Error: "Empty username or password",
		})
		return
	}

	valid, err := am.AuthenticateUser(request.Username, request.Password)
	if err != nil {
		ctx.JSON(401, loginResponse{
			Token: "",
			Error: "Invalid username or password",
		})
		return
	}

	if !valid {
		ctx.JSON(401, loginResponse{
			Token: "",
			Error: "Invalid username or password",
		})
		return
	}

	// Now we need to create a token for the user

	// Get the user from the database
	user, err := am.DB.GetUserByUsername(request.Username)

	userToken, err := am.CreateToken(user.IndexID)
	if err != nil {
		ctx.JSON(500, loginResponse{
			Token: "",
			Error: "Failed to create token",
		})
		return
	}

	ctx.JSON(200, loginResponse{
		Token: userToken,
		Error: "",
	})
}

type meResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Balance   int    `json:"balance"`
	CreatedAt string `json:"created_at"`
}

// GetUserHandler godoc
// @Summary Get the current user
// @Description Get the current user from the token
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} meResponse
// @Failure 401 {object} NotSignedInResponse
// @Security ApiKeyAuth
// @Router /api/auth/me [get]
func (am AuthenticationService) GetUserHandler(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(200, gin.H{"id": user.(DatabaseAbstraction.User).IndexID, "username": user.(DatabaseAbstraction.User).Username})
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterUserHandler godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param registerRequest body registerRequest true "Register request"
// @Success 200 {object} loginResponse
// @Failure 400 {object} loginResponse
// @Failure 500 {object} loginResponse
// @Router /api/auth/register [post]
func (am AuthenticationService) RegisterUserHandler(c *gin.Context) {
	registerRequest := registerRequest{}
	err := c.ShouldBindJSON(&registerRequest)
	if err != nil {
		c.JSON(400, loginResponse{
			Token: "",
			Error: "Invalid request",
		})
		return
	}

	if registerRequest.Username == "" || registerRequest.Password == "" {
		c.JSON(400, loginResponse{
			Token: "",
			Error: "Empty username or password",
		})
		return
	}

	// Check if the user already exists
	user, err := am.DB.GetUserByUsername(registerRequest.Username)
	if err == nil {
		c.JSON(400, loginResponse{
			Token: "",
			Error: "User already exists",
		})
		return
	}

	// Create the user
	err = am.CreateUser(registerRequest.Username, registerRequest.Password)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, loginResponse{
			Token: "",
			Error: "Failed to create user",
		})
		return
	}

	// Get the user from the database
	user, err = am.DB.GetUserByUsername(registerRequest.Username)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, loginResponse{
			Token: "",
			Error: "Failed to get user",
		})
		return
	}

	// Generate a token for the user
	token, err := am.CreateToken(user.IndexID)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, loginResponse{
			Token: "",
			Error: "Failed to create token",
		})
		return
	}

	c.JSON(200, loginResponse{
		Token: token,
		Error: "",
	})
}
