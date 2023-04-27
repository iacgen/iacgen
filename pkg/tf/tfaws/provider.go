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

const tfAwsProvider TFAWSResType = "provider"

func init() {
	addResourceType(tfAwsProvider, NewTfAwsProvider())
}

type TfAwsProvider struct{}

func NewTfAwsProvider() *TfAwsProvider {
	return &TfAwsProvider{}
}

func (*TfAwsProvider) addVersionsBlock(body *hclwrite.Body) {
	body.AppendNewline()
	// add terraform block for required versions
	terraformBlock := body.AppendNewBlock("terraform", nil)
	terraformBlock.Body().SetAttributeValue("required_versions", cty.StringVal(">= 1.0.0, < 2.0.0"))

	// add required providers block
	requiredProvidersBlock := terraformBlock.Body().AppendNewBlock("required_providers", nil)
	requiredProvidersBlock.Body().SetAttributeValue("aws", cty.MapVal(
		map[string]cty.Value{
			"source":  cty.StringVal("hashicorp/aws"),
			"version": cty.StringVal("~> 4.0"),
		},
	))
}

func (*TfAwsProvider) addProviderBlock(basedir string, body *hclwrite.Body) error {
	body.AppendNewline()
	provider := body.AppendNewBlock("provider", []string{"aws"})
	provider.Body().SetAttributeTraversal("region", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "var",
		},
		hcl.TraverseAttr{
			Name: "aws_region",
		},
	})
	return tfutils.AddVariable(basedir, "aws_region", "AWS region")
}

func (g *TfAwsProvider) Generate(basedir, name string) error {
	// create new empty hcl file object
	hclFile, file, err := tfutils.GetHCLFile(filepath.Join(basedir, constant.ProviderTf))
	if err != nil {
		return fmt.Errorf("failed to create terraform file for aws provider: %w", err)
	}
	defer file.Close()
	body := hclFile.Body()

	// add versions block
	g.addVersionsBlock(body)

	// add provider block
	if err := g.addProviderBlock(basedir, body); err != nil {
		return fmt.Errorf("failed to add provider block: %w", err)
	}

	// write contents to file
	file.Write(hclFile.Bytes())
	return nil
}
