package sample

type MoreComponent interface {
	Exec() error
}

type moreComponent struct {
	s SampleComponent
	o OtherComponent
}

func NewMoreComponent(s SampleComponent, o OtherComponent) MoreComponent {
	return &moreComponent{
		s: s,
		o: o,
	}
}

func (m *moreComponent) Exec() error {
	return nil
}
