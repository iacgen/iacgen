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

type TfAwsSubnet struct {
	Prefix       string
	VpcId        hcl.Traversal
	CidrBlock    hcl.Traversal
	AZ           hcl.Traversal
	PublicSubnet bool
}

func (g *TfAwsSubnet) GetId() hcl.Traversal {
	return hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_subnet",
		},
		hcl.TraverseAttr{
			Name: g.Name(),
		},
		hcl.TraverseAttr{
			Name: "id",
		},
	}
}

func (g *TfAwsSubnet) createSubnet(body *hclwrite.Body) {
	body.AppendNewline()
	block := body.AppendNewBlock("resource", []string{"aws_subnet", g.Name()})
	block.Body().SetAttributeTraversal("vpc_id", g.VpcId)
	block.Body().SetAttributeTraversal("cidr_block", g.CidrBlock)
	block.Body().SetAttributeTraversal("availability_zone", g.AZ)
	if g.PublicSubnet {
		block.Body().SetAttributeValue("map_public_ip_on_launch", cty.BoolVal(true))
	}
}

func (g *TfAwsSubnet) Name() string {
	if g.PublicSubnet {
		return fmt.Sprintf("%s-public-subnet", g.Prefix)
	}
	return fmt.Sprintf("%s-priv-subnet", g.Prefix)
}

func (g *TfAwsSubnet) Generate(basedir string) error {
	// create new empty hcl file object
	hclFile, file, err := tfutils.GetHCLFile(filepath.Join(basedir, constant.MainTf))
	if err != nil {
		return fmt.Errorf("failed to create terraform file for aws provider: %w", err)
	}
	defer file.Close()
	body := hclFile.Body()

	g.createSubnet(body)

	// write contents to file
	file.Write(hclFile.Bytes())
	return nil
}
