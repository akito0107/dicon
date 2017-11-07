package sample

import "github.com/akito0107/dicon/sample2"

// +DICON
type DIContainer interface {
	SampleComponent() SampleComponent
	OtherComponent() OtherComponent
	MoreComponent() MoreComponent
	Sample2Component() sample2.Sample2Component
}
