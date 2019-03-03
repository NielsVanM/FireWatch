package models

import (
	"database/sql"

	"github.com/nielsvanm/firewatch/core/database"
	"github.com/nielsvanm/firewatch/core/tools"
	log "github.com/sirupsen/logrus"
)

// APIToken provides access to the API of the program
type APIToken struct {
	ID     int
	UserID int
	Token  string
}

// NewToken generates a new API token for the provided user
func NewToken(a *Account) *APIToken {
	at := APIToken{
		-1, a.ID, tools.RandomToken(128),
	}

	return &at
}

// GetAllTokens returns a list of tokens associated with the provided account
func GetAllAPITokensByAccount(a *Account) []*APIToken {
	rows, err := database.DB.Query(`
	SELECT * FROM apitoken
	WHERE user_id = $1;`, a.ID)
	if err != nil {
		log.Warn("Failed to retrieve api tokens from database ", err.Error())
		return nil
	}

	tokenList := []*APIToken{}
	for rows.Next() {
		at := APIToken{}
		err := rows.Scan(
			&at.ID,
			&at.UserID,
			&at.Token,
		)
		if err != nil {
			log.Warn("Failed to scan row into struct")
			continue
		}

		tokenList = append(tokenList, &at)
	}

	return tokenList
}

func GetUserByAPIToken(token string) {

}

// Save either inserts a new entry in the api token table or updates an existing
// record
func (at *APIToken) Save() {
	var rows *sql.Rows
	var err error

	if at.ID == -1 {
		// Insert new
		rows, err = database.DB.Query(`
		INSERT INTO apitoken (user_id, token)
		VALUES ($1, $2)
		RETURING id;`, at.UserID, at.Token)
		if err != nil {
			log.Warn(err.Error())
			return
		}

	} else {
		// Update
		rows, err = database.DB.Query(`
		UPDATE apitoken
		SET user_id = $1, token = $2
		WHERE id = $3`, at.UserID, at.Token, at.ID)
		if err != nil {
			log.Warn(err.Error())
			return
		}
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(
			&at.ID,
		)
	}
}
