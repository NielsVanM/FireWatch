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
func (a *Account) NewSession(address string) (*Session, error) {
	// Check if ID is set
	if !a.HasID() {
		return nil, errors.New("failed to create a new session because the user has no ID")
	}

	// Create session
	s := Session{
		-1,
		a.ID,
		tools.RandomToken(64),
		address,
		time.Now().AddDate(0, 0, 14),
	}

	return &s, nil
}

// GetSessions retrieves a list of all sessions currently in the database
func (a *Account) GetSessions() []*Session {
	if !a.HasID() {
		return nil
	}

	rows := database.DB.Query(`
	SELECT * FROM session
	WHERE user_id == $1`, a.ID)

	sessions := []*Session{}

	for rows.Next() {
		ses := Session{}
		rows.Scan(
			&ses.ID,
			&ses.UserID,
			&ses.SessionToken,
			&ses.Address,
			&ses.ExpiryDate,
		)

		sessions = append(sessions, &ses)
	}

	return sessions
}

// GetAccountByUsername retrieves a account from the database and parses
// it into the account struct
func GetAccountByUsername(username string) *Account {
	rows := database.DB.Query(`
	SELECT * FROM account
	WHERE username = $1;`, username)

	a := Account{}
	for rows.Next() {
		rows.Scan(
			&a.ID,
			&a.UserName,
			&a.Password,
		)
	}

	// Check if no user account is found
	if a.UserName == "" {
		return nil
	}

	return &a
}

// Save saves the user to the database.
func (a *Account) Save() {
	rows := database.DB.Query(`
	INSERT INTO account (username, password)
	VALUES ($1, $2)
	ON CONFLICT (username)
	DO UPDATE SET password = $2
	RETURNING id;`, a.UserName, a.Password)

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
	Address      string
	ExpiryDate   time.Time
}

// Verify checks various attributes of the session to check if it's valid
func (s *Session) Verify(address string) bool {
	// Fail if the source IP is different
	if address != s.Address {
		return false
	}

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
	rows := database.DB.Query(`
	SELECT * FROM session
	WHERE token = $1;`, token)

	s := Session{}
	for rows.Next() {
		rows.Scan(
			&s.ID,
			&s.UserID,
			&s.SessionToken,
			&s.Address,
			&s.ExpiryDate,
		)
		break
	}

	return &s
}

// Save stores the session in the database
func (s *Session) Save() {
	s.UpdateExpiryDate()

	rows := database.DB.Query(`
	INSERT INTO session (user_id, token, address, expiry_date)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (token) DO NOTHING;
	`, s.UserID, s.SessionToken, s.Address, s.ExpiryDate)

	for rows.Next() {
		rows.Scan(
			&s.ID,
		)
	}
}

// Delete deletes the session from the database
func (s *Session) Delete() {
	database.DB.Query(`
	DELETE FROM session
	WHERE id = $1`, s.ID)
}
