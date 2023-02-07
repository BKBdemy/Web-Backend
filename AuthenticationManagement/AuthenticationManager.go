package AuthenticationManagement

import (
	"EntitlementServer/DatabaseAbstraction"
	"github.com/gin-gonic/gin"
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
	/*r.POST("/api/logout", am.Logout)
	r.POST("/api/register", am.Register)*/
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
