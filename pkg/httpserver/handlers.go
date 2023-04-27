package httpserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cafi-dev/iac-gen/pkg/archiver"
	"github.com/cafi-dev/iac-gen/pkg/logging"
	"github.com/cafi-dev/iac-gen/pkg/model"
	"github.com/cafi-dev/iac-gen/pkg/tf/tfaws"
	"go.uber.org/zap"
)

type APIHandler struct {
	tgz archiver.Archiver
}

func NewAPIHandler() *APIHandler {
	return &APIHandler{
		tgz: archiver.NewZIP(),
	}
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
	defer os.RemoveAll(basedir)

	// initiate k8s discovery
	aws := tfaws.NewTfAws()
	if err := aws.GenerateIac(basedir, req); err != nil {
		errMsg := "failed to discover k8s resources"
		logger.Error(errMsg, zap.Error(err))
		sendErrResponse(w, http.StatusInternalServerError, fmt.Errorf("%s: %w", errMsg, err))
		return
	}

	// creat tgz file
	filePath := filepath.Join(os.TempDir(), fmt.Sprintf("%s.zip", filepath.Base(basedir)))
	if err := h.tgz.Compress(filePath, basedir); err != nil {
		errMsg := "failed to create tgz file"
		logger.Error(errMsg, zap.Error(err))
		sendErrResponse(w, http.StatusInternalServerError, fmt.Errorf("%s: %w", errMsg, err))
		return
	}
	defer os.RemoveAll(filePath)

	logger.Info("successfully generate terraform", zap.String("output", filePath))
	sendFileResponse(w, r, filePath)
}
