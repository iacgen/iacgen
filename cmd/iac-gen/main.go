package main

import (
	"os"
	"path/filepath"

	"github.com/cafi-dev/iac-gen/pkg/logging"
	"github.com/cafi-dev/iac-gen/pkg/tf/tfaws"
	"go.uber.org/zap"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("generating iac files")

	if err := tfaws.CreateProjectInfra(filepath.Join(os.Getenv("PWD"), ".terraform"), "foo"); err != nil {
		logger.Error("failed to generate iac files", zap.Error(err))
		return
	}
}
