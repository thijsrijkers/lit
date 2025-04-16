package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"lit/namespacecgroup"
)

func RunContainer(config *ContainerConfig) error {
	// 1. Set up OverlayFS for the root filesystem
	rootfs, err := SetupOverlayFS(config.Image)
	if err != nil {
		return fmt.Errorf("overlayfs error: %v", err)
	}

	// 2. Apply cgroup settings via namespacecgroup package
	cgroupConfig := namespacecgroup.RuntimeConfig{
		CgroupMemoryLimit: config.MemoryLimit,
		CgroupCPULimit:    config.CPULimit,
		NamespaceType:     config.NamespaceType,
	}
	if err := namespacecgroup.SetupCgroup(cgroupConfig); err != nil {
		return fmt.Errorf("cgroup setup error: %v", err)
	}

	// 3. Set up Linux namespaces
	if err := syscall.Unshare(syscall.CLONE_NEWUTS |
		syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
		syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC); err != nil {
		return fmt.Errorf("namespace error: %v", err)
	}

	// 4. Chroot into container root filesystem
	if err := syscall.Chroot(rootfs); err != nil {
		return fmt.Errorf("chroot error: %v", err)
	}
	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("chdir error: %v", err)
	}

	// 5. Prepare command to run inside container
	cmd := exec.Command(config.Entrypoint, config.Args...)
	cmd.Env = os.Environ()
	for k, v := range config.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWIPC,
	}

	// 6. Start the container process
	return cmd.Run()
}
