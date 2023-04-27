package httpserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cafi-dev/iac-gen/pkg/logging"
	"github.com/cafi-dev/iac-gen/pkg/model"
	"github.com/cafi-dev/iac-gen/pkg/tf/tfaws"
	"go.uber.org/zap"
)

type APIHandler struct{}

func NewAPIHandler() *APIHandler {
	return &APIHandler{}
}

func (h *APIHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *APIHandler) GenerateIac(w http.ResponseWriter, r *http.Request) {
	var (
		req    model.ProjectDetails
		logger = logging.GetLogger()
	)

	// decode request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errMsg := "failed to decode request body"
		logger.Error(errMsg, zap.Error(err))
		sendErrResponse(w, http.StatusBadRequest, fmt.Errorf("%s: %w", errMsg, err))
		return
	}

	// create temp dir
	basedir, err := ioutil.TempDir("", "remove-me-*")
	if err != nil {
		errMsg := "failed to create temp directory"
		logger.Error(errMsg, zap.Error(err))
		sendErrResponse(w, http.StatusInternalServerError, fmt.Errorf("%s: %w", errMsg, err))
		return
	}
	// defer os.RemoveAll(basedir)

	// initiate k8s discovery
	aws := tfaws.NewTfAws()
	if err := aws.GenerateIac(basedir, req); err != nil {
		errMsg := "failed to discover k8s resources"
		logger.Error(errMsg, zap.Error(err))
		sendErrResponse(w, http.StatusInternalServerError, fmt.Errorf("%s: %w", errMsg, err))
		return
	}

	logger.Info("successfully generate terraform", zap.String("output", basedir))
	sendJSONResponse(w, http.StatusOK, nil)
}
