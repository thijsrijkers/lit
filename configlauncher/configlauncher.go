package configlauncher

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"io/ioutil"
	"lit/namespacecgroup"
)

// Config struct for parsing lit.yml
type Config struct {
	NamespaceType  string `yaml:"namespace_type"`
	MemoryLimit    int64  `yaml:"memory_limit"`
	CPULimit       int64  `yaml:"cpu_limit"`
	Image          string `yaml:"image"`
}

func ParseConfig(filePath string) (*Config, error) {
	// Read the YAML config file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse the YAML data into the Config struct
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

func ApplyConfig(config *Config) error {
	// Apply configurations, for example, create namespaces and cgroups
	namespaceConfig := namespacecgroup.RuntimeConfig{
		CgroupMemoryLimit: config.MemoryLimit,
		CgroupCPULimit:    config.CPULimit,
		NamespaceType:     config.NamespaceType,
	}

	// Apply namespace and cgroup setup
	err := namespacecgroup.CreateNamespace()
	if err != nil {
		return fmt.Errorf("failed to create namespace: %v", err)
	}

	err = namespacecgroup.SetupCgroup(namespaceConfig)
	if err != nil {
		return fmt.Errorf("failed to set up cgroup: %v", err)
	}

	// More configuration can be applied here, e.g., pulling and launching the container image
	return nil
}
