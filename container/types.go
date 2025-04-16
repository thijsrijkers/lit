package container

type ContainerConfig struct {
	Image         string            `yaml:"image"`
	Entrypoint    string            `yaml:"entrypoint"`
	Args          []string          `yaml:"args"`
	NamespaceType string            `yaml:"namespace_type"`
	MemoryLimit   int64             `yaml:"memory_limit"`
	CPULimit      int64             `yaml:"cpu_limit"`
	Env           map[string]string `yaml:"env"`
	Mounts        []string          `yaml:"mounts"`
	Network       string            `yaml:"network"`
}
