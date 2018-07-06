package di

import (
	"github.com/akito0107/dicon/examples/helloworld/component"
)

// +DICON
type DIContainer interface {
	SampleComponent() (component.SampleComponent, error)
	DependencyComponent() (component.DependencyComponent, error)
}
