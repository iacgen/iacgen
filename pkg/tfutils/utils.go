package tfutils

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

func CreateEmptyHCLFile(filePath string) (*hclwrite.File, *os.File, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create file: %w", err)
	}
	return hclwrite.NewEmptyFile(), file, nil
}
