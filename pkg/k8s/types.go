package k8s

// ClusterWithName represents a cluster with its name
type ClusterWithName struct {
	Cluster *Cluster `yaml:"cluster"`
	Name    string   `yaml:"name"`
}

// Cluster represents a cluster in a kubeconfig file
type Cluster struct {
	CertificatAuthorityData string `yaml:"certificate-authority-data"`
	Server                  string `yaml:"server"`
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
	User *User  `yaml:"user"`
	Name string `yaml:"name"`
}

// User represents a user in a kubeconfig file
type User struct {
	ClientCertificateData string `yaml:"client-certificate-data,omitempty"`
	ClientKeyData         string `yaml:"client-key-data,omitempty"`
	Username              string `yaml:"username,omitempty"`
	Password              string `yaml:"password,omitempty"`
	Exec                  *Exec  `yaml:"exec,omitempty"`
	Token                 string `yaml:"token,omitempty"`
}

// Exec represents a exec section witin a user
type Exec struct {
	APIVersion string   `yaml:"apiVersion"`
	Args       []string `yaml:"args"`
	Command    string   `yaml:"command"`
	Env        []*Env   `yaml:"env"`
}

// Env represents the environment variable injected
// during the exec command
type Env struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
