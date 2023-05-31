package DatabaseAbstraction

import "context"

type Comment struct {
	IndexID   int
	Username  string
	ProductID int
	Comment   string
}

// Get comments of a product
func (dbc DBConnector) GetCommentsByProductID(productID int) ([]Comment, error) {
	// Get the comments from the database
	rows, err := dbc.DB.Query(context.Background(), "SELECT id, username, product_id, comment FROM product_comments JOIN users ON users.id = user_id WHERE product_id = $1", productID)
	if err != nil {
		return []Comment{}, err
	}

	// Iterate over the rows and add them to the slice
	var comments []Comment

	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.IndexID, &comment.Username, &comment.ProductID, &comment.Comment)
		if err != nil {
			return []Comment{}, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

// Add a comment to a product
func (dbc DBConnector) AddComment(userID int, productID int, comment string) error {
	_, err := dbc.DB.Exec(context.Background(), "INSERT INTO product_comments (user_id, product_id, comment) VALUES ($1, $2, $3)", userID, productID, comment)
	if err != nil {
		return err
	}

	return nil
}
