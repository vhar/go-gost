package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Err struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func ErrorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]Err)

	resp["error"] = Err{
		Message: message,
		Code:    httpStatusCode,
	}

	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func DataResponse(w http.ResponseWriter, body any, httpStatusCode int) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]any)

	resp["payload"] = body

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err.Error())
		return err

	}

	w.Write(jsonResp)

	return nil
}
