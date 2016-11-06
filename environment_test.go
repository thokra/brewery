package brewery

import (
	"testing"
)

func TestEnvironmentContainerName(t *testing.T) {
	env := Environment{
		Build:      3,
		Repository: "https://git.d09.no/lager/brewery.git",
	}

	brew := &Brew{
		Key: "lager",
	}

	expected := "lager-9c6c0e8202cb418d25bd7124d67f26ba-3"
	if name := env.ContainerName(brew.Key); name != expected {
		t.Errorf("Got %q, expected %q", name, expected)
	}
}
