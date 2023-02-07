package ProductService

import (
	"EntitlementServer/DatabaseAbstraction"
	"EntitlementServer/LicenseKeyManager"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Product struct {
	ID          int
	Name        string
	Description string
	Price       int
	Image       string
	MPD_URL     string
	LicenseKeys []LicenseKeyManager.LicenseKey
	CreatedAt   string
	UpdatedAt   string
}

type ProductServiceProvider interface {
	GetProduct(ProductID int) (Product, error)
	GetAllProducts() []Product
	PurchaseProduct(ProductID int, user DatabaseAbstraction.User) error
	GetOwnedProducts(user DatabaseAbstraction.User) []Product
	GetProductLicenseKeys(ProductID int) []LicenseKeyManager.LicenseKey
	AddProduct(Product Product) (int, error)
}

type ProductService struct {
	DB DatabaseAbstraction.DBOrm
}

func (p ProductService) RegisterHandlers(r *gin.Engine, middleware ...gin.HandlerFunc) {

	r.GET("/api/products", middleware[0], p.GetAllProductsHandler)
	//r.GET("/api/products/:id", p.GetProductHandler)
	r.GET("/api/products/:id/licensekeys", middleware[0], p.GetProductLicenseKeysHandler)
	//r.POST("/api/products/:id/purchase", p.PurchaseProductHandler)

}

func (p ProductService) GetProductLicenseKeysHandler(c *gin.Context) {
	productID := c.Param("id")
	convertedProductID, err := strconv.Atoi(productID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid product id"})
		return
	}

	licenseKeys := p.GetProductLicenseKeys(convertedProductID)
	if licenseKeys == nil {
		c.JSON(404, gin.H{"error": "product not found"})
		return
	}

	userOwnedProducts := p.GetOwnedProducts(c.MustGet("user").(DatabaseAbstraction.User))

	// Check if the user owns the product
	var productOwned bool
	for _, product := range userOwnedProducts {
		if product.ID == convertedProductID {
			productOwned = true
			break
		}
	}

	if !productOwned {
		c.JSON(403, gin.H{"error": "product not owned"})
	} else {
		c.JSON(200, licenseKeys)
	}
}

func (p ProductService) GetAllProductsHandler(c *gin.Context) {
	products := p.GetAllProducts()

	// Strip the license keys from the products
	for i := range products {
		products[i].LicenseKeys = nil
	}

	c.JSON(200, products)
}

func (p ProductService) GetLabel() string {
	return "Product Service"
}
