package brewery

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Auth describes how to authenticate the docker registry
type Auth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Project contains project settings
type Project struct {
	Mount   string        `yaml:"mount"`   // Mountpoint for the repository
	Timeout time.Duration `yaml:"timeout"` // Timeout of the entire build, includes fetching images
}

func (p *Project) Brew() *Brew {
	return &Brew{
		Key:     "project",
		Name:    "Repository",
		Image:   "thokra/repo:latest",
		Workdir: p.Mount,
	}
}

// Brew describes how to test and build an application
type Brew struct {
	Key      string              `yaml:"-"`        // Key used for the brew. This is the YAML key
	Name     string              `yaml:"name"`     // Name of the build
	Image    string              `yaml:"image"`    // Docker image to build from
	Build    string              `yaml:"build"`    // Dockerfile to build
	Workdir  string              `yaml:"workdir"`  // Custom working directory
	Links    []string            `yaml:"link"`     // Linked containers. Array should contain keys for other brews
	Envs     []string            `yaml:"env"`      // Environment variables
	Commands []string            `yaml:"commands"` // Commands to run
	Auth     *Auth               `yaml:"auth"`     // Registry authentication
	Volumes  map[string][]string `yaml:"volumes"`  // Mount volumes from other containers

	dependencies map[string]struct{} // Array containing keys of brews this brew depend on
}

// Valid ensures that the brew contains all the required fields.
func (b *Brew) Valid() error {
	switch {
	case b.Name == "":
		return errors.Errorf("Invalid name")
	case b.Name == "project":
		return errors.Errorf("Name (%q) is reserved.", b.Name)
	case b.Image == "" && b.Build == "":
		return errors.Errorf("Image or build must be set")
	case b.Image != "" && b.Build != "":
		return errors.Errorf("Cannot set both image and build")
	}
	return nil
}

func (b *Brew) calculateDependencies() {
	b.dependencies = make(map[string]struct{})

	for _, link := range b.Links {
		b.dependencies[link] = struct{}{}
	}

	for k := range b.Volumes {
		b.dependencies[k] = struct{}{}
	}
}

func (b *Brew) validateDependencies(brews []string) error {
	b.calculateDependencies()

	dependencies := []string{}
	for d := range b.dependencies {
		dependencies = append(dependencies, d)
	}

	// Brews not found
	notFound := strSliceExclude(dependencies, brews...)
	if len(notFound) > 0 {
		return &ErrDependencyMissing{Missing: notFound}
	}

	return nil
}

// ErrDependencyMissing is used if there are missing dependencies in the recipe.
//
// Contains the names of the missing dependencies.
type ErrDependencyMissing struct {
	Missing []string
}

func (e *ErrDependencyMissing) Error() string {
	return fmt.Sprintf("Dependency not found: %v", e.Missing)
}
