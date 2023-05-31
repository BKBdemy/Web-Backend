package AuthenticationManagement

import (
	"EntitlementServer/DatabaseAbstraction"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

func (am AuthenticationService) CreateToken(userid int) (string, error) {
	// Get the user from the database
	user, err := am.DB.GetUserByIndexID(userid)

	if err != nil {
		return "", err
	}

	// Use random bytes as salt
	randomBytes := make([]byte, argon2settings.saltLength)
	_, err = rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Create a token by hashing user+their hashed password+the random bytes
	token := fmt.Sprintf("%x", sha256.Sum256([]byte(user.Username+user.Password+string(randomBytes))))

	// by default, tokens expire after 7 days
	err = am.DB.AddToken(user.IndexID, token, time.Now().Add(time.Hour*24*7))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (am AuthenticationService) ValidateToken(token string) (bool, DatabaseAbstraction.User, error) {
	// Get the user from the database
	userToken, err := am.DB.GetTokenByHash(token)
	if err != nil {
		logrus.Errorf("Error getting token from database: %v", err)
		return false, DatabaseAbstraction.User{}, err
	}
	if userToken == (DatabaseAbstraction.Token{}) {
		return false, DatabaseAbstraction.User{}, nil
	}

	// Get the user from the database
	user, err := am.DB.GetUserByIndexID(userToken.UserID)
	if err != nil {
		return false, DatabaseAbstraction.User{}, err
	}

	return true, user, nil
}
