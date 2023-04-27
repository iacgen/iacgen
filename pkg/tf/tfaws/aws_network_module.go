package tfaws

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cafi-dev/iac-gen/pkg/constant"
	"github.com/cafi-dev/iac-gen/pkg/tf/tfutils"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

var networkVariables = []string{"environment", "vpc_cidr", "public_subnet_cidr", "private_subnet_cidr", "availability_zone"}

func (*TfAws) createNetworkModuleCall(basedir string) error {
	// create new empty hcl file object
	hclFile, file, err := tfutils.GetHCLFile(filepath.Join(basedir, constant.MainTf))
	if err != nil {
		return fmt.Errorf("failed to create terraform file for aws provider: %w", err)
	}
	defer file.Close()
	body := hclFile.Body()

	body.AppendNewline()
	block := body.AppendNewBlock("module", []string{"networking"})
	block.Body().SetAttributeValue("source", cty.StringVal("./modules/network"))
	for _, tfvar := range networkVariables {
		block.Body().SetAttributeTraversal(tfvar, hcl.Traversal{
			hcl.TraverseRoot{Name: "var"},
			hcl.TraverseAttr{Name: tfvar},
		})
	}

	// write contents to file
	file.Write(hclFile.Bytes())
	return nil
}

func (*TfAws) CreateNetworkModuleVariables(basedir string) {
	tfutils.AddVariable(basedir, "prefix", "string", "Prefix to be added to every resource created")
	tfutils.AddVariable(basedir, "environment", "string", "Deployment environment (viz. dev, qa, stage, prod)")
	tfutils.AddVariable(basedir, "vpc_cidr", "string", "CIDR block of VPC")
	tfutils.AddVariable(basedir, "public_subnet_cidr", "string", "CIDR block of public subnet")
	tfutils.AddVariable(basedir, "private_subnet_cidr", "string", "CIDR block of private subnet")
	tfutils.AddVariable(basedir, "availability_zone", "string", "Availability Zone (AZ) in which all resources would be deployed")
}

func (t *TfAws) CreateNetworkModule(basedir, prefix string) error {
	// create module directory
	modulesDir := filepath.Join(basedir, "modules", "network")
	if err := os.MkdirAll(modulesDir, 0744); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// create variables for network module
	t.CreateNetworkModuleVariables(modulesDir)

	// create networking infrastructure
	provider := TfAwsProvider{
		AddVersionBlock: true,
	}
	if err := provider.Generate(modulesDir); err != nil {
		return fmt.Errorf("failed to create aws terraform versions for modules: %w", err)
	}

	// create vpc
	vpc := TfAwsVpc{
		Prefix: prefix,
	}
	if err := vpc.Generate(modulesDir); err != nil {
		return fmt.Errorf("failed to create aws_vpc: %w", err)
	}

	// create public subnet
	pubSubnet := TfAwsSubnet{
		Prefix: prefix,
		VpcId:  vpc.GetId(),
		CidrBlock: hcl.Traversal{
			hcl.TraverseRoot{Name: "var"},
			hcl.TraverseAttr{Name: "public_subnet_cidr"},
		},
		AZ: hcl.Traversal{
			hcl.TraverseRoot{Name: "var"},
			hcl.TraverseAttr{Name: "availability_zone"},
		},
		PublicSubnet: true,
	}
	if err := pubSubnet.Generate(modulesDir); err != nil {
		return fmt.Errorf("failed to create public aws_subnet: %w", err)
	}

	// create private subnet
	privSubnet := TfAwsSubnet{
		Prefix: prefix,
		VpcId:  vpc.GetId(),
		CidrBlock: hcl.Traversal{
			hcl.TraverseRoot{Name: "var"},
			hcl.TraverseAttr{Name: "private_subnet_cidr"},
		},
		AZ: hcl.Traversal{
			hcl.TraverseRoot{Name: "var"},
			hcl.TraverseAttr{Name: "availability_zone"},
		},
	}
	if err := privSubnet.Generate(modulesDir); err != nil {
		return fmt.Errorf("failed to create private aws_subnet: %w", err)
	}

	// create internet gateway
	igw := TfAwsInternetGateway{
		Prefix: prefix,
		VPC:    vpc,
	}
	if err := igw.Generate(modulesDir); err != nil {
		return fmt.Errorf("failed to create aws_internet_gateway: %w", err)
	}

	// create eip for nat
	eip := TfAwsEip{
		Prefix: fmt.Sprintf("%s-nat", prefix),
	}
	if err := eip.Generate(modulesDir); err != nil {
		return fmt.Errorf("failed to create aws_eip: %w", err)
	}

	// create nat gateway with eip
	natgw := TfAwsNatGateway{
		Prefix:   prefix,
		EipId:    eip.GetId(),
		SubnetId: pubSubnet.GetId(),
	}
	if err := natgw.Generate(modulesDir); err != nil {
		return fmt.Errorf("failed to create aws_nat_gateway: %w", err)
	}

	// create route table for public subnet
	publicRouteTable := TfAwsRouteTable{
		Prefix: fmt.Sprintf("%s-public", prefix),
		VpcId:  vpc.GetId(),
	}
	if err := publicRouteTable.Generate(modulesDir); err != nil {
		return fmt.Errorf("failed to create route table for public subnet: %w", err)
	}

	// create route table for private subnet
	privateRouteTable := TfAwsRouteTable{
		Prefix: fmt.Sprintf("%s-private", prefix),
		VpcId:  vpc.GetId(),
	}
	if err := privateRouteTable.Generate(modulesDir); err != nil {
		return fmt.Errorf("failed to create route table for private subnet: %w", err)
	}

	// add routes in route tables
	igwRoute := TfAwsRoute{
		Prefix:               fmt.Sprintf("%s-igw", prefix),
		RouteTableId:         publicRouteTable.GetId(),
		DestinationCidrBlock: "0.0.0.0/0",
		GatewayId:            igw.GetId(),
	}
	if err := igwRoute.Generate(modulesDir); err != nil {
		return fmt.Errorf("failed to create aws_route for aws_internet_gateway: %w", err)
	}

	natgwRoute := TfAwsRoute{
		Prefix:               fmt.Sprintf("%s-nat-gw", prefix),
		RouteTableId:         privateRouteTable.GetId(),
		DestinationCidrBlock: "0.0.0.0/0",
		NatGatewayId:         natgw.GetId(),
	}
	if err := natgwRoute.Generate(modulesDir); err != nil {
		return fmt.Errorf("failed to create aws_route for aws_nat_gateway: %w", err)
	}

	// create route table association with subnet
	publicRouteTableAssoc := TfAwsRouteTableAssociation{
		Prefix:       fmt.Sprintf("%s-public-subnet", prefix),
		SubnetId:     pubSubnet.GetId(),
		RouteTableId: publicRouteTable.GetId(),
	}
	if err := publicRouteTableAssoc.Generate(modulesDir); err != nil {
		return fmt.Errorf("failed to create aws_route_table_association for public subnet: %w", err)
	}

	privRouteTableAssoc := TfAwsRouteTableAssociation{
		Prefix:       fmt.Sprintf("%s-priv-subnet", prefix),
		SubnetId:     privSubnet.GetId(),
		RouteTableId: privateRouteTable.GetId(),
	}
	if err := privRouteTableAssoc.Generate(modulesDir); err != nil {
		return fmt.Errorf("failed to create aws_route_table_association for private subnet: %w", err)
	}

	return nil
}
