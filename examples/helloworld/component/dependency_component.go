package component

import "log"

type DependencyComponent interface {
	Exec() error
}

type dependencyComponent struct {
}

func NewDependencyComponent() (DependencyComponent, error) {
	return &dependencyComponent{}, nil
}

func (s *dependencyComponent) Exec() error {
	log.Println("Hello World from DependencyComponent!")
	return nil
}
