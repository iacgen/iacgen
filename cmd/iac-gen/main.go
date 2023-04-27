package main

import "github.com/cafi-dev/iac-gen/pkg/logging"

func main() {
	logger := logging.GetLogger()
	logger.Info("generating iac files")
}
