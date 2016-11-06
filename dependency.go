package brewery

import (
	"errors"
)

type node struct {
	name string
	deps []string
}

type graph []*node

var (
	errCircularDependency = errors.New("Circular dependency found")
)

func (g *graph) add(name string, deps ...string) {
	(*g) = append((*g), &node{name, deps})
}

func (g graph) resolve() (graph, error) {
	nodeNames := make(map[string]*node)
	nodeDependencies := make(map[string][]string)

	for _, n := range g {
		nodeNames[n.name] = n
		nodeDependencies[n.name] = n.deps
	}

	resolved := graph{}
	for len(nodeDependencies) > 0 {
		ready := []string{}
		for name, deps := range nodeDependencies {
			if len(deps) == 0 {
				ready = append(ready, name)
			}
		}

		if len(ready) == 0 {
			ng := graph{}
			for name, deps := range nodeDependencies {
				ng.add(name, deps...)
			}
			return ng, errCircularDependency
		}

		for _, name := range ready {
			delete(nodeDependencies, name)
			resolved = append(resolved, nodeNames[name])
		}

		for name, deps := range nodeDependencies {
			nodeDependencies[name] = strSliceExclude(deps, ready...)
		}
	}

	return resolved, nil
}

func strSliceExclude(s []string, exclude ...string) []string {
	n := []string{}
OUTER:
	for _, v := range s {
		for _, ex := range exclude {
			if v == ex {
				continue OUTER
			}
		}
		n = append(n, v)
	}
	return n
}
