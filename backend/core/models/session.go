package models

import (
	"time"

	"github.com/nielsvanm/firewatch/core/database"
	"github.com/nielsvanm/firewatch/core/tools"
	log "github.com/sirupsen/logrus"
)

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

	// Check if we actually got a result
	if s.SessionToken == "" {
		return nil
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

// NewSession creates a new session for the user with a session token and IP
// address combination. The expiry date is set to two weeks from now.
func NewSession(a *Account) (*Session, error) {
	// Check if ID is set
	if !a.HasID() {
		return nil, ErrAccountHasNoID
	}

	// Create session
	s := Session{
		-1,
		a.ID,
		tools.RandomToken(64),
		time.Now().AddDate(0, 0, 1),
	}

	return &s, nil
}

// GetSessions retrieves a list of all sessions currently in the database
func GetSessions(a *Account) []*Session {
	if !a.HasID() {
		return nil
	}

	rows, err := database.DB.Query(`
	SELECT * FROM session
	WHERE user_id == $1`, a.ID)
	if err != nil {
		log.Warn("Failed to retrieve sessions from the database ", err.Error())
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
func DeleteAllSessions(a *Account) {
	database.DB.Exec(`
	DELETE FROM session
	WHERE user_id = $1`, a.ID)
}
