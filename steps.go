package brewery

import (
	"context"
	"io"
	"io/ioutil"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var containerTimeout = 10 * time.Second

// Step defines an action to run in or with containers in a provided environment.
type Step interface {
	Run(context.Context, Environment) error
}

// StepFunc is a function that can act as a Step.
type StepFunc func(context.Context, Environment) error

// Run the step function
func (s StepFunc) Run(ctx context.Context, env Environment) error {
	return s(ctx, env)
}

// EnsureContainer is a step which ensures that the provided brews containers are running.
//
// Will create new containers if the provided brew has no created containers.
type EnsureContainer struct {
	Containers []*Brew
}

// Run will boot every brews containers
func (e *EnsureContainer) Run(ctx context.Context, env Environment) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, cn := range e.Containers {
		cn := cn
		g.Go(func() error {
			return e.run(ctx, env, cn)
		})
	}

	return g.Wait()
}

func (e *EnsureContainer) run(ctx context.Context, env Environment, brew *Brew) error {
	containerName := env.ContainerName(brew.Key)
	cs, err := env.Client.ContainerInspect(ctx, containerName)
	if err != nil && !client.IsErrContainerNotFound(err) {
		return err
	}

	switch {
	case err != nil:
		return e.create(ctx, env, brew)
	case cs.State.Running:
		return nil
	default:
		return e.start(ctx, env, cs.ID)
	}
}

func (e *EnsureContainer) build(ctx context.Context, env Environment, brew *Brew) error {
	return errors.Errorf("Currently not able to build from dockerfiles")
}

func (e *EnsureContainer) create(ctx context.Context, env Environment, brew *Brew) error {
	if brew.Build != "" {
		return e.build(ctx, env, brew)
	}

	rc, err := env.Client.ImagePull(ctx, brew.Image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	if _, err = io.Copy(ioutil.Discard, rc); err != nil {
		return err
	}
	if err = rc.Close(); err != nil {
		return err
	}

	containerName := env.ContainerName(brew.Key)
	config := &container.Config{
		AttachStdout: true,
		AttachStderr: true,
		Env:          brew.Envs,
		Image:        brew.Image,
	}
	hostConfig := &container.HostConfig{
		AutoRemove: true,
	}
	nwConfig := &network.NetworkingConfig{}

	cbody, err := env.Client.ContainerCreate(ctx, config, hostConfig, nwConfig, containerName)
	if err != nil {
		return err
	}

	return e.start(ctx, env, cbody.ID)
}

func (e *EnsureContainer) start(ctx context.Context, env Environment, cid string) error {
	return env.Client.ContainerRestart(ctx, cid, &containerTimeout)
}

// EnsureContainers creates a new EnsureContainer step with the provided containers.
func EnsureContainers(containers ...*Brew) *EnsureContainer {
	return &EnsureContainer{
		Containers: containers,
	}
}

type StepBundle struct {
	Steps []Step
}

func (s StepBundle) Run(ctx context.Context, env Environment) error {
	for _, step := range s.Steps {
		if err := step.Run(ctx, env); err != nil {
			return err
		}
	}
	return nil
}

func (s *StepBundle) add(step Step) {
	s.Steps = append(s.Steps, step)
}
