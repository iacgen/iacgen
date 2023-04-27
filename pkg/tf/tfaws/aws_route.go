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

type TfAwsRoute struct {
	Prefix               string
	RouteTableId         hcl.Traversal
	DestinationCidrBlock string
	GatewayId            hcl.Traversal
	NatGatewayId         hcl.Traversal
}

func (g *TfAwsRoute) createRoute(body *hclwrite.Body) {
	body.AppendNewline()
	block := body.AppendNewBlock("resource", []string{"aws_route", g.Name()})
	block.Body().SetAttributeTraversal("route_table_id", g.RouteTableId)
	block.Body().SetAttributeValue("destination_cidr_block", cty.StringVal(g.DestinationCidrBlock))
	if len(g.GatewayId) > 0 {
		block.Body().SetAttributeTraversal("gateway_id", g.GatewayId)
	}
	if len(g.NatGatewayId) > 0 {
		block.Body().SetAttributeTraversal("nat_gateway_id", g.NatGatewayId)
	}
}

func (g *TfAwsRoute) Name() string {
	return fmt.Sprintf("%s-route", g.Prefix)
}

func (g *TfAwsRoute) Generate(basedir string) error {
	// create new empty hcl file object
	hclFile, file, err := tfutils.GetHCLFile(filepath.Join(basedir, constant.MainTf))
	if err != nil {
		return fmt.Errorf("failed to create terraform file for aws provider: %w", err)
	}
	defer file.Close()
	body := hclFile.Body()

	g.createRoute(body)

	// write contents to file
	file.Write(hclFile.Bytes())
	return nil
}
