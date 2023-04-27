package tfaws

import (
	"fmt"
	"os"
)

type ITfAws interface {
	CreateConfig(basedir, resources TFAWSResourceMap) error
}

type TfAws struct{}

func NewTfAws() *TfAws {
	return &TfAws{}
}

func (*TfAws) CreateConfig(basedir string, resources TFAWSResourceMap) error {
	if err := os.MkdirAll(basedir, 0744); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	for restype, resnames := range resources {
		generator, present := resourceGenerators[TFAWSResType(restype)]
		if !present {
			return fmt.Errorf("terraform aws resource type %q not supported", restype)
		}
		for _, name := range resnames {
			if err := generator.Generate(basedir, name); err != nil {
				return fmt.Errorf("failed to create terraform aws resource for type %q: %w", restype, err)
			}
		}
	}
	return nil
}
