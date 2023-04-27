package tfaws

import (
	"fmt"
	"path/filepath"

	"github.com/cafi-dev/iac-gen/pkg/constant"
	"github.com/cafi-dev/iac-gen/pkg/tf/tfutils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type TfAwsNatGateway struct {
	Prefix   string
	EipId    hcl.Traversal
	SubnetId hcl.Traversal
}

func (g *TfAwsNatGateway) GetId() hcl.Traversal {
	return hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_nat_gateway",
		},
		hcl.TraverseAttr{
			Name: g.Name(),
		},
		hcl.TraverseAttr{
			Name: "id",
		},
	}
}

func (g *TfAwsNatGateway) createNatGw(body *hclwrite.Body) {
	body.AppendNewline()
	block := body.AppendNewBlock("resource", []string{"aws_nat_gateway", g.Prefix})
	block.Body().SetAttributeTraversal("allocation_id", g.EipId)
	block.Body().SetAttributeTraversal("subnet_id", g.SubnetId)
}

func (g *TfAwsNatGateway) Name() string {
	return fmt.Sprintf("%s-nat-gw", g.Prefix)
}

func (g *TfAwsNatGateway) Generate(basedir string) error {
	// create new empty hcl file object
	hclFile, file, err := tfutils.GetHCLFile(filepath.Join(basedir, constant.MainTf))
	if err != nil {
		return fmt.Errorf("failed to create terraform file for aws provider: %w", err)
	}
	defer file.Close()
	body := hclFile.Body()

	g.createNatGw(body)

	// write contents to file
	file.Write(hclFile.Bytes())
	return nil
}
