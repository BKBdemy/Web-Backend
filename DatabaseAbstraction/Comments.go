package DatabaseAbstraction

import (
	"context"
	"time"
)

type Comment struct {
	IndexID   int
	Username  string
	ProductID int
	Comment   string
	CreatedAt time.Time
}

// Get comments of a product
func (dbc DBConnector) GetCommentsByProductID(productID int) ([]Comment, error) {
	// Get the comments from the database
	rows, err := dbc.DB.Query(context.Background(), "SELECT product_comments.id, username, product_comments.course_id, comment, product_comments.created_at FROM product_comments JOIN users ON users.id = user_id WHERE course_id = $1 ORDER BY product_comments.id DESC", productID)
	if err != nil {
		return []Comment{}, err
	}

	// Iterate over the rows and add them to the slice
	var comments []Comment

	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.IndexID, &comment.Username, &comment.ProductID, &comment.Comment, &comment.CreatedAt)
		if err != nil {
			return []Comment{}, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

// AddComment adds a comment to a product
func (dbc DBConnector) AddComment(userID int, productID int, comment string) error {
	_, err := dbc.DB.Exec(context.Background(), "INSERT INTO product_comments (user_id, course_id, comment) VALUES ($1, $2, $3)", userID, productID, comment)
	if err != nil {
		return err
	}

	return nil
}
