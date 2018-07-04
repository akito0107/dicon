package dicon

import (
	"strings"
)

type CyclicDependencyError struct {
	cyclicPath []string
}

func (e *CyclicDependencyError) Error() string {
	return "detect cyclic dependency '" + strings.Join(e.cyclicPath, "' -> '") + "'"
}

type walkState struct {
	stack  []string
	marked map[string]struct{}
}

func (ws *walkState) visit(name string) error {
	ws.stack = append(ws.stack, name)
	if _, ok := ws.marked[name]; ok {
		for i, x := range ws.stack {
			if x == name {
				return &CyclicDependencyError{cyclicPath: ws.stack[i:]}
			}
		}
		panic("unreachable")
	}
	ws.marked[name] = struct{}{}
	return nil
}

func (ws *walkState) leave() {
	last := ws.stack[len(ws.stack)-1]
	ws.stack = ws.stack[:len(ws.stack)-1]
	delete(ws.marked, last)
}

type cyclicDetector struct {
	dependencies map[string][]string
	visited      map[string]struct{}
}

func (cd *cyclicDetector) detect() error {
	for name, _ := range cd.dependencies {
		if _, ok := cd.visited[name]; ok {
			continue
		}
		ws := &walkState{
			marked: make(map[string]struct{}),
		}
		if err := cd.walk(name, ws); err != nil {
			return err
		}
	}
	return nil
}

func (cd *cyclicDetector) walk(name string, state *walkState) error {
	if err := state.visit(name); err != nil {
		return err
	}
	if _, ok := cd.visited[name]; ok {
		return nil
	}
	for _, dependency := range cd.dependencies[name] {
		if err := cd.walk(dependency, state); err != nil {
			return err
		}
	}
	state.leave()
	cd.visited[name] = struct{}{}
	return nil
}

func DetectCyclicDependency(funcs []FuncType) error {
	dependencies := make(map[string][]string, len(funcs))
	for _, fn := range funcs {
		name := fn.ReturnTypes[0].SimpleName()

		deps := make([]string, 0, len(fn.ArgumentTypes))
		for _, dep := range fn.ArgumentTypes {
			deps = append(deps, dep.SimpleName())
		}
		dependencies[name] = deps
	}
	cd := &cyclicDetector{
		dependencies: dependencies,
		visited:      make(map[string]struct{}),
	}
	return cd.detect()
}
