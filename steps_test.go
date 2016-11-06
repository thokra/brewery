package brewery

import (
	"context"
	"testing"

	"github.com/docker/docker/client"
)

func TestEnsureContainer(t *testing.T) {
	client, err := client.NewEnvClient()
	if err != nil {
		t.Fatal(err)
	}

	brew := &Brew{
		Key:   "rethinkdb",
		Image: "rethinkdb:latest",
	}

	env := Environment{
		Client:     client,
		Build:      1,
		Repository: "https://git.d09.no/lager/brewery.git",
	}

	ec := EnsureContainers(brew)
	if err := ec.Run(context.TODO(), env); err != nil {
		t.Error(err)
	}
}
