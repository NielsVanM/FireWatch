package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/nielsvanm/firewatch/core/models"
	log "github.com/sirupsen/logrus"
)

// Login view handles the login sequence of a user, it receives a username and
// password and returns a newly generated token for the session
func Login(w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Load data
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to read data from request body")
		NewResp(false, StatusInternalError).Write(w)
		return
	}

	// Parse login request
	lr := LoginRequest{}
	err = json.Unmarshal(data, &lr)
	if err != nil {
		log.WithFields(log.Fields{
			"request_data": string(data),
			"error":        err.Error(),
		}).Warn("Failed to unmarshal HTTP request body")

		NewResp(false, StatusInvalidRequest).
			AddData("message", "no username/password provided").
			Write(w)
		return
	}

	// Retrieve account
	acc := models.GetAccountByUsername(lr.Username)
	if acc == nil {
		log.Warn("Failed to retrieve user account")

		NewResp(false, StatusInvalidCredentials).Write(w)
		return
	}

	// Verify password and create session
	succ, err := acc.VerifyPassword(lr.Password)
	if !succ || err != nil {
		// Invalid creds
		NewResp(false, StatusInvalidCredentials).Write(w)
		return
	}

	// Valid credentials
	sess, _ := models.NewSession(acc)
	go sess.Save()
	NewResp(true, StatusOkay).
		AddData("token", sess.SessionToken).
		Write(w)

}

// Logout deletes the retreived token from the database and therefore ending
// a session
func Logout(w http.ResponseWriter, r *http.Request) {
	type LogoutRequest struct {
		Token string `json:"token"`
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		NewResp(false, StatusInternalError).
			AddData("message", "failed to read request body").
			Write(w)
		return
	}

	lr := LogoutRequest{}
	err = json.Unmarshal(data, &lr)
	if err != nil {
		NewResp(false, StatusInvalidRequest).
			Write(w)
		return
	}

	sess := models.GetSessionByToken(lr.Token)
	if sess == nil {
		NewResp(false, StatusInvalidToken).Write(w)
		return
	}

	sess.Delete()
	NewResp(true, StatusOkay).Write(w)
}
