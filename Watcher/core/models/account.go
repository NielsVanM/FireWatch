package models

import (
	"database/sql"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/nielsvanm/firewatch/core/database"
	"golang.org/x/crypto/bcrypt"
)

// ErrAccountHasNoID must be returned when an account doesn't have an ID, i.e.
// it isn't saved in the database
var ErrAccountHasNoID = errors.New("Account has no ID")

// Account is an end-user in the application, it provides a username/password based
// authentication method.
type Account struct {
	ID       int
	UserName string
	Password []byte
}

// NewAccount creates a new Account and hashes the password
func NewAccount(username, password string) *Account {
	a := Account{
		-1,
		username,
		[]byte{},
	}

	a.SetPassword(password)

	return &a
}

// SetPassword hashes the input password and sets it as the new password
func (a *Account) SetPassword(password string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Warn("Failed to hash user password", err.Error())
		return
	}
	a.Password = hashedPassword
}

// VerifyPassword verifies a plaintext password against the hashed password
func (a *Account) VerifyPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(a.Password, []byte(password))
	if err != nil {

		// Log any other error than mismatch
		if err != bcrypt.ErrMismatchedHashAndPassword {
			log.Warn("Failed to verify password", err.Error())
		}
		return false, err
	}

	return true, nil
}

// HasID is a helper function for checking if the user has been stored in the
// database
func (a *Account) HasID() bool {
	return a.ID != -1
}

// GetAccountByID retrieves a account form the database by the provided by the ID
// it returns a pointer to the account if it exists, if the account can't be found
// it returns nil
func GetAccountByID(id int) *Account {
	rows, err := database.DB.Query(`
	SELECT * FROM account
	WHERE id = $1
	LIMIT 1;`, id)
	if err != nil {
		log.Warn("Failed to retrieve user from database", err.Error())
		return nil
	}

	defer rows.Close()

	a := Account{}
	for rows.Next() {
		rows.Scan(
			&a.ID,
			&a.UserName,
			&a.Password,
		)
		break
	}

	// Check if the username and password are still empty
	if a.UserName == "" || a.Password == nil {
		return nil
	}

	return &a
}

// GetAccountByUsername retrieves a account from the database and parses
// it into the account struct. The function returns a pointer to the Account
// object if the account is found in the database, otherwise nil is returned
func GetAccountByUsername(username string) *Account {
	rows, err := database.DB.Query(`
	SELECT * FROM account
	WHERE username = $1;`, username)

	if err != nil {
		log.Info("Failed to retrieve user from database", err.Error())
		return nil
	}

	defer rows.Close()

	a := Account{}
	for rows.Next() {
		rows.Scan(
			&a.ID,
			&a.UserName,
			&a.Password,
		)
		break
	}

	// Check if the username and password are still empty
	if a.UserName == "" || a.Password == nil {
		return nil
	}

	// User exists, return it
	return &a
}

// Save saves the user to the database.
func (a *Account) Save() {
	var rows *sql.Rows
	var err error

	if a.ID == -1 {
		rows, err = database.DB.Query(`
		INSERT INTO account (username, password)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
		RETURNING id;`, a.UserName, a.Password)

		if err != nil {
			log.Warning("Failed to save account to the database", err.Error())
			return
		}
	} else {
		rows, err = database.DB.Query(`
		UPDATE account
		SET password=$1
		WHERE id = $2;`, a.Password, a.ID)
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(
			&a.ID,
		)
	}
}

// Delete removes the account from the database along with it's sessions, the
// object still exists in memory
func (a *Account) Delete() {
	// Delete sessions
	database.DB.Exec(`
	DELETE FROM session
	WHERE user_id = $1`, a.ID)

	// Delete account
	database.DB.Exec(`
	DELETE FROM account
	WHERE id = $1;`, a.ID)
}
