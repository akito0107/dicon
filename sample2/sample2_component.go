package sample2

type Sample2Component interface {
	Exec() error
}

type sampleComponent struct {
}

func NewSample2Component() Sample2Component {
	return &sampleComponent {}
}

func (s *sampleComponent) Exec() error {
	return nil
}