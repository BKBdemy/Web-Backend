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
	Difficulty  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (dbc DBConnector) GetAllProducts() ([]Product, error) {
	// Get all the products from the database
	rows, err := dbc.DB.Query(context.Background(), "SELECT id, name, description, price, image, created_at, updated_at, difficulty FROM products")
	if err != nil {
		return []Product{}, err
	}

	// Iterate over the rows and add them to the slice
	var products []Product

	for rows.Next() {
		var product Product
		err := rows.Scan(&product.IndexID, &product.Name, &product.Description, &product.Price, &product.Image, &product.CreatedAt, &product.UpdatedAt, &product.Difficulty)
		if err != nil {
			return []Product{}, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (dbc DBConnector) GetProductByIndexID(indexID int) (Product, error) {
	// Get the product from the database
	row := dbc.DB.QueryRow(context.Background(), "SELECT id, name, description, price, image, created_at, updated_at, difficulty FROM products WHERE id = $1", indexID)

	var product Product
	err := row.Scan(&product.IndexID, &product.Name, &product.Description, &product.Price, &product.Image, &product.CreatedAt, &product.UpdatedAt, &product.Difficulty)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

func (dbc DBConnector) AddProduct(NewProduct Product) (int, error) {
	// Insert the product into the database
	var indexID int
	err := dbc.DB.QueryRow(context.Background(), "INSERT INTO products (name, description, price, image) VALUES ($1, $2, $3, $4) RETURNING id", NewProduct.Name, NewProduct.Description, NewProduct.Price, NewProduct.Image).Scan(&indexID)
	if err != nil {
		return -1, err
	}

	return indexID, nil
}

func (dbc DBConnector) GetProductVideos(indexID int) ([]Video, error) {
	// Get all the videos from the database
	rows, err := dbc.DB.Query(context.Background(), "SELECT id, name, description, points, thumbnail, filename FROM video WHERE parent_product_id = $1", indexID)
	if err != nil {
		return []Video{}, err
	}

	// Iterate over the rows and add them to the slice
	var videos []Video

	for rows.Next() {
		var video Video
		err := rows.Scan(&video.IndexID, &video.Name, &video.Description, &video.Points, &video.Thumbnail, &video.Filename)
		if err != nil {
			return []Video{}, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}
