package AuthenticationManagement

func (am AuthenticationService) AuthenticateUser(username string, password string) (bool, error) {
	// Get the user from the database
	user, err := am.DB.GetUserByUsername(username)
	if err != nil {
		return false, err
	}

	// Verify password
	valid, err := am.ComparePasswords(user.Password, password)
	if err != nil {
		return false, err
	}

	if !valid {
		return false, nil
	}

	return true, nil
}

func (am AuthenticationService) CreateUser(username string, password string) error {
	passwordHash, err := am.HashPassword(password)
	if err != nil {
		return err
	}

	err = am.DB.AddUser(username, passwordHash)
	if err != nil {
		return err
	}

	return nil
}
