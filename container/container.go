package container

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
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
		BaseLayer:     "/mnt/base",
		WritableLayer: "/mnt/write",
		MountPoint:    "/mnt/mount",
	}
	

	// Function to handle OverlayFS setup with cleanup and retries
	setupOverlayFS := func() error {
		// Cleanup any previous OverlayFS mounts if they exist
		if _, err := os.Stat(fs.MountPoint); err == nil {
			// Unmount previous OverlayFS if it's still mounted
			if unmountErr := overlayfs.UnmountOverlayFS(fs.MountPoint); unmountErr != nil {
				log.Printf("Warning: failed to unmount previous OverlayFS mount point at %s: %v", fs.MountPoint, unmountErr)
			} else {
				log.Printf("Successfully unmounted previous OverlayFS mount point at %s", fs.MountPoint)
			}
		}

		// Retry the OverlayFS setup a few times if needed
		for i := 0; i < 3; i++ {
			err := overlayfs.CreateOverlayFS(fs)
			if err == nil {
				log.Println("Successfully mounted OverlayFS")
				return nil
			}

			// Log the error and attempt again
			log.Printf("Failed to mount OverlayFS (attempt %d): %v", i+1, err)
			time.Sleep(2 * time.Second) // Sleep before retrying
		}

		// If all attempts fail, return an error
		return fmt.Errorf("failed to mount OverlayFS after multiple attempts")
	}

	// Attempt to set up OverlayFS
	err = setupOverlayFS()
	if err != nil {
		return fmt.Errorf("overlayfs setup: %w", err)
	}

	defer func() {
		// Cleanup unmounting of OverlayFS when done
		if _, err := os.Stat(fs.MountPoint); err == nil {
			if unmountErr := overlayfs.UnmountOverlayFS(fs.MountPoint); unmountErr != nil {
				log.Printf("Warning: failed to unmount OverlayFS mount point at %s: %v", fs.MountPoint, unmountErr)
			} else {
				log.Printf("Successfully unmounted OverlayFS mount point at %s", fs.MountPoint)
			}
		} else {
			log.Printf("Mount point %s does not exist, skipping unmount.", fs.MountPoint)
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
