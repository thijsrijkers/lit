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
	// Set up the cgroup for memory and CPU limits
	// Cgroup file paths can differ based on your system (this is an example)
	cgroupPath := "/sys/fs/cgroup/memory/lit_container"

	err := os.MkdirAll(cgroupPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create cgroup: %v", err)
	}

	// Apply memory limit
	err = os.WriteFile(cgroupPath+"/memory.limit_in_bytes", []byte(fmt.Sprintf("%d", config.CgroupMemoryLimit)), 0644)
	if err != nil {
		return fmt.Errorf("failed to set memory limit: %v", err)
	}

	// Apply CPU limit
	err = os.WriteFile(cgroupPath+"/cpu.cfs_quota_us", []byte(fmt.Sprintf("%d", config.CgroupCPULimit)), 0644)
	if err != nil {
		return fmt.Errorf("failed to set CPU limit: %v", err)
	}

	return nil
}
