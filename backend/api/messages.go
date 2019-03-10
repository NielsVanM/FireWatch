package api

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Constants for common error codes
const (
	StatusOkay               = "okay"
	StatusInternalError      = "internal_error"
	StatusInvalidRequest     = "invalid_request"
	StatusInvalidCredentials = "invalid_credentials"
	StatusInvalidToken       = "invalid_token"
)

// Resp is a response that should be marshaled into a responsewriter
type Resp struct {
	Success    bool                   `json:"success"`
	StatusCode string                 `json:"status_code"`
	Data       map[string]interface{} `json:"data"`
}

// NewResp is the constructor of a Resp
func NewResp(succ bool, statuscode string) *Resp {
	return &Resp{
		succ, statuscode, map[string]interface{}{},
	}
}

// AddData adds a key/value pair to the json response
func (r *Resp) AddData(key string, value interface{}) *Resp {
	r.Data[key] = value

	return r
}

func (r *Resp) Write(w http.ResponseWriter) {
	dat, err := json.Marshal(r)
	if err != nil {
		log.WithFields(log.Fields{
			"resp":  r,
			"error": err.Error(),
		}).Error("Failed to marshal response")
		w.Write([]byte("internal server error"))
	}

	w.Write(dat)
}
