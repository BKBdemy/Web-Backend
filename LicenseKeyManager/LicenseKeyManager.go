package LicenseKeyManager

import (
	"EntitlementServer/DatabaseAbstraction"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type MediaLicenseService interface {
	GetEncryptionKey(KeyID string) (EncryptionKey, error)
	GetAllLicenseKeys() []LicenseKey
	GetLicenseKey(KeyID string) (LicenseKey, error)
	AddLicenseKey(KeyID string, EncryptionKey string, ProductID int) (int, error)
	GetProductLicenseKeys(ProductID int) ([]LicenseKey, error)
}

func (lsc LicenseService) GetProductLicenseKeys(ProductID int) ([]LicenseKey, error) {
	// Return the license keys for a specific product
	keys, err := lsc.DB.GetProductLicenseKeys(ProductID)
	if err != nil {
		return nil, err
	}

	var licenseKeys []LicenseKey
	for _, key := range keys {
		licenseKeys = append(licenseKeys, LicenseKey{
			KeyID:         key.KeyID,
			EncryptionKey: EncryptionKey{Hex: key.EncryptionKey},
			ProductID:     key.ProductID,
		})
	}

	return licenseKeys, nil
}

func (lsc LicenseService) AddLicenseKey(KeyID string, NewEncryptionKey string, ProductID int) (int, error) {
	if KeyID == "" || NewEncryptionKey == "" || ProductID < 0 {
		return 0, fmt.Errorf("no arguments may be empty")
	}
	// Add the license key to the database
	key, err := lsc.DB.AddLicenseKey(DatabaseAbstraction.LicenseKey{
		KeyID:         KeyID,
		EncryptionKey: NewEncryptionKey,
		ProductID:     ProductID,
	})
	if err != nil {
		return 0, err
	}

	return key, nil
}

type LicenseService struct {
	DB DatabaseAbstraction.DBOrm
}

func (lsc LicenseService) GetEncryptionKey(keyid string) (EncryptionKey, error) {
	// Get the license key from the database
	key, err := lsc.DB.GetLicenseKeyByKeyID(keyid)
	if err != nil {
		return EncryptionKey{}, err
	}

	return EncryptionKey{
		Hex: key.EncryptionKey,
	}, nil
}

func (lsc LicenseService) GetAllLicenseKeys() []LicenseKey {
	// Get all the license keys from the database
	keys, err := lsc.DB.GetAllLicenseKeys()
	if err != nil {
		return nil
	}

	var licenseKeys []LicenseKey
	for _, key := range keys {
		licenseKeys = append(licenseKeys, LicenseKey{
			KeyID:         key.KeyID,
			EncryptionKey: EncryptionKey{Hex: key.EncryptionKey},
			ProductID:     key.ProductID,
		})
	}

	return licenseKeys
}

func (lsc LicenseService) GetLicenseKey(keyid string) (LicenseKey, error) {
	// Get the license key from the database
	key, err := lsc.DB.GetLicenseKeyByKeyID(keyid)
	if err != nil {
		return LicenseKey{}, err
	}

	return LicenseKey{
		KeyID:         key.KeyID,
		EncryptionKey: EncryptionKey{Hex: key.EncryptionKey},
		ProductID:     key.ProductID,
	}, nil
}

func corsMiddleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "http://localhost:9000")
	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Header("Access-Control-Expose-Headers", "Content-Length")
	c.Header("Access-Control-Allow-Credentials", "true")

	// Stop here if its Preflighted OPTIONS request
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	// Handle request
	c.Next()
}

func (lsc LicenseService) RegisterHandlers(r *gin.Engine, middleware ...gin.HandlerFunc) {

	// Create the HTTP handlers for the service endpoints
	r.OPTIONS("/api/licensing/keyserver/keyrequest", corsMiddleware)
	r.POST("/api/licensing/keyserver/keyrequest", corsMiddleware, middleware[0], lsc.DRMLicenseKeyRequestHandler)
}

func (lsc LicenseService) GetLabel() string {
	return "License Key Manager"
}

type shakaLicenseRequest struct {
	KeyIDs []string `json:"kids"`
	Type   string   `json:"type"`
}

type shakaLicenseResponse struct {
	Keys []shakaLicenseKeypair `json:"keys"`
	Type string                `json:"type"`
}

type shakaLicenseKeypair struct {
	KeyType string `json:"kty"`
	Key     string `json:"k"`
	KeyID   string `json:"kid"`
}

type licenseRequestError struct {
	Error string `json:"error"`
}

// DRMLicenseKeyRequestHandler godoc
//
//			@Summary		Get license key in key request format
//			@Description	ShakaPlayer license key request handler
//			@Tags			LicenseKeyManager
//			@Accept			json
//			@Produce		json
//			@Param			licenseRequest		body	shakaLicenseRequest	true	"License request"
//			@Success		200		{object}	shakaLicenseResponse
//			@Failure		400		{object}	licenseRequestError
//			@Failure		500		{object}	licenseRequestError
//	     @Security		ApiKeyAuth
//			@Router			/api/licensing/keyserver/keyrequest [post]
func (lsc LicenseService) DRMLicenseKeyRequestHandler(c *gin.Context) {
	var licenseRequest shakaLicenseRequest
	err := c.BindJSON(&licenseRequest)
	if err != nil {
		c.JSON(400, licenseRequestError{Error: "invalid request"})
		return
	}

	// Get the license keys the user has access to
	user := c.MustGet("user").(DatabaseAbstraction.User)
	licenseKeys, err := lsc.DB.GetOwnedLicenseKeys(user.IndexID)
	if err != nil {
		c.JSON(500, licenseRequestError{Error: "failed to get license keys"})
		return
	}

	if licenseRequest.Type != "temporary" {
		c.JSON(400, licenseRequestError{Error: "invalid license request type"})
		return
	}

	var licenseResponse shakaLicenseResponse
	licenseResponse.Type = "temporary"

	for _, keyID := range licenseRequest.KeyIDs {
		keyIDdecoded, err := Base64ToHex(keyID)
		if err != nil {
			c.JSON(400, licenseRequestError{Error: "couldn't decode key ID: " + err.Error()})
			return
		}

		licenseKey, err := lsc.GetEncryptionKey(keyIDdecoded)
		if err != nil {
			c.JSON(400, licenseRequestError{Error: "failed to get license key " + keyID})
			return
		}

		// Check if the user has access to the license key
		// FIXME: user keys not correctly found
		var hasAccess bool
		for _, licenseKey := range licenseKeys {
			logrus.Print(licenseKey.KeyID + " " + keyIDdecoded)
			if licenseKey.KeyID == keyIDdecoded {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			c.JSON(400, licenseRequestError{Error: "user does not have access to license key"})
			return
		}

		licenseKeyBase64, err := licenseKey.GetBase64()
		if err != nil {
			c.JSON(500, licenseRequestError{Error: "failed to get license key"})
			return
		}

		licenseResponse.Keys = append(licenseResponse.Keys, shakaLicenseKeypair{
			KeyType: "oct",
			Key:     licenseKeyBase64,
			KeyID:   keyID,
		})
	}

	c.JSON(200, licenseResponse)
}
