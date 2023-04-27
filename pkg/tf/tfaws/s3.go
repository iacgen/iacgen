package tfaws

const tfAwsS3 TFAWSResType = "s3"

func init() {
	addResourceType(tfAwsS3, NewTfAwsS3())
}

type TfAwsS3 struct{}

func NewTfAwsS3() *TfAwsS3 {
	return &TfAwsS3{}
}

func (g *TfAwsS3) Generate(basedir, name string) error {
	return nil
}
