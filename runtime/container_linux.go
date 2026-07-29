package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/docker/containerd/specs"
	"github.com/docker/containerd/subreaper/exec"
	"github.com/opencontainers/runc/libcontainer"
	ocs "github.com/opencontainers/runtime-spec/specs-go"
)

func (c *container) getLibctContainer() (libcontainer.Container, error) {
	runtimeRoot := "/run/runc"

	// Check that the root wasn't changed
	for _, opt := range c.runtimeArgs {
		if strings.HasPrefix(opt, "--root=") {
			runtimeRoot = strings.TrimPrefix(opt, "--root=")
			break
		}
	}

	f, err := libcontainer.New(runtimeRoot)
	if err != nil {
		return nil, err
	}
	return f.Load(c.id)
}

func u64Ptr(i uint64) *uint64 { return &i }

func (c *container) UpdateResources(r *Resource) error {
	sr := ocs.Resources{
		Memory: &ocs.Memory{
			Limit:       u64Ptr(uint64(r.Memory)),
			Reservation: u64Ptr(uint64(r.MemoryReservation)),
			Swap:        u64Ptr(uint64(r.MemorySwap)),
			Kernel:      u64Ptr(uint64(r.KernelMemory)),
		},
		CPU: &ocs.CPU{
			Shares: u64Ptr(uint64(r.CPUShares)),
			Quota:  u64Ptr(uint64(r.CPUQuota)),
			Period: u64Ptr(uint64(r.CPUPeriod)),
			Cpus:   &r.CpusetCpus,
			Mems:   &r.CpusetMems,
		},
		BlockIO: &ocs.BlockIO{
			Weight: &r.BlkioWeight,
		},
	}

	srStr := bytes.NewBuffer(nil)
	if err := json.NewEncoder(srStr).Encode(&sr); err != nil {
		return err
	}

	args := c.runtimeArgs
	args = append(args, "update", "-r", "-", c.id)
	cmd := exec.Command(c.runtime, args...)
	cmd.Stdin = srStr
	b, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(b))
	}
	return nil
}

func getRootIDs(s *specs.Spec) (int, int, error) {
	if s == nil {
		return 0, 0, nil
	}
	var hasUserns bool
	for _, ns := range s.Linux.Namespaces {
		if ns.Type == ocs.UserNamespace {
			hasUserns = true
			break
		}
	}
	if !hasUserns {
		return 0, 0, nil
	}
	uid := hostIDFromMap(0, s.Linux.UIDMappings)
	gid := hostIDFromMap(0, s.Linux.GIDMappings)
	return uid, gid, nil
}
