package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func ErrorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]any)
	resp["error"] = message
	resp["code"] = httpStatusCode
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func DataResponse(w http.ResponseWriter, body any, httpStatusCode int) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatusCode)

	jsonResp, err := json.Marshal(body)

	if err != nil {
		fmt.Println(err.Error())
		return err

	}

	w.Write(jsonResp)

	return nil
}
