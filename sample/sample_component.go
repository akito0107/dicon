package sample

type SampleComponent interface {
	Exec() error
}

type sampleComponent struct {
}

func NewSampleComponent() (SampleComponent, error) {
	return &sampleComponent{}, nil
}

func (s *sampleComponent) Exec() error {
	return nil
}
