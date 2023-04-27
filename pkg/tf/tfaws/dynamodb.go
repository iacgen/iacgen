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

const tfAwsDynamodb TFAWSResType = "dynamodb"

func init() {
	addResourceType(tfAwsDynamodb, NewTfAwsDynamodb())
}

type TfAwsDynamodb struct{}

func NewTfAwsDynamodb() *TfAwsDynamodb {
	return &TfAwsDynamodb{}
}

func (g *TfAwsDynamodb) createDynamodb(body *hclwrite.Body, basedir, tablename string) error {
	body.AppendNewline()

	resBlock := body.AppendNewBlock("resource", []string{"aws_dynamodb_table", tablename})
	resBlock.Body().SetAttributeTraversal("name", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "var",
		},
		hcl.TraverseAttr{
			Name: "table_name",
		},
	})
	resBlock.Body().SetAttributeValue("billing_mode", cty.StringVal("PAY_PER_REQUEST"))
	resBlock.Body().SetAttributeValue("hash_key", cty.StringVal("LockID"))

	attributeBlock := resBlock.Body().AppendNewBlock("attribute", nil)
	attributeBlock.Body().SetAttributeValue("name", cty.StringVal("LockID"))
	attributeBlock.Body().SetAttributeValue("type", cty.StringVal("S"))

	return tfutils.AddVariable(basedir, "table_name", "Dynamodb table name")
}

func (g *TfAwsDynamodb) Generate(basedir, name string) error {
	// create new empty hcl file object
	hclFile, file, err := tfutils.GetHCLFile(filepath.Join(basedir, constant.MainTf))
	if err != nil {
		return fmt.Errorf("failed to create terraform file for aws provider: %w", err)
	}
	defer file.Close()
	body := hclFile.Body()

	// create dynamodb table
	if err := g.createDynamodb(body, basedir, name); err != nil {
		return fmt.Errorf("failed to create aws dynamodb table: %w", err)
	}

	if err := tfutils.AddOutput(basedir, "dynamodb_table_name", "Name of the dynamodb table", []string{"aws_dynamodb_table", name, "name"}); err != nil {
		return fmt.Errorf("failed to add output variables for aws dynamodb table name: %w", err)
	}

	// write contents to file
	file.Write(hclFile.Bytes())
	return nil
}
