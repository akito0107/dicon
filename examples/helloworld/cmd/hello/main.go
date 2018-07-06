package main

import "github.com/akito0107/dicon/examples/helloworld/di"

func main() {
	container := di.NewDIContainer()
	s, _ := container.SampleComponent()
	s.Exec()
}
