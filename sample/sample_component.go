package sample

type SampleComponent interface {
	Exec() error
}

type sampleComponent struct {
}

func NewSampleComponent() SampleComponent {
	return &sampleComponent{}
}

func (s *sampleComponent) Exec() error {
	return nil
}
