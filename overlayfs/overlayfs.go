package overlayfs

import (
	"fmt"
	"os"
	"syscall"
)

type OverlayFS struct {
	BaseLayer    string
	WritableLayer string
	MountPoint   string
}

func CreateOverlayFS(fs OverlayFS) error {
	// Ensure the directories exist for BaseLayer, WritableLayer, and MountPoint
	if err := os.MkdirAll(fs.BaseLayer, 0755); err != nil {
		return fmt.Errorf("failed to create base layer directory: %w", err)
	}
	if err := os.MkdirAll(fs.WritableLayer, 0755); err != nil {
		return fmt.Errorf("failed to create writable layer directory: %w", err)
	}
	if err := os.MkdirAll(fs.MountPoint, 0755); err != nil {
		return fmt.Errorf("failed to create mount point directory: %w", err)
	}

	// Create the work directory inside the mount point
	workDir := fs.MountPoint + "/work"
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return fmt.Errorf("failed to create work directory: %w", err)
	}

	// Prepare the OverlayFS mount command
	mountCmd := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", fs.BaseLayer, fs.WritableLayer, workDir)

	// Ensure the mount point is isolated before mounting
	err := syscall.Mount("overlay", fs.MountPoint, "overlay", syscall.MS_MGC_VAL, mountCmd)
	if err != nil {
		return fmt.Errorf("failed to mount OverlayFS: %w", err)
	}

	// Mount the work directory (if necessary) and return
	return nil
}

func UnmountOverlayFS(mountPoint string) error {
	// Unmount the OverlayFS mount point
	err := syscall.Unmount(mountPoint, 0)
	if err != nil {
		return fmt.Errorf("failed to unmount OverlayFS: %v", err)
	}
	return nil
}
