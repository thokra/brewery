package brewery

import (
	"crypto/md5"
	"fmt"

	"github.com/docker/docker/client"
)

// Environment specifies the environment in which the build steps are run
type Environment struct {
	Client     *client.Client
	Build      int
	Repository string
}

// ContainerName generates a container name from the provided name, repository and build number
func (e Environment) ContainerName(name string) string {
	return fmt.Sprintf("%s-%x-%d", name, md5.Sum([]byte(e.Repository)), e.Build)
}
