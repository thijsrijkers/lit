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

func CreateOverlayFS(overlay OverlayFS) error {
	// Check if the base and writable layers exist
	if _, err := os.Stat(overlay.BaseLayer); os.IsNotExist(err) {
		return fmt.Errorf("base layer does not exist: %v", err)
	}
	if _, err := os.Stat(overlay.WritableLayer); os.IsNotExist(err) {
		return fmt.Errorf("writable layer does not exist: %v", err)
	}

	// Prepare OverlayFS mount
	err := syscall.Mount(overlay.BaseLayer, overlay.MountPoint, "overlay", syscall.MS_MGC_VAL, "lowerdir="+overlay.BaseLayer+",upperdir="+overlay.WritableLayer+",workdir="+overlay.MountPoint+"/work")
	if err != nil {
		return fmt.Errorf("failed to mount OverlayFS: %v", err)
	}

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
