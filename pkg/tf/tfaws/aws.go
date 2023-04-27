package tfaws

import (
	"fmt"
	"os"
	"path/filepath"
)

type ITfAws interface {
	CreateRemoteBackendConfig(basedir string) error
}

type TfAws struct{}

func NewTfAws() *TfAws {
	return &TfAws{}
}

func (*TfAws) CreateRemoteBackendConfig(basedir string) error {
	prefix := "remote-backend"
	dir := filepath.Join(basedir, "remote-backend")
	if err := os.MkdirAll(dir, 0744); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// create aws provider
	provider := TfAwsProvider{
		AddProviderBlock: true,
		AddVersionBlock:  true,
	}
	if err := provider.Generate(dir); err != nil {
		return fmt.Errorf("failed to create aws terraform provider for remote backend: %w", err)
	}

	// create s3 bucket storing terraform state
	s3 := TfAwsS3{
		Prefix:               prefix,
		EnableVersioning:     true,
		EnableEncryption:     true,
		RestrictPublicAccess: true,
	}
	if err := s3.Generate(dir); err != nil {
		return fmt.Errorf("failed to create s3 bucket for aws remote backend: %w", err)
	}

	// create dynamodb for locking
	dynamodb := TfAwsDynamodb{
		Prefix: prefix,
	}
	dynamodb.Generate(dir)

	return nil
}

func (t *TfAws) CreateMain(basedir string) error {
	// create aws provider
	provider := TfAwsProvider{
		AddProviderBlock: true,
		AddVersionBlock:  true,
	}
	if err := provider.Generate(basedir); err != nil {
		return fmt.Errorf("failed to create aws terraform provider for remote backend: %w", err)
	}

	// create variables for networking module
	t.CreateNetworkModuleVariables(basedir)

	// create main.tf calling network module
	if err := t.createNetworkModuleCall(basedir); err != nil {
		return fmt.Errorf("failed to create network module call: %w", err)
	}

	return nil
}

func CreateProjectInfra(basedir, projectname string) error {
	dir := filepath.Join(basedir, projectname)
	tfAws := NewTfAws()
	if err := tfAws.CreateRemoteBackendConfig(dir); err != nil {
		return fmt.Errorf("failed to create remote backend for aws terraform: %w", err)
	}
	if err := tfAws.CreateNetworkModule(dir, projectname); err != nil {
		return fmt.Errorf("failed to create project's network infrastructure terraform module: %w", err)
	}
	if err := tfAws.CreateMain(dir); err != nil {
		return fmt.Errorf("failed to create project's main module: %w", err)
	}
	return nil
}
