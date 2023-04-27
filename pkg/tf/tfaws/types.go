package tfaws

type (
	TFAWSResType     string
	TFAWSResName     string
	TFAWSResourceMap map[string][]string
)

type IResourceGenerator interface {
	Generate(basedir, name string) error
}
type ResourceGenerators map[TFAWSResType]IResourceGenerator

var resourceGenerators = make(ResourceGenerators)

func (m ResourceGenerators) addResourceType(resType TFAWSResType, generator IResourceGenerator) {
	m[resType] = generator
}

func addResourceType(resType TFAWSResType, generator IResourceGenerator) {
	resourceGenerators.addResourceType(resType, generator)
}
