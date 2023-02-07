package DatabaseAbstraction

import (
	"context"
	"time"
)

type Product struct {
	IndexID     int
	Name        string
	Description string
	Price       int
	Image       string
	MPDURL      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (dbc DBConnector) GetAllProducts() ([]Product, error) {
	// Get all the products from the database
	rows, err := dbc.DB.Query(context.Background(), "SELECT * FROM products")
	if err != nil {
		return []Product{}, err
	}

	// Iterate over the rows and add them to the slice
	var products []Product

	for rows.Next() {
		var product Product
		err := rows.Scan(&product.IndexID, &product.Name, &product.Description, &product.Price, &product.Image, &product.MPDURL, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return []Product{}, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (dbc DBConnector) GetProductByIndexID(indexID int) (Product, error) {
	// Get the product from the database
	row := dbc.DB.QueryRow(context.Background(), "SELECT * FROM products WHERE id = $1", indexID)

	var product Product
	err := row.Scan(&product.IndexID, &product.Name, &product.Description, &product.Price, &product.Image, &product.MPDURL, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

// GetProductLicenseKeys returns all the license keys associated with a product
func (dbc DBConnector) GetProductLicenseKeys(indexID int) ([]LicenseKey, error) {
	// Get all the license keys from the database related to the product
	rows, err := dbc.DB.Query(context.Background(), "SELECT key_id, encryption_key FROM product_licenses JOIN licenses l on product_licenses.license_id = l.id WHERE product_licenses.product_id = $1", indexID)
	if err != nil {
		return []LicenseKey{}, err
	}

	type licenseKeyQuery struct {
		KeyID         string
		EncryptionKey string
	}

	// Iterate over the rows and add them to the slice
	var licenseKeys []licenseKeyQuery
	for rows.Next() {
		var licenseKey licenseKeyQuery
		err := rows.Scan(&licenseKey.KeyID, &licenseKey.EncryptionKey)
		if err != nil {
			return []LicenseKey{}, err
		}
		licenseKeys = append(licenseKeys, licenseKey)
	}

	// Convert the license keys to the correct format
	var completeLicenseKeys []LicenseKey
	for _, licenseKey := range licenseKeys {
		fullLicense, _ := dbc.GetLicenseKeyByKeyID(licenseKey.KeyID)
		completeLicenseKeys = append(completeLicenseKeys, fullLicense)
	}

	return completeLicenseKeys, nil
}

func (dbc DBConnector) AddProduct(NewProduct Product) (int, error) {
	// Insert the product into the database
	var indexID int
	err := dbc.DB.QueryRow(context.Background(), "INSERT INTO products (name, description, price, image, mpd_url) VALUES ($1, $2, $3, $4, $5) RETURNING id", NewProduct.Name, NewProduct.Description, NewProduct.Price, NewProduct.Image, NewProduct.MPDURL).Scan(&indexID)
	if err != nil {
		return -1, err
	}

	return indexID, nil
}
