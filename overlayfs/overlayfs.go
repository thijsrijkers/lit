package overlayfs

import (
	"fmt"
	"os"
	"syscall"
	"log"
)

type OverlayFS struct {
	BaseLayer    string
	WritableLayer string
	MountPoint   string
}

// CreateOverlayFS sets up the overlay filesystem
func CreateOverlayFS(fs OverlayFS) error {

	log.Printf("BaseLayer: %s", fs.BaseLayer)
	log.Printf("WritableLayer: %s", fs.WritableLayer)
	log.Printf("MountPoint: %s", fs.MountPoint)
	// Ensure directories exist (create them if they don't)
	if err := os.MkdirAll(fs.BaseLayer, 0755); err != nil {
		return fmt.Errorf("failed to create base layer directory: %w", err)
	}
	if err := os.MkdirAll(fs.WritableLayer, 0755); err != nil {
		return fmt.Errorf("failed to create writable layer directory: %w", err)
	}
	if err := os.MkdirAll(fs.MountPoint, 0755); err != nil {
		return fmt.Errorf("failed to create mount point directory: %w", err)
	}

	// Ensure the work directory exists and is empty
	workDir := fs.MountPoint + "/work"
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return fmt.Errorf("failed to create work directory: %w", err)
	}

	// Clean the work directory (if not empty)
	if err := removeContents(workDir); err != nil {
		return fmt.Errorf("failed to clean work directory: %w", err)
	}

	// Prepare the OverlayFS mount command
	mountCmd := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", fs.BaseLayer, fs.WritableLayer, workDir)

	// Perform the mount
	err := syscall.Mount("overlay", fs.MountPoint, "overlay", syscall.MS_MGC_VAL, mountCmd)
	if err != nil {
		return fmt.Errorf("failed to mount OverlayFS: %w", err)
	}

	return nil
}

// removeContents removes all files and directories inside a given directory
func removeContents(dir string) error {
	// Read the contents of the directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// Iterate over all entries and remove them
	for _, entry := range entries {
		path := dir + "/" + entry.Name()
		if entry.IsDir() {
			// Recursively remove directories
			if err := os.RemoveAll(path); err != nil {
				return fmt.Errorf("failed to remove directory %s: %w", path, err)
			}
		} else {
			// Remove regular files
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove file %s: %w", path, err)
			}
		}
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
