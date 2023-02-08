package ProductService

import (
	"EntitlementServer/DatabaseAbstraction"
	"EntitlementServer/LicenseKeyManager"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type Product struct {
	ID          int
	Name        string
	Description string
	Price       int
	Image       string
	MPD_URL     string
	LicenseKeys []LicenseKeyManager.LicenseKey
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
	r.GET("/api/products", p.GetAllProductsHandler)
	r.GET("/api/products/:id", p.GetProductHandler)
	r.GET("/api/products/:id/licensekeys", middleware[0], p.GetProductLicenseKeysHandler)
	r.POST("/api/products/:id/purchase", middleware[0], p.PurchaseProductHandler)
	r.GET("/api/products/owned", middleware[0], p.GetOwnedProductsHandler)
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

type productResponse struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ID          int
	Name        string
	Description string
	Price       int
	Image       string
	MPD_URL     string
}

type productErrorResponse struct {
	Error string
}

// GetAllProductsHandler godoc
// @Summary Get all products
// @Description Get all products
// @Tags Products
// @Accept  json
// @Produce  json
// @Success 200 {object} []productResponse
// @Failure 500 {object} string
// @Router /api/products [get]
func (p ProductService) GetAllProductsHandler(c *gin.Context) {
	products := p.GetAllProducts()

	productResponses := make([]productResponse, len(products))
	for i, product := range products {
		productResponses[i] = productResponse{
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Image:       product.Image,
			MPD_URL:     product.MPD_URL,
		}
	}

	c.JSON(200, products)
}

// GetProductHandler godoc
// @Summary Get a product
// @Description Get a product
// @Tags Products
// @Accept  json
// @Produce  json
// @Param id path int true "Product ID"
// @Success 200 {object} productResponse
// @Failure 400 {object} productErrorResponse
// @Failure 404 {object} productErrorResponse
// @Router /api/products/{id} [get]
func (p ProductService) GetProductHandler(c *gin.Context) {
	productID := c.Param("id")
	convertedProductID, err := strconv.Atoi(productID)
	if err != nil {
		c.JSON(400, productErrorResponse{Error: "invalid product id"})
		return
	}

	product, err := p.GetProduct(convertedProductID)
	if err != nil {
		c.JSON(404, productErrorResponse{Error: "product not found"})
		return
	}

	responseProduct := productResponse{
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Image:       product.Image,
		MPD_URL:     product.MPD_URL,
	}

	c.JSON(200, responseProduct)
}

func (p ProductService) GetLabel() string {
	return "Product Service"
}

type purchaseProductResponse struct {
	error string
	msg   string
}

// PurchaseProductHandler godoc
// @Summary Purchase a product
// @Description Purchase a product and add it to the user's owned products
// @Tags Products
// @Accept  json
// @Produce  json
// @Param id path int true "Product ID"
// @Success 200 {object} purchaseProductResponse
// @Failure 400 {object} purchaseProductResponse
// @Failure 500 {object} purchaseProductResponse
// @Security ApiKeyAuth
// @Router /api/products/{id}/purchase [post]
func (p ProductService) PurchaseProductHandler(c *gin.Context) {
	// TODO: Fix error 500 on purchase
	productID := c.Param("id")
	convertedProductID, err := strconv.Atoi(productID)
	if err != nil {
		c.JSON(400, purchaseProductResponse{error: "invalid product id"})
		return
	}

	user := c.MustGet("user").(DatabaseAbstraction.User)

	err = p.PurchaseProduct(convertedProductID, user)
	if err != nil {
		c.JSON(500, purchaseProductResponse{error: "error purchasing product: " + err.Error()})
		return
	}

	c.JSON(200, purchaseProductResponse{msg: "product purchased"})
}

// GetOwnedProductsHandler godoc
// @Summary Get owned products
// @Description Get products owned by the user
// @Tags Products
// @Accept  json
// @Produce  json
// @Success 200 {object} []productResponse
// @Failure 500 {object} string
// @Security ApiKeyAuth
// @Router /api/products/owned [get]
func (p ProductService) GetOwnedProductsHandler(c *gin.Context) {
	user := c.MustGet("user").(DatabaseAbstraction.User)

	products := p.GetOwnedProducts(user)

	productResponses := make([]productResponse, len(products))
	for i, product := range products {
		productResponses[i] = productResponse{
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Image:       product.Image,
			MPD_URL:     product.MPD_URL,
		}
	}

	c.JSON(200, productResponses)
}
