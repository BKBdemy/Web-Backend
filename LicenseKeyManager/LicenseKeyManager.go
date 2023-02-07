package LicenseKeyManager

import (
	"EntitlementServer/DatabaseAbstraction"
	"fmt"
	"github.com/gin-gonic/gin"
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

func (lsc LicenseService) RegisterHandlers(r *gin.Engine, middleware ...gin.HandlerFunc) {

	// Create the HTTP handlers for the service endpoints
	r.POST("/api/licensing/keyserver/keyrequest", middleware[0], DRMLicenseKeyRequestHandler)
}

func (lsc LicenseService) GetLabel() string {
	return "License Key Manager"
}

type licenseRequest struct {
}

// DRMLicenseKeyRequestHandler
// @Summary Get license key in key request format
// @Description Internal videoplayer API
// @Tags Licensing
// @Success 200 {array} model.User
// @Failure 401 {object} object
// @Router / [get]
func DRMLicenseKeyRequestHandler(c *gin.Context) {
	// User DRM client is requesting a license key
	// for a specific content item

}
