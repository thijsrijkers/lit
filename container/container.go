package container

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"lit/configlauncher"
	"lit/overlayfs"
)

type Container struct {
	ConfigPath string
	MountBase  string
}

func (c *Container) Run() error {
	// Step 1: Parse Config
	config, err := configlauncher.ParseConfig(c.ConfigPath)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	// Step 2: Apply namespace & cgroup config
	err = configlauncher.ApplyConfig(config)
	if err != nil {
		return fmt.Errorf("apply config: %w", err)
	}

	// Step 3: Setup OverlayFS
	fs := overlayfs.OverlayFS{
		BaseLayer:     c.MountBase + "/mnt/base",
		WritableLayer: c.MountBase + "/mnt/write",
		MountPoint:    c.MountBase + "/mnt/mount",
	}

	err = overlayfs.CreateOverlayFS(fs)
	if err != nil {
		return fmt.Errorf("overlayfs setup: %w", err)
	}

	defer func() {
		// Defer unmount with better error handling and logging.
		if unmountErr := overlayfs.UnmountOverlayFS(fs.MountPoint); unmountErr != nil {
			log.Printf("Warning: failed to unmount OverlayFS mount point at %s: %v", fs.MountPoint, unmountErr)
		} else {
			log.Printf("Successfully unmounted OverlayFS mount point at %s", fs.MountPoint)
		}
	}()

	// Step 4: Chroot into container and run the app
	if err := syscall.Chroot(fs.MountPoint); err != nil {
		return fmt.Errorf("failed to chroot into container at %s: %w", fs.MountPoint, err)
	}
	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("failed to chdir after chroot: %w", err)
	}

	// Step 5: Run the containerized application
	cmd := exec.Command("/bin/" + config.Image) // image = binary name from config
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Log the start of the application and its execution
	log.Println("Launching containerized app...")

	// Run the containerized app and check for any execution errors
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run containerized app: %w", err)
	}

	log.Println("Containerized app finished execution.")
	return nil
}
