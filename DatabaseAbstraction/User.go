package DatabaseAbstraction

import (
	"context"
	"time"
)

type User struct {
	IndexID   int
	Username  string
	Password  string
	Balance   int
	CreatedAt time.Time
	UpdatedAt time.Time
	Points    int
}

func (dbc DBConnector) GetAllUsers() ([]User, error) {
	// Get all the users from the database
	rows, err := dbc.DB.Query(context.Background(), "SELECT id, username, password, balance, created_at, updated_at, points FROM users")
	if err != nil {
		return []User{}, err
	}

	// Iterate over the rows and add them to the slice
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.IndexID, &user.Username, &user.Password, &user.Balance, &user.CreatedAt, &user.UpdatedAt, &user.Points)
		if err != nil {
			return []User{}, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (dbc DBConnector) GetUserByUsername(username string) (User, error) {
	// Get the user from the database
	row := dbc.DB.QueryRow(context.Background(), "SELECT id, username, password, balance, created_at, updated_at, points FROM users WHERE username = $1", username)

	var user User
	err := row.Scan(&user.IndexID, &user.Username, &user.Password, &user.Balance, &user.CreatedAt, &user.UpdatedAt, &user.Points)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (dbc DBConnector) GetUserByIndexID(indexID int) (User, error) {
	// Get the user from the database
	row := dbc.DB.QueryRow(context.Background(), "SELECT id, username, password, balance, created_at, updated_at, points FROM users WHERE id = $1", indexID)

	var user User
	err := row.Scan(&user.IndexID, &user.Username, &user.Password, &user.Balance, &user.CreatedAt, &user.UpdatedAt, &user.Points)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (dbc DBConnector) AddUser(username string, password string) error {
	// Add the user to the database
	_, err := dbc.DB.Exec(context.Background(), "INSERT INTO users (username, password) VALUES ($1, $2)", username, password)
	if err != nil {
		return err
	}

	return nil
}

func (dbc DBConnector) DeleteUser(indexID int) error {
	// Delete the user from the database
	_, err := dbc.DB.Exec(context.Background(), "DELETE FROM users WHERE id = $1", indexID)
	if err != nil {
		return err
	}
	// Delete the user's license keys from the database
	_, err = dbc.DB.Exec(context.Background(), "DELETE FROM user_purchases WHERE user_id = $1", indexID)
	if err != nil {
		return err
	}
	// Delete the user's tokens
	_, err = dbc.DB.Exec(context.Background(), "DELETE FROM user_tokens WHERE user_id = $1", indexID)
	if err != nil {
		return err
	}

	return nil
}

func (dbc DBConnector) UpdateUserPassword(indexID int, password string) error {
	_, err := dbc.DB.Exec(context.Background(), "UPDATE users SET password = $1 WHERE id = $2", password, indexID)
	if err != nil {
		return err
	}
	return nil
}

func (dbc DBConnector) UpdateUserUsername(indexID int, username string) error {
	_, err := dbc.DB.Exec(context.Background(), "UPDATE users SET username = $1 WHERE id = $2", username, indexID)
	if err != nil {
		return err
	}
	return nil
}

func (dbc DBConnector) GetOwnedProducts(indexID int) ([]Product, error) {
	// Get all the products from the database
	rows, err := dbc.DB.Query(context.Background(), "SELECT products.* FROM products INNER JOIN user_purchases ON products.id = user_purchases.product_id WHERE user_purchases.user_id = $1", indexID)
	if err != nil {
		return []Product{}, err
	}

	// Iterate over the rows and add them to the slice
	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.IndexID, &product.Name, &product.Description, &product.Price, &product.Image, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return []Product{}, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (dbc DBConnector) IncreaseUserBalance(indexID int, amount int) error {
	_, err := dbc.DB.Exec(context.Background(), "UPDATE users SET balance = balance + $1 WHERE id = $2", amount, indexID)
	if err != nil {
		return err
	}
	return nil
}

func (dbc DBConnector) DecreaseUserBalance(indexID int, amount int) error {
	_, err := dbc.DB.Exec(context.Background(), "UPDATE users SET balance = balance - $1 WHERE id = $2", amount, indexID)
	if err != nil {
		return err
	}

	return nil
}

func (dbc DBConnector) AddOwnedProduct(indexID int, productID int) error {
	_, err := dbc.DB.Exec(context.Background(), "INSERT INTO user_purchases (user_id, product_id) VALUES ($1, $2)", indexID, productID)
	if err != nil {
		return err
	}

	return nil
}
