package dockertest

import (
	"context"
	"errors"
	"log"
	"time"
)

var (
	ErrWaitedTooLong = errors.New("dockertest: waited too long for container from pool")
)

// A ContainerPool wraps multiple containers to allow t.Parallel to work well
type ContainerPool struct {
	containers chan *Container
}

// NewContainerPool is exactly the same as RunContainer, but takes a number of containers to boot
func NewContainerPool(num int, container string, port string, waitFunc func(addr string) error, args ...string) (*ContainerPool, error) {
	containerBuffer := make(chan *Container, num)
	for i := 0; i < num; i++ {
		c, err := RunContainer(container, port, waitFunc, args...)
		if err != nil {
			return nil, err
		}
		containerBuffer <- c
	}

	return &ContainerPool{
		containers: containerBuffer,
	}, nil
}

func (cp *ContainerPool) GetContainer(ctx context.Context) (*Container, error) {
	select {
	case <-ctx.Done():
		return nil, ErrWaitedTooLong
	case c := <-cp.containers:
		return c, nil
	}
}

func (cp *ContainerPool) ReleaseContainer(c *Container) {
	cp.containers <- c
}

func (cp *ContainerPool) Shutdown() {
	// ensure we have all containers back

	maxPoll := 30
	for {
		if maxPoll == 0 {
			log.Println("dockertest: cannot shutdown pool, not all containers returned")
		}

		if len(cp.containers) == cap(cp.containers) {
			break
		}

		time.Sleep(50 * time.Millisecond)
		maxPoll--
	}

	close(cp.containers)
	for c := range cp.containers {
		c.Shutdown()
	}
}
