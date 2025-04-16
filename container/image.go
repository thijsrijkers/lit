package container

import (
	"fmt"
	"os"
	"syscall"
)

func SetupOverlayFS(image string) (string, error) {
	base := fmt.Sprintf("/var/lib/lit/images/%s/layer", image)
	upper := fmt.Sprintf("/var/lib/lit/images/%s/writable", image)
	work := fmt.Sprintf("/var/lib/lit/images/%s/work", image)
	mountPoint := fmt.Sprintf("/var/lib/lit/containers/%s/rootfs", image)

	err := os.MkdirAll(mountPoint, 0755)
	if err != nil {
		return "", err
	}

	opts := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", base, upper, work)
	err = syscall.Mount("overlay", mountPoint, "overlay", 0, opts)
	if err != nil {
		return "", err
	}

	return mountPoint, nil
}
