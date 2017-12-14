package sample

type MoreComponent interface {
	Exec() error
	ExecFun(a int, b string) func() SampleComponent
}

type moreComponent struct {
	s SampleComponent
	o OtherComponent
}

func NewMoreComponent(s SampleComponent, o OtherComponent) (MoreComponent, error) {
	return &moreComponent{
		s: s,
		o: o,
	}, nil
}

func (m *moreComponent) Exec() error {
	return nil
}

func (m *moreComponent) ExecFun(a int, b string) func() SampleComponent {
	return func() SampleComponent {
		return nil
	}
}
