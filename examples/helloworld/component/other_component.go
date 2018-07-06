package component

type SampleComponent interface {
	Exec() error
}

type sampleComponent struct {
	d DependencyComponent
}

func NewSampleComponent(d DependencyComponent) (SampleComponent, error) {
	return &sampleComponent{
		d: d,
	}, nil
}

func (s *sampleComponent) Exec() error {
	// using dependency
	if err := s.d.Exec(); err != nil {
		return err
	}

	return nil
}
