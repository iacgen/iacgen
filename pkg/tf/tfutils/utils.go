package tfutils

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func GetExistingHCLFile(filePath string) (*hclwrite.File, *os.File, error) {
	// open file
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file %q: %w", filePath, err)
	}

	// read file
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file %q: %w", filePath, err)
	}

	hclFile, diags := hclwrite.ParseConfig(fileBytes, filePath, hcl.InitialPos)
	if diags.HasErrors() {
		return nil, nil, fmt.Errorf("failed to read existing HCL file %q: %w", filePath, err)
	}

	return hclFile, file, nil
}

func GetHCLFile(filePath string) (*hclwrite.File, *os.File, error) {
	// check if file exists
	if _, err := os.Stat(filePath); err != nil {
		if !os.IsNotExist(err) {
			return nil, nil, fmt.Errorf("failed to check if file %q exists: %w", filePath, err)
		}

		// file path does not exist
		file, err := os.Create(filePath)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create file: %w", err)
		}
		return hclwrite.NewEmptyFile(), file, nil
	}
	return GetExistingHCLFile(filePath)
}
