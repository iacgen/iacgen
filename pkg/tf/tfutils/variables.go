package tfutils

import (
	"fmt"
	"path/filepath"

	"github.com/cafi-dev/iac-gen/pkg/constant"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func AddVariable(basedir, varname, description string) error {
	// add input variable into variables.tf file
	hclFile, file, err := GetHCLFile(filepath.Join(basedir, constant.VariablesTf))
	if err != nil {
		return fmt.Errorf("failed to get %q HCL file: %w", constant.VariablesTf, err)
	}
	defer file.Close()
	body := hclFile.Body()

	body.AppendNewline()
	block := body.AppendNewBlock("variable", []string{varname})
	block.Body().SetAttributeValue("description", cty.StringVal(description))
	block.Body().SetAttributeTraversal("type", hcl.Traversal{
		hcl.TraverseRoot{Name: "string"},
	})

	file.Write(hclwrite.Format(hclFile.Bytes()))
	return nil
}
