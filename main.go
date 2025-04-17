package main

import (
	"log"
	"lit/container"
)

func main() {
	runner := container.Container{
		ConfigPath: "lit.yml",
		MountBase:  ".", // Or absolute path to your rootfs layout
	}
	if err := runner.Run(); err != nil {
		log.Fatalf("Container run failed: %v", err)
	}
}

