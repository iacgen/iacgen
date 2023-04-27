package tfaws

import (
	"fmt"
	"path/filepath"

	"github.com/cafi-dev/iac-gen/pkg/constant"
	"github.com/cafi-dev/iac-gen/pkg/tf/tfutils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type TfAwsEip struct {
	Prefix string
}

func (g *TfAwsEip) GetId() hcl.Traversal {
	return hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_eip",
		},
		hcl.TraverseAttr{
			Name: g.Name(),
		},
		hcl.TraverseAttr{
			Name: "id",
		},
	}
}

func (g *TfAwsEip) createEIP(body *hclwrite.Body) {
	body.AppendNewline()
	block := body.AppendNewBlock("resource", []string{"aws_eip", g.Name()})
	block.Body().SetAttributeValue("vpc", cty.BoolVal(true))
}

func (g *TfAwsEip) Name() string {
	return fmt.Sprintf("%s-eip", g.Prefix)
}

func (g *TfAwsEip) Generate(basedir string) error {
	// create new empty hcl file object
	hclFile, file, err := tfutils.GetHCLFile(filepath.Join(basedir, constant.MainTf))
	if err != nil {
		return fmt.Errorf("failed to create terraform file for aws provider: %w", err)
	}
	defer file.Close()
	body := hclFile.Body()

	g.createEIP(body)

	// write contents to file
	file.Write(hclFile.Bytes())
	return nil
}
