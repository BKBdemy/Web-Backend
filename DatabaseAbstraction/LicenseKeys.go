package DatabaseAbstraction

import (
	"context"
	"log"
	"time"
)

type LicenseKey struct {
	IndexID       int
	KeyID         string
	EncryptionKey string
	ProductID     int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (dbc *DBConnector) GetAllLicenseKeys() ([]LicenseKey, error) {
	err := dbc.DB.Ping(context.Background())

	// Get all the license keys from the database
	rows, err := dbc.DB.Query(context.Background(), "SELECT * FROM licenses")
	if err != nil {
		log.Fatal(err)
		return []LicenseKey{}, err
	}

	// Iterate over the rows and add them to the slice
	var licenseKeys []LicenseKey

	for rows.Next() {
		var licenseKey LicenseKey
		err := rows.Scan(&licenseKey.IndexID, &licenseKey.KeyID, &licenseKey.EncryptionKey, &licenseKey.ProductID, &licenseKey.CreatedAt, &licenseKey.UpdatedAt)
		if err != nil {
			return []LicenseKey{}, err
		}
		licenseKeys = append(licenseKeys, licenseKey)
	}

	return licenseKeys, nil
}

func (dbc DBConnector) GetLicenseKeyByKeyID(keyID string) (LicenseKey, error) {
	// Get the license key from the database
	row := dbc.DB.QueryRow(context.Background(), "SELECT * FROM licenses WHERE key_id = $1", keyID)

	var licenseKey LicenseKey
	err := row.Scan(&licenseKey.IndexID, &licenseKey.KeyID, &licenseKey.EncryptionKey, &licenseKey.ProductID, &licenseKey.CreatedAt, &licenseKey.UpdatedAt)
	if err != nil {
		return LicenseKey{}, err
	}

	return licenseKey, nil
}

func (dbc DBConnector) GetLicenseKeyByIndexID(indexID int) (LicenseKey, error) {
	// Get the license key from the database
	row := dbc.DB.QueryRow(context.Background(), "SELECT * FROM licenses WHERE id = $1", indexID)

	var licenseKey LicenseKey
	err := row.Scan(&licenseKey.IndexID, &licenseKey.KeyID, &licenseKey.EncryptionKey, &licenseKey.ProductID, &licenseKey.CreatedAt, &licenseKey.UpdatedAt)
	if err != nil {
		return LicenseKey{}, err
	}

	return licenseKey, nil
}

// AddLicenseKey adds a license key to the database, returned is the index ID of the newly added license key
func (dbc DBConnector) AddLicenseKey(licenseKey LicenseKey) (int, error) {
	// Add the license key to the database

	result := dbc.DB.QueryRow(context.Background(), "INSERT INTO licenses (key_id, encryption_key, product_id) VALUES ($1, $2, $3)", licenseKey.KeyID, licenseKey.EncryptionKey, licenseKey.ProductID)

	var indexID int
	err := result.Scan(&indexID)
	if err != nil {
		return -1, err
	}

	return int(indexID), nil
}

func (dbc DBConnector) DeleteLicenseKey(indexID int) error {
	// Delete the license key from the database
	_, err := dbc.DB.Exec(context.Background(), "DELETE FROM licenses WHERE id = $1", indexID)
	if err != nil {
		return err
	}

	return nil
}
