package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/nielsvanm/firewatch/internal/database"
	"github.com/nielsvanm/firewatch/internal/tools"
	"golang.org/x/crypto/bcrypt"
)

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
		fmt.Println("Failed to hash user password")
		return
	}
	a.Password = hashedPassword
}

// VerifyPassword verifies a plaintext password against the hashed password
func (a *Account) VerifyPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(a.Password, []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}

// HasID is a helper function for checking if the user has been stored in the
// database
func (a *Account) HasID() bool {
	return a.ID != -1
}

// NewSession creates a new session for the user with a session token and IP
// address combination. The expiry date is set to two weeks from now.
func (a *Account) NewSession() (*Session, error) {
	// Check if ID is set
	if !a.HasID() {
		return nil, errors.New("failed to create a new session because the user has no ID")
	}

	// Create session
	s := Session{
		-1,
		a.ID,
		tools.RandomToken(64),
		time.Now().AddDate(0, 0, 14),
	}

	return &s, nil
}

// GetSessions retrieves a list of all sessions currently in the database
func (a *Account) GetSessions() []*Session {
	if !a.HasID() {
		return nil
	}

	rows, err := database.DB.Query(`
	SELECT * FROM session
	WHERE user_id == $1`, a.ID)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	sessions := []*Session{}

	for rows.Next() {
		ses := Session{}
		rows.Scan(
			&ses.ID,
			&ses.UserID,
			&ses.SessionToken,
			&ses.ExpiryDate,
		)

		sessions = append(sessions, &ses)
	}

	return sessions
}

// GetAccountByUsername retrieves a account from the database and parses
// it into the account struct
func GetAccountByUsername(username string) *Account {
	rows, err := database.DB.Query(`
	SELECT * FROM account
	WHERE username = $1;`, username)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

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
	rows, err := database.DB.Query(`
	INSERT INTO account (username, password)
	VALUES ($1, $2)
	ON CONFLICT (username)
	DO UPDATE SET password = $2
	RETURNING id;`, a.UserName, a.Password)

	if err != nil {
		fmt.Println(err.Error())
	}

	for rows.Next() {
		rows.Scan(
			&a.ID,
		)
	}
}

// Session is a structure representing a source ip an d token, it can be used
// to keep track of whether the user is logged in.
type Session struct {
	ID           int
	UserID       int
	SessionToken string
	ExpiryDate   time.Time
}

// Verify checks various attributes of the session to check if it's valid
func (s *Session) Verify() bool {
	// Check expiry date, if expired delete from DB
	if time.Until(s.ExpiryDate).Minutes() < 0.0 {
		s.Delete()
		return false
	}

	// If all checks were succesful we can safely assume the token is okay
	return true
}

// UpdateExpiryDate updates the expiry date of the session to be two weeks from
// now
func (s *Session) UpdateExpiryDate() {
	s.ExpiryDate = time.Now().AddDate(0, 0, 14)
}

// GetSessionByToken retrieves a session object based on the provided token
func GetSessionByToken(token string) *Session {
	rows, err := database.DB.Query(`
	SELECT * FROM session
	WHERE token = $1;`, token)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	s := Session{}
	for rows.Next() {
		rows.Scan(
			&s.ID,
			&s.UserID,
			&s.SessionToken,
			&s.ExpiryDate,
		)
		break
	}

	return &s
}

// Save stores the session in the database
func (s *Session) Save() {
	s.UpdateExpiryDate()

	rows, err := database.DB.Query(`
	INSERT INTO session (user_id, token, expiry_date)
	VALUES ($1, $2, $3)
	ON CONFLICT (token)
	DO 
		UPDATE SET expiry_date = $3;
	`, s.UserID, s.SessionToken, s.ExpiryDate)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for rows.Next() {
		rows.Scan(
			&s.ID,
		)
	}
}

// Delete deletes the session from the database
func (s *Session) Delete() {
	_, err := database.DB.Query(`
	DELETE FROM session
	WHERE id = $1`, s.ID)
	if err != nil {
		fmt.Println(err.Error())
	}
}
