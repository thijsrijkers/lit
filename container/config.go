package container

import (
	"os"
	"gopkg.in/yaml.v2"
)

func LoadConfig(path string) (*ContainerConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config ContainerConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
