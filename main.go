package main

import (
	"EntitlementServer/AuthenticationManagement"
	"EntitlementServer/DatabaseAbstraction"
	"EntitlementServer/LicenseKeyManager"
	"EntitlementServer/ProductService"
	_ "EntitlementServer/docs"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

type HTTPService interface {
	RegisterHandlers(r *gin.Engine, middleware ...gin.HandlerFunc)
	GetLabel() string
}

// @title BKBdemy API documentation
// @version 1.0.0
// @host localhost:8080
// @BasePath /
func main() {
	DB := DatabaseAbstraction.DBConnector{}
	conn, err := DatabaseAbstraction.Connect()
	if err != nil {
		logrus.Fatal(err)
	}
	DB.DB = conn
	defer conn.Close()

	licenseSvc := LicenseKeyManager.LicenseService{DB: &DB}
	authenticationSvc := AuthenticationManagement.AuthenticationService{DB: &DB}
	productSvc := ProductService.ProductService{DB: &DB}

	r := gin.Default()

	// Register the HTTP handlers for the services
	// The authentication service is always first, it may be ignored if authentication is not needed by the service endpoint
	// Even if you don't need authentication, you still need to register the service BEFORE the other services
	// Middleware registration must happen in every route, because all middleware ties into a central router and a .Use call will apply to all routes
	licenseSvc.RegisterHandlers(r, authenticationSvc.AuthenticationMiddleware)
	authenticationSvc.RegisterHandlers(r)
	productSvc.RegisterHandlers(r, authenticationSvc.AuthenticationMiddleware)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
