package main

import (
	"fmt"
	"log"
	"os"

	"lit/container"
)

func main() {
	configPath := "lit.yml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	config, err := container.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Running container from image: %s\n", config.Image)
	if err := container.RunContainer(config); err != nil {
		log.Fatalf("Container failed: %v", err)
	}
}
