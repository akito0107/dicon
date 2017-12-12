package sample

type OtherComponent interface {
	Exec() error
}

type otherComponent struct {
	s SampleComponent
}

func NewOtherComponent(s SampleComponent) (OtherComponent, error) {
	return &otherComponent{
		s: s,
	}, nil
}

func (s *otherComponent) Exec() error {
	return nil
}

// interface継承をつかった名前解決
// 上位のinterfaceを明示的に使っていたらそっちを入れる
