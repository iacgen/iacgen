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

type TfAwsVpc struct {
	Prefix string
}

func (g *TfAwsVpc) GetId() hcl.Traversal {
	return hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_vpc",
		},
		hcl.TraverseAttr{
			Name: g.Name(),
		},
		hcl.TraverseAttr{
			Name: "id",
		},
	}
}

func (g *TfAwsVpc) createVpc(body *hclwrite.Body, basedir, vpcname string) {
	body.AppendNewline()
	block := body.AppendNewBlock("resource", []string{"aws_vpc", vpcname})
	block.Body().SetAttributeTraversal("cidr_block", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "var",
		},
		hcl.TraverseAttr{
			Name: "vpc_cidr",
		},
	})
	block.Body().SetAttributeValue("enable_dns_hostnames", cty.BoolVal(true))
	block.Body().SetAttributeValue("enable_dns_support", cty.BoolVal(true))
}

func (g *TfAwsVpc) Name() string {
	return fmt.Sprintf("%s-vpc", g.Prefix)
}

func (g *TfAwsVpc) Generate(basedir string) error {
	vpcname := g.Name()
	// create new empty hcl file object
	hclFile, file, err := tfutils.GetHCLFile(filepath.Join(basedir, constant.MainTf))
	if err != nil {
		return fmt.Errorf("failed to create terraform file for aws provider: %w", err)
	}
	defer file.Close()
	body := hclFile.Body()

	g.createVpc(body, basedir, vpcname)

	// write contents to file
	file.Write(hclFile.Bytes())
	return nil
}
