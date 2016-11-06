package brewery

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type (
	// Artifacts describes the artifacts to pull from containers for deployment
	Artifacts map[string][]string

	// EnvironmentVariables contains environment variables provided by the builder
	EnvironmentVariables map[string]string
)

// ErrMissingEnvVariables is used when environment variables is not found during
// recipe creation. Contains the environment variables not found.
type ErrMissingEnvVariables struct {
	Missing []string
}

func (e *ErrMissingEnvVariables) Error() string {
	return fmt.Sprintf("Undefined environment variables %v", e.Missing)
}

// Read in a brewfile and create a recipe
func Read(r io.Reader, env EnvironmentVariables) (*Recipe, error) {
	rec, err := read(r, env)
	if err != nil {
		return nil, err
	}

	if err := rec.Valid(); err != nil {
		return nil, err
	}

	return rec, nil
}

func read(r io.Reader, env EnvironmentVariables) (*Recipe, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	b, err = envReplacer(b, env)
	if err != nil {
		return nil, err
	}

	rec := &Recipe{
		Brews: make(map[string]*Brew),
	}

	// TODO(thokra): Create custom yaml parser so that we can add distributors without
	// chainging this package. Might also remove the usage of regex in envReplacer
	if err = yaml.Unmarshal(b, rec); err != nil {
		return nil, err
	}

	return rec, nil
}

var brewEnvMatcher = regexp.MustCompile(`\$\$([\w\d]+)`)

// envReplacer matches environment variables defined by double `$` and replaces them with
// environment variables provided by the caller.
//
// Returns an error of type ErrMissingEnvVariables if one or more environment variables is missing
func envReplacer(b []byte, env EnvironmentVariables) ([]byte, error) {
	missing := []string{}
	b = brewEnvMatcher.ReplaceAllFunc(b, func(in []byte) []byte {
		e := strings.TrimLeft(string(in), "$")
		val, ok := env[e]
		if !ok {
			missing = append(missing, e)
		}

		return []byte(val)
	})

	if len(missing) > 0 {
		return b, &ErrMissingEnvVariables{Missing: missing}
	}

	return b, nil
}
