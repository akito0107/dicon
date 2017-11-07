package sample

/*
func TestDependency(t *testing.T) {
	di := NewDIContainer()
	dep := di.Dependency()
	if e := dep.Run(); e != nil {
		t.Error(e)
	}
}

func TestSampleComponent(t *testing.T) {
	di := NewDIContainer()
	s := di.SampleComponent()
	if e := s.Exec(); e != nil {
		t.Error(e)
	}
	sc, ok := s.(*sampleComponent)
	if !ok {
		t.Error("sc is not a sampleComponent")
	}

	if sc.dep == nil {
		t.Error("dependency is nil")
	}

}

func TestOtherComponent(t *testing.T) {
	di := NewDIContainer()
	o := di.OtherComponent()
	oc, _ := o.(*otherComponent)

	if oc.s == nil {
		t.Error("samplecomponent is nil")
	}

	sc, ok := oc.s.(*sampleComponent)
	if !ok {
		t.Error("other.sample is invalid")
	}
	if sc.dep == nil {
		t.Error("other.sample.dep is nil")
	}
	_, ok = sc.dep.(*dependency)
	if !ok {
		t.Error("other.sample.dep is invalid")
	}
}
*/
