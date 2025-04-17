package container

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"lit/configlauncher"
	"lit/overlayfs"
	"lit/namespacecgroup"
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
		BaseLayer:     c.MountBase + "/base",
		WritableLayer: c.MountBase + "/write",
		MountPoint:    c.MountBase + "/mount",
	}
	err = overlayfs.CreateOverlayFS(fs)
	if err != nil {
		return fmt.Errorf("overlayfs: %w", err)
	}
	defer func() {
		if unmountErr := overlayfs.UnmountOverlayFS(fs.MountPoint); unmountErr != nil {
			log.Printf("Warning: failed to unmount: %v", unmountErr)
		}
	}()

	// Step 4: chroot into container and run the app
	if err := syscall.Chroot(fs.MountPoint); err != nil {
		return fmt.Errorf("chroot: %w", err)
	}
	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("chdir: %w", err)
	}

	// Step 5: Run the container
	cmd := exec.Command("/bin/" + config.Image) // image = binary name
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println("Launching containerized app...")
	return cmd.Run()
}

