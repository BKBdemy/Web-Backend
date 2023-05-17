package ProductService

import (
	"EntitlementServer/DatabaseAbstraction"
	"EntitlementServer/VideoService"
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
	Difficulty  int
	Videos      []VideoService.VSVideo
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductServiceProvider interface {
	GetProduct(ProductID int) (Product, error)
	GetAllProducts() []Product
	PurchaseProduct(ProductID int, user DatabaseAbstraction.User) error
	GetOwnedProducts(user DatabaseAbstraction.User) []Product
	AddProduct(Product Product) (int, error)
}

type ProductService struct {
	DB DatabaseAbstraction.DBOrm
}

func (p ProductService) RegisterHandlers(r *gin.Engine, middleware ...gin.HandlerFunc) {
	r.GET("/api/products", p.GetAllProductsHandler)
	r.GET("/api/products/:id", p.GetProductHandler)
	r.POST("/api/products/:id/purchase", middleware[0], p.PurchaseProductHandler)
	r.GET("/api/products/owned", middleware[0], p.GetOwnedProductsHandler)
}

type productResponse struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ID          int
	Name        string
	Description string
	Price       int
	Image       string
	Difficulty  int
	Videos      []VideoService.VSVideo
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
			Difficulty:  product.Difficulty,
			Videos:      product.Videos,
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
		Videos:      product.Videos,
		Difficulty:  product.Difficulty,
	}

	c.JSON(200, responseProduct)
}

func (p ProductService) GetLabel() string {
	return "Product Service"
}

type purchaseProductResponse struct {
	Error   string
	Message string
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
	productID := c.Param("id")
	convertedProductID, err := strconv.Atoi(productID)
	if err != nil {
		c.JSON(400, purchaseProductResponse{Error: "invalid product id"})
		return
	}

	user := c.MustGet("user").(DatabaseAbstraction.User)

	err = p.PurchaseProduct(convertedProductID, user)
	if err != nil {
		c.JSON(400, purchaseProductResponse{Error: "Error purchasing product: " + err.Error()})
		return
	}

	c.JSON(200, purchaseProductResponse{Message: "product purchased"})
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
		// Get the videos for the product
		videos, err := p.DB.GetVideosByProductIndexID(product.ID)
		if err != nil {
			c.JSON(500, gin.H{"Error": "Error getting videos"})
			return
		}
		// Convert the videos to video responses
		videoResponses := make([]VideoService.VSVideo, len(videos))
		for i, video := range videos {
			videoResponses[i] = VideoService.VSVideo{
				IndexID:     video.IndexID,
				Name:        video.Name,
				Description: video.Description,
				Points:      video.Points,
				Thumbnail:   video.Thumbnail,
				Filename:    video.Filename,
			}
		}

		productResponses[i] = productResponse{
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Image:       product.Image,
			Videos:      videoResponses,
			Difficulty:  product.Difficulty,
		}
	}

	c.JSON(200, productResponses)
}
