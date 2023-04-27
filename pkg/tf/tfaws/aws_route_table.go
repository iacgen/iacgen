package tfaws

import (
	"fmt"
	"path/filepath"

	"github.com/cafi-dev/iac-gen/pkg/constant"
	"github.com/cafi-dev/iac-gen/pkg/tf/tfutils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type TfAwsRouteTable struct {
	Prefix string
	VpcId  hcl.Traversal
}

func (g *TfAwsRouteTable) GetId() hcl.Traversal {
	return hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_route_table",
		},
		hcl.TraverseAttr{
			Name: g.Name(),
		},
		hcl.TraverseAttr{
			Name: "id",
		},
	}
}

func (g *TfAwsRouteTable) createRouteTable(body *hclwrite.Body) {
	body.AppendNewline()
	block := body.AppendNewBlock("resource", []string{"aws_route_table", g.Name()})
	block.Body().SetAttributeTraversal("vpc_id", g.VpcId)
}

func (g *TfAwsRouteTable) Name() string {
	return fmt.Sprintf("%s-route-table", g.Prefix)
}

func (g *TfAwsRouteTable) Generate(basedir string) error {
	// create new empty hcl file object
	hclFile, file, err := tfutils.GetHCLFile(filepath.Join(basedir, constant.MainTf))
	if err != nil {
		return fmt.Errorf("failed to create terraform file for aws provider: %w", err)
	}
	defer file.Close()
	body := hclFile.Body()

	g.createRouteTable(body)

	// write contents to file
	file.Write(hclFile.Bytes())
	return nil
}
