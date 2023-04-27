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

const tfAwsS3 TFAWSResType = "s3"

func init() {
	addResourceType(tfAwsS3, NewTfAwsS3())
}

type TfAwsS3 struct{}

func NewTfAwsS3() *TfAwsS3 {
	return &TfAwsS3{}
}

func (g *TfAwsS3) getBucketID(bucketname string) hcl.Traversal {
	return hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_s3_bucket",
		},
		hcl.TraverseAttr{
			Name: bucketname,
		},
		hcl.TraverseAttr{
			Name: "id",
		},
	}
}

func (g *TfAwsS3) restrictPublicAccess(body *hclwrite.Body, bucketname string) {
	body.AppendNewline()
	resBlock := body.AppendNewBlock("resource", []string{"aws_s3_bucket_public_access_block", fmt.Sprintf("%s-default", bucketname)})
	resBlock.Body().SetAttributeTraversal("bucket", g.getBucketID(bucketname))
	resBlock.Body().SetAttributeValue("block_public_acls", cty.BoolVal(true))
	resBlock.Body().SetAttributeValue("block_public_policy", cty.BoolVal(true))
	resBlock.Body().SetAttributeValue("ignore_public_acls", cty.BoolVal(true))
	resBlock.Body().SetAttributeValue("restrict_public_buckets", cty.BoolVal(true))
}

func (g *TfAwsS3) enableEncryption(body *hclwrite.Body, bucketname string) {
	body.AppendNewline()
	resBlock := body.AppendNewBlock("resource", []string{"aws_s3_bucket_server_side_encryption_configuration", fmt.Sprintf("%s-default", bucketname)})
	resBlock.Body().SetAttributeTraversal("bucket", g.getBucketID(bucketname))

	ruleBlock := resBlock.Body().AppendNewBlock("rule", nil)
	encryptionBlock := ruleBlock.Body().AppendNewBlock("apply_server_side_encryption_by_default", nil)
	encryptionBlock.Body().SetAttributeValue("sse_algorithm", cty.StringVal("AES256"))
}

func (g *TfAwsS3) enableVersioning(body *hclwrite.Body, bucketname string) {
	body.AppendNewline()
	versioningBlock := body.AppendNewBlock("resource", []string{"aws_s3_bucket_versioning", fmt.Sprintf("%s-enabled", bucketname)})
	versioningBlock.Body().SetAttributeTraversal("bucket", g.getBucketID(bucketname))
	configBlock := versioningBlock.Body().AppendNewBlock("versioning_configuraion", nil)
	configBlock.Body().SetAttributeValue("status", cty.StringVal("Enabled"))
}

func (g *TfAwsS3) createAwsS3(body *hclwrite.Body, basedir, name string) error {
	body.AppendNewline()
	block := body.AppendNewBlock("resource", []string{"aws_s3_bucket", name})
	block.Body().SetAttributeTraversal("bucket", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "var",
		},
		hcl.TraverseAttr{
			Name: "bucket_name",
		},
	})
	return tfutils.AddVariable(basedir, "bucket_name", "AWS S3 bucket to store terraform state files")
}

func (g *TfAwsS3) Generate(basedir, bucketname string) error {
	// get main.tf hcl file
	hclFile, file, err := tfutils.GetHCLFile(filepath.Join(basedir, constant.MainTf))
	if err != nil {
		return fmt.Errorf("failed to get hcl file: %w", err)
	}
	defer file.Close()
	body := hclFile.Body()

	// create aws s3 bucket
	if err := g.createAwsS3(body, basedir, bucketname); err != nil {
		return fmt.Errorf("failed to create terraform configuration for aws s3 bucket: %w", err)
	}

	// enable versioning
	g.enableVersioning(body, bucketname)

	// enable server side encryption
	g.enableEncryption(body, bucketname)

	// block public access
	g.restrictPublicAccess(body, bucketname)

	if err := tfutils.AddOutput(basedir, "aws_s3_bucket_arn", "ARN of the S3 bucket which store terraform state", []string{"aws_s3_bucket", bucketname, "arn"}); err != nil {
		return fmt.Errorf("failed to add output variables for aws s3 bucket ARN: %w", err)
	}

	file.Write(hclwrite.Format(hclFile.Bytes()))
	return nil
}
