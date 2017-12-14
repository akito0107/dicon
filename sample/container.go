package sample

import "github.com/akito0107/dicon/sample2"

// +DICON
type DIContainer interface {
	SampleComponent() (SampleComponent, error)
	OtherComponent() (OtherComponent, error)
	MoreComponent() (MoreComponent, error)
	Sample2Component() (sample2.Sample2Component, error)
}
