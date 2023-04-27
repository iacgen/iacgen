package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/cafi-dev/iac-gen/pkg/logging"
	"go.uber.org/zap"
)

func sendResponse(w http.ResponseWriter, statusCode int, body []byte) {
	w.WriteHeader(statusCode)
	if _, err := w.Write(body); err != nil {
		logging.GetLogger().Error("failed to write http response", zap.Error(err))
	}
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	resp, err := json.Marshal(body)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}
	sendResponse(w, statusCode, resp)
}

type errResp struct {
	Error string `json:"error"`
}

func sendErrResponse(w http.ResponseWriter, statusCode int, err error) {
	sendJSONResponse(w, statusCode, errResp{
		Error: err.Error(),
	})
}
