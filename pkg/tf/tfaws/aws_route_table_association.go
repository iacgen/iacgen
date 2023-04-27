package tfaws

import (
	"fmt"
	"path/filepath"

	"github.com/cafi-dev/iac-gen/pkg/constant"
	"github.com/cafi-dev/iac-gen/pkg/tf/tfutils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type TfAwsRouteTableAssociation struct {
	Prefix       string
	SubnetId     hcl.Traversal
	RouteTableId hcl.Traversal
}

func (g *TfAwsRouteTableAssociation) createRouteTableAssoc(body *hclwrite.Body) {
	body.AppendNewline()
	block := body.AppendNewBlock("resource", []string{"aws_route_table_association", g.Name()})
	block.Body().SetAttributeTraversal("subnet_id", g.SubnetId)
	block.Body().SetAttributeTraversal("route_table_id", g.RouteTableId)
}

func (g *TfAwsRouteTableAssociation) Name() string {
	return fmt.Sprintf("%s-route-table-assoc", g.Prefix)
}

func (g *TfAwsRouteTableAssociation) Generate(basedir string) error {
	// create new empty hcl file object
	hclFile, file, err := tfutils.GetHCLFile(filepath.Join(basedir, constant.MainTf))
	if err != nil {
		return fmt.Errorf("failed to create terraform file for aws provider: %w", err)
	}
	defer file.Close()
	body := hclFile.Body()

	g.createRouteTableAssoc(body)

	// write contents to file
	file.Write(hclFile.Bytes())
	return nil
}
