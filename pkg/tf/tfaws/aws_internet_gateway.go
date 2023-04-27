package tfaws

import (
	"fmt"
	"path/filepath"

	"github.com/cafi-dev/iac-gen/pkg/constant"
	"github.com/cafi-dev/iac-gen/pkg/tf/tfutils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type TfAwsInternetGateway struct {
	Prefix string
	VPC    TfAwsVpc
}

func (g *TfAwsInternetGateway) GetId() hcl.Traversal {
	return hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_internet_gateway",
		},
		hcl.TraverseAttr{
			Name: g.Name(),
		},
		hcl.TraverseAttr{
			Name: "id",
		},
	}
}

func (g *TfAwsInternetGateway) createIGW(body *hclwrite.Body, igwname string) {
	body.AppendNewline()
	block := body.AppendNewBlock("resource", []string{"aws_internet_gateway", igwname})
	block.Body().SetAttributeTraversal("vpc_id", g.VPC.GetId())
}

func (g *TfAwsInternetGateway) Name() string {
	return fmt.Sprintf("%s-igw", g.Prefix)
}

func (g *TfAwsInternetGateway) Generate(basedir string) error {
	igname := g.Name()
	// create new empty hcl file object
	hclFile, file, err := tfutils.GetHCLFile(filepath.Join(basedir, constant.MainTf))
	if err != nil {
		return fmt.Errorf("failed to create terraform file for aws provider: %w", err)
	}
	defer file.Close()
	body := hclFile.Body()

	g.createIGW(body, igname)

	// write contents to file
	file.Write(hclFile.Bytes())
	return nil
}
