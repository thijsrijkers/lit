package namespacecgroup

import (
	"fmt"
	"os"
	"syscall"
)

// Namespace and Cgroup struct
type RuntimeConfig struct {
	CgroupMemoryLimit int64
	CgroupCPULimit    int64
	NamespaceType     string
}

func CreateNamespace() error {
	// Create a new namespace (PID, Network, Mount, etc.)
	err := syscall.Unshare(syscall.CLONE_NEWNET | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID)
	if err != nil {
		return fmt.Errorf("failed to create namespace: %v", err)
	}
	return nil
}

func SetupCgroup(config RuntimeConfig) error {
	// Memory cgroup
	memoryPath := "/sys/fs/cgroup/memory/lit_container"
	if err := os.MkdirAll(memoryPath, 0755); err != nil {
		return fmt.Errorf("failed to create memory cgroup: %v", err)
	}
	if err := os.WriteFile(memoryPath+"/memory.limit_in_bytes", []byte(fmt.Sprintf("%d", config.CgroupMemoryLimit)), 0644); err != nil {
		return fmt.Errorf("failed to set memory limit: %v", err)
	}

	// CPU cgroup
	cpuPath := "/sys/fs/cgroup/cpu,cpuacct/lit_container"
	if err := os.MkdirAll(cpuPath, 0755); err != nil {
		return fmt.Errorf("failed to create CPU cgroup: %v", err)
	}
	if err := os.WriteFile(cpuPath+"/cpu.cfs_quota_us", []byte(fmt.Sprintf("%d", config.CgroupCPULimit)), 0644); err != nil {
		return fmt.Errorf("failed to set CPU limit: %v", err)
	}

	// Assign current process (or another PID) to both cgroups
	pid := fmt.Sprintf("%d", os.Getpid())

	if err := os.WriteFile(memoryPath+"/cgroup.procs", []byte(pid), 0644); err != nil {
		return fmt.Errorf("failed to assign process to memory cgroup: %v", err)
	}
	if err := os.WriteFile(cpuPath+"/cgroup.procs", []byte(pid), 0644); err != nil {
		return fmt.Errorf("failed to assign process to CPU cgroup: %v", err)
	}

	return nil
}
