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

	if _, err := file.Write(hclwrite.Format(hclFile.Bytes())); err != nil {
		return fmt.Errorf("failed to write to file %q: %w", file.Name(), err)
	}
	return nil
}

func AddOutput(basedir, varname, description string, traversePath []string) error {
	if len(traversePath) == 0 {
		return fmt.Errorf("empty traverse path for output variable %q", varname)
	}
	// add output variable into output.tf file
	hclFile, file, err := GetHCLFile(filepath.Join(basedir, constant.OutputTf))
	if err != nil {
		return fmt.Errorf("failed to get %q HCL file: %w", constant.OutputTf, err)
	}
	defer file.Close()
	body := hclFile.Body()

	body.AppendNewline()
	block := body.AppendNewBlock("output", []string{varname})
	block.Body().SetAttributeValue("description", cty.StringVal(description))
	block.Body().SetAttributeTraversal("type", hcl.Traversal{
		hcl.TraverseRoot{Name: "string"},
	})

	hclTraversal := hcl.Traversal{
		hcl.TraverseRoot{Name: traversePath[0]},
	}
	for i := range traversePath {
		if i == 0 {
			continue
		}
		hclTraversal = hcl.TraversalJoin(hclTraversal, hcl.Traversal{
			hcl.TraverseAttr{Name: traversePath[i]},
		})
	}
	block.Body().SetAttributeTraversal("value", hclTraversal)

	if _, err := file.Write(hclwrite.Format(hclFile.Bytes())); err != nil {
		return fmt.Errorf("failed to write to file %q: %w", file.Name(), err)
	}
	return nil
}
