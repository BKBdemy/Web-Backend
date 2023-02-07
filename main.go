package main

import (
	"EntitlementServer/AuthenticationManagement"
	"EntitlementServer/DatabaseAbstraction"
	"EntitlementServer/LicenseKeyManager"
	"EntitlementServer/ProductService"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
)

type HTTPService interface {
	RegisterHandlers(r *gin.Engine, middleware ...gin.HandlerFunc)
	GetLabel() string
}

// @title User API documentation
// @version 1.0.0
// @host localhost:8000
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

	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
