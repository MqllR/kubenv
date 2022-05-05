package k8s

// ClusterWithName represents a cluster with its name
type ClusterWithName struct {
	Cluster interface{} `yaml:"cluster"`
	Name    string      `yaml:"name"`
}

// ContextWithName represents a context with its name
type ContextWithName struct {
	Context *Context `yaml:"context"`
	Name    string   `yaml:"name"`
}

// Context represents a context in a kubeconfig file
// It's a mapping between a cluster name and a user name.
// Optionally a namespace can be added
type Context struct {
	Cluster   string `yaml:"cluster"`
	Namespace string `yaml:"namespace,omitempty"`
	User      string `yaml:"user"`
}

// UserWithName represents a user with its name
type UserWithName struct {
	User interface{} `yaml:"user"`
	Name string      `yaml:"name"`
}
