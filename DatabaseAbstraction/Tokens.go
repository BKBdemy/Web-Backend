package DatabaseAbstraction

import (
	"context"
	"time"
)

type Token struct {
	IndexID int
	UserID  int
	Token   string
	Expiry  time.Time
}

func (dbc DBConnector) GetTokenByTokenID(tokenID string) (Token, error) {
	// Get the token from the database
	row := dbc.DB.QueryRow(context.Background(), "SELECT * FROM user_tokens WHERE id = $1 AND expiry > now()", tokenID)

	var userToken Token
	err := row.Scan(&userToken.IndexID, &userToken.UserID, &userToken.Token, &userToken.Expiry)
	if err != nil {
		return Token{}, err
	}

	return userToken, nil
}

func (dbc DBConnector) GetTokenByHash(token string) (Token, error) {
	// Get the token from the database
	row := dbc.DB.QueryRow(context.Background(), "SELECT * FROM user_tokens WHERE token = $1 AND expiry > now()", token)

	var userToken Token
	err := row.Scan(&userToken.IndexID, &userToken.UserID, &userToken.Token, &userToken.Expiry)
	if err != nil {
		return Token{}, err
	}

	return userToken, nil
}

func (dbc DBConnector) AddToken(userID int, token string, expiry time.Time) error {
	// Add the token to the database
	_, err := dbc.DB.Exec(context.Background(), "INSERT INTO user_tokens (user_id, token, expiry) VALUES ($1, $2, $3)", userID, token, expiry)
	if err != nil {
		return err
	}

	return nil
}

func (dbc DBConnector) DeleteToken(tokenID int) error {
	// Delete the token from the database
	_, err := dbc.DB.Exec(context.Background(), "DELETE FROM user_tokens WHERE id = $1", tokenID)
	if err != nil {
		return err
	}

	return nil
}
