package models

import (
	"database/sql"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/nielsvanm/firewatch/internal/database"
	"github.com/nielsvanm/firewatch/internal/tools"
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

// NewSession creates a new session for the user with a session token and IP
// address combination. The expiry date is set to two weeks from now.
func (a *Account) NewSession() (*Session, error) {
	// Check if ID is set
	if !a.HasID() {
		return nil, ErrAccountHasNoID
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
		log.Warn("Failed to retrieve sessions from the database", err.Error())
		return nil
	}

	defer rows.Close()

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

// DeleteAllSessions deletes all the sessions from the current account
func (a *Account) DeleteAllSessions() {
	database.DB.Exec(`
	DELETE FROM session
	WHERE user_id = $1`, a.ID)
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
		log.Info("Failed to retrieve token from database", err.Error())
		return nil
	}

	defer rows.Close()

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
		log.Warning("Failed to save session to the database", err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(
			&s.ID,
		)
	}
}

// Delete deletes the session from the database
func (s *Session) Delete() {
	database.DB.Exec(`
	DELETE FROM session
	WHERE id = $1`, s.ID)
}
