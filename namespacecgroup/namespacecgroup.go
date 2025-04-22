package namespacecgroup

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

// RuntimeConfig struct holds configuration options for namespaces and cgroups.
type RuntimeConfig struct {
	CgroupMemoryLimit int64
	CgroupCPULimit    int64
	NamespaceType     string
}

// CreateNamespace creates new namespaces (PID, Network, Mount).
func CreateNamespace() error {
	// Create a new namespace (PID, Network, Mount, etc.)
	err := syscall.Unshare(syscall.CLONE_NEWNET | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID)
	if err != nil {
		return fmt.Errorf("failed to create namespace: %v", err)
	}
	return nil
}

// SetupCgroup sets up the memory and CPU cgroups based on the provided config.
func SetupCgroup(config RuntimeConfig) error {
	// Ensure cgroups are mounted
	if err := mountCgroupSubsystems(); err != nil {
		return err
	}

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

// mountCgroupSubsystems checks if cgroup subsystems are mounted and mounts them if necessary.
func mountCgroupSubsystems() error {
	// Check if cgroup subsystems are mounted. If not, try mounting them.
	mounts := []string{
		"/sys/fs/cgroup/memory",
		"/sys/fs/cgroup/cpu,cpuacct",
	}

	for _, path := range mounts {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// Attempt to mount the cgroup subsystem if it does not exist.
			log.Printf("Cgroup path %s does not exist, attempting to mount it.", path)
			if err := syscall.Mount("cgroup", path, "cgroup", 0, ""); err != nil {
				return fmt.Errorf("failed to mount cgroup subsystem %v: %v", path, err)
			}
		} else {
			log.Printf("Cgroup path %s already exists and is mounted.", path)
		}
	}

	return nil
}
