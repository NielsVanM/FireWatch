package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/context"
	"github.com/nielsvanm/firewatch/core/models"

	log "github.com/sirupsen/logrus"
)

// ChangePassword verifies the old password and tries to set the new password
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	type ChangePasswordRequest struct {
		OldPassword       string `json:"old_password"`
		NewPassword       string `json:"new_password"`
		RepeatNewPassword string `json:"repeat_new_password"`
	}

	// Load data
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Error("Failed to read data from request body")
		NewResp(false, StatusInternalError).Write(w)
		return
	}

	cpr := ChangePasswordRequest{}
	err = json.Unmarshal(data, &cpr)
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

	u := context.Get(r, "user").(*models.Account)

	succ, _ := u.VerifyPassword(cpr.OldPassword)
	if succ {
		if cpr.NewPassword == cpr.RepeatNewPassword {
			u.SetPassword(cpr.NewPassword)
			go u.Save()

			NewResp(true, StatusOkay).Write(w)
			return
		}

		NewResp(false, StatusInvalidCredentials).
			AddData("message", "New passwords are not the same").
			Write(w)
		return

	}
	NewResp(false, StatusInvalidCredentials).
		AddData("message", "Invalid Password").
		Write(w)
	return
}
