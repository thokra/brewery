package brewery

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/thokra/brewery/distributors"
)

// Recipe contains instructions for all builds and artifacts
type Recipe struct {
	Project   Project          `yaml:"project"`
	Brews     map[string]*Brew `yaml:"brews"`
	Artifacts Artifacts        `yaml:"artifacts"`
	Publish   Publish          `yaml:"publish"`

	dependencyGraph graph
}

// Valid validates the recipe to find errors. Will also update some required fields.
func (r *Recipe) Valid() error {
	invalid := []string{}
	brewNames := []string{"project"}

	for b := range r.Brews {
		brewNames = append(brewNames, b)
	}

	r.Brews["project"] = r.Project.Brew()

	dg := graph{}
	dg.add("project")
	for k, brew := range r.Brews {
		brew.Key = k

		if brew.Name == "" {
			brew.Name = k
		}

		if err := brew.Valid(); err != nil {
			invalid = append(invalid, fmt.Sprintf("%v: %v", k, err))
		}

		if err := brew.validateDependencies(brewNames); err != nil {
			invalid = append(invalid, fmt.Sprintf("%v: %v", k, err))
		}

		deps := []string{}
		for dep := range brew.dependencies {
			deps = append(deps, dep)
		}
		dg.add(brew.Key, deps...)
	}

	// Calculate dependencygraph, check for circular dependencies
	var err error
	r.dependencyGraph, err = dg.resolve()
	if err != nil {
		circular := []string{}
		for _, n := range r.dependencyGraph {
			circular = append(circular, fmt.Sprintf("%v and %v", n.name, strings.Join(n.deps, ", ")))
		}

		invalid = append(invalid, "Circular dependency found. Might be between:", strings.Join(circular, "\n- "))
	}

	if len(invalid) > 0 {
		return errors.Errorf("One or more brew is invalid:\n-%v", strings.Join(invalid, "\n- "))
	}

	return nil
}

// Ingredients calculates the order brews should be run
func (r *Recipe) Ingredients() (*Ingredientslist, error) {
	if r.dependencyGraph == nil {
		if err := r.Valid(); err != nil {
			return nil, err
		}
	}

	added := make(map[string]struct{})
	root := &Ingredientslist{}
	current := root
	for _, node := range r.dependencyGraph {
		makeNew := false
		for _, d := range node.deps {
			for _, added := range current.Brews {
				if added == d {
					makeNew = true
				}
			}
		}

		if makeNew {
			for _, n := range current.Brews {
				added[n] = struct{}{}
			}

			nw := &Ingredientslist{}
			current.Next = nw
			current = nw
		}

		current.Brews = append(current.Brews, node.name)
	}

	return root, nil
}

func (r *Recipe) Steps() ([]Step, error) {
	il, err := r.Ingredients()
	if err != nil {
		return nil, err
	}

	steps := []Step{}
	for list := il; list != nil; list = list.Next {
		sb := &StepBundle{}
		brews := []*Brew{}
		for _, b := range list.Brews {
			brews = append(brews, r.Brews[b])
		}
		sb.add(EnsureContainers(brews...))
		steps = append(steps, sb)
	}

	return steps, nil
}

// Publish contains targets for deployment
type Publish struct {
	GCS distributors.GCS `yaml:"gcs"`
}
