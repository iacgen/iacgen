package tfaws

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cafi-dev/iac-gen/pkg/model"
	cp "github.com/otiai10/copy"
)

type ITfAws interface {
	CreateRemoteBackendConfig(basedir string) error
	GenerateIac(basedir string, projectDetails model.ProjectDetails) error
}

type TfAws struct{}

func NewTfAws() *TfAws {
	return &TfAws{}
}

func (t *TfAws) GenerateIac(basedir string, projectDetails model.ProjectDetails) error {
	// copy terraform templates
	templateDir := filepath.Join(os.Getenv("PWD"), "terraform")

	for _, project := range projectDetails.Projects {
		outputDir := filepath.Join(basedir, project.Metadata.Name)
		if err := cp.Copy(templateDir, outputDir); err != nil {
			return fmt.Errorf("failed to copy template directory: %w", err)
		}

		// create ecs resources
		for _, svc := range project.Services {
			app := TfAwsECSTemplate{
				AppName:       svc.Name,
				Image:         svc.Image,
				ContainerPort: svc.Ports[0].Listen,
			}
			if err := app.RenderTemplate(outputDir); err != nil {
				return fmt.Errorf("failed to render ecs template: %w", err)
			}
		}
	}
	return nil
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
