package ProductService

import (
	"EntitlementServer/DatabaseAbstraction"
	"EntitlementServer/VideoService"
	"errors"
	"github.com/sirupsen/logrus"
)

func (p ProductService) GetProduct(ProductID int) (Product, error) {
	product, err := p.DB.GetProductByIndexID(ProductID)
	if err != nil {
		return Product{}, err
	}
	if product.IndexID == 0 {
		return Product{}, errors.New("product not found")
	}

	return p.enrichDatabaseProducts([]DatabaseAbstraction.Product{product})[0], nil
}

func (p ProductService) GetAllProducts() []Product {
	// Get all the products from the database
	products, err := p.DB.GetAllProducts()
	if err != nil {
		return []Product{}
	}

	return p.enrichDatabaseProducts(products)
}

var ErrNotEnoughMoney = errors.New("Not enough money")

func (p ProductService) PurchaseProduct(ProductID int, user DatabaseAbstraction.User) error {
	product, err := p.DB.GetProductByIndexID(ProductID)
	if err != nil {
		return err
	}

	// Get a fresh user object from the database
	user, err = p.DB.GetUserByIndexID(user.IndexID)
	if err != nil {
		return err
	}

	// Check if the user already owns the product
	ownedProducts, err := p.DB.GetOwnedProducts(user.IndexID)
	if err != nil {
		return err
	}
	for _, ownedProduct := range ownedProducts {
		if ownedProduct.IndexID == product.IndexID {
			return errors.New("user already owns product")
		}
	}

	// Check if the user has enough money to purchase the product
	if user.Balance < product.Price {
		return ErrNotEnoughMoney
	}

	// Update the user's balance
	err = p.DB.DecreaseUserBalance(user.IndexID, product.Price)
	if err != nil {
		logrus.Error(err)
		// If this Error happens, we have probably just prevented a racy purchase
		return err
	}

	// Add the product to the user's owned products
	err = p.DB.AddOwnedProduct(user.IndexID, product.IndexID)
	if err != nil {
		logrus.Error(err)
		return err
	}

	logrus.Println("User", user.IndexID, "purchased product", product.IndexID)

	return nil
}

func (p ProductService) GetOwnedProducts(user DatabaseAbstraction.User) []Product {
	// Get owned products from the database
	ownedProducts, err := p.DB.GetOwnedProducts(user.IndexID)
	if err != nil {
		logrus.Error(err)
		return []Product{}
	}

	// Convert the products to the correct format
	return p.enrichDatabaseProducts(ownedProducts)
}

// enrichDatabaseProducts takes a slice of database products and converts them to the ProductManagement format, including the license keys
func (p ProductService) enrichDatabaseProducts(products []DatabaseAbstraction.Product) []Product {
	// Convert the products to the correct format
	var convertedProducts []Product
	for _, product := range products {
		// Fetch the videos for the product
		videos, err := p.DB.GetVideosByProductIndexID(product.IndexID)
		if err != nil {
			continue
		}
		vsvideos := make([]VideoService.VSVideo, len(videos))
		for i, video := range videos {
			vsvideos[i] = VideoService.VSVideo{
				IndexID:     video.IndexID,
				Name:        video.Name,
				Description: video.Description,
				Points:      video.Points,
				Thumbnail:   video.Thumbnail,
			}
		}

		convertedProducts = append(convertedProducts, Product{
			ID:          product.IndexID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Image:       product.Image,
			Difficulty:  product.Difficulty,
			Videos:      vsvideos,
			PreviewURL:  product.PreviewURL,
		})
	}

	return convertedProducts
}
