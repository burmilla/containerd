package supervisor

import (
	"errors"

	"github.com/docker/containerd/runtime"
)

func NewMonitor() (*Monitor, error) {
	return &Monitor{}, errors.New("Monitor NewMonitor() not implemented on Solaris")
}

type Monitor struct {
	ooms chan string
}

func (m *Monitor) Exits() chan runtime.Process {
	return nil
}

func (m *Monitor) Monitor(p runtime.Process) error {
	return errors.New("Monitor Monitor() not implemented on Solaris")
}

func (m *Monitor) Close() error {
	return errors.New("Monitor Close() not implemented on Solaris")
}

func (m *Monitor) start() {
}
