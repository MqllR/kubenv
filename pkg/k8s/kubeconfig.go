package k8s

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Cluster ...
type ClusterWithName struct {
	Cluster *Cluster `yaml:"cluster"`
	Name    string   `yaml:"name"`
}

type Cluster struct {
	CertificatAuthorityData string `yaml:"certificate-authority-data"`
	Server                  string `yaml:"server"`
}

// Context ...
type ContextWithName struct {
	Context *Context `yaml:"context"`
	Name    string   `yaml:"name"`
}

type Context struct {
	Cluster   string `yaml:"cluster"`
	Namespace string `yaml:"namespace,omitempty"`
	User      string `yaml:"user"`
}

// User ...
type UserWithName struct {
	User *User  `yaml:"user"`
	Name string `yaml:"name"`
}

type User struct {
	ClientCertificateData string `yaml:"client-certificate-data,omitempty"`
	ClientKeyData         string `yaml:"client-key-data,omitempty"`
	Username              string `yaml:"username,omitempty"`
	Password              string `yaml:"password,omitempty"`
	Exec                  *Exec  `yaml:"exec,omitempty"`
	Token                 string `yaml:"token,omitempty"`
}

// Exec ...
type Exec struct {
	APIVersion string   `yaml:"apiVersion"`
	Args       []string `yaml:"args"`
	Command    string   `yaml:"command"`
	Env        []*Env   `yaml:"env"`
}

// Env ...
type Env struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// KubeConfig represent a kubernetes client configuration
type KubeConfig struct {
	APIVersion     string                 `yaml:"apiVersion"`
	Clusters       []*ClusterWithName     `yaml:"clusters"`
	Contexts       []*ContextWithName     `yaml:"contexts"`
	CurrentContext string                 `yaml:"current-context,omitempty"`
	Kind           string                 `yaml:"kind"`
	Preferences    map[string]interface{} `yaml:"preferences"`
	Users          []*UserWithName        `yaml:"users"`
}

// NewKubeConfig create a new struct KubeConfig
func NewKubeConfig() *KubeConfig {
	return &KubeConfig{
		APIVersion:  "v1", // Initiale values
		Preferences: map[string]interface{}{},
		Kind:        "Config",
	}
}

// NewKubeConfig create a new struct KubeConfig
func NewKubeConfigFromFile(kubeconfig string) (*KubeConfig, error) {
	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		return nil, fmt.Errorf("File doesn't exist: %s", err)
	}

	content, err := ioutil.ReadFile(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("Error when reading kubeconfig file: %s", err)
	}

	k := NewKubeConfig()

	if err = k.Unmarshal(content); err != nil {
		return nil, fmt.Errorf("Can't unmarshal the kubeconfig file: %s", err)
	}

	return k, nil
}

// Unmarshal fill a kubeConfig struct with yaml.Unmarshal
func (kubeConfig *KubeConfig) Unmarshal(config []byte) error {
	err := yaml.Unmarshal(config, kubeConfig)
	if err != nil {
		return err
	}

	return nil
}

// Marshal convert to []byte a KubeConfig
func (kubeConfig *KubeConfig) Marshal() ([]byte, error) {
	return yaml.Marshal(&kubeConfig)
}

// WriteFile Marshal KubeConfig in a file
func (kubeConfig *KubeConfig) WriteFile(file string) error {
	config, err := kubeConfig.Marshal()
	if err != nil {
		return nil
	}

	return ioutil.WriteFile(file, config, 0644)
}

// Append each Clusters, Users and Contexts into another KubeConfig
func (kubeConfig *KubeConfig) Append(config *KubeConfig) {
	kubeConfig.Clusters = append(kubeConfig.Clusters, config.Clusters...)
	kubeConfig.Contexts = append(kubeConfig.Contexts, config.Contexts...)
	kubeConfig.Users = append(kubeConfig.Users, config.Users...)
}

// GetContextNames returns all context names
func (kubeConfig *KubeConfig) GetContextNames() []string {
	var contexts []string
	for _, context := range kubeConfig.Contexts {
		contexts = append(contexts, context.Name)
	}

	return contexts
}

// ToString convert a KubeConfig in a string
func (kubeConfig *KubeConfig) ToString() (string, error) {
	config, err := kubeConfig.Marshal()
	return string(config), err
}

// ToString convert a KubeConfig in a string
func (kubeConfig *KubeConfig) IsContextExist(context string) bool {
	exist := func(slice []string, item string) bool {
		for _, s := range slice {
			if item == s {
				return true
			}
		}
		return false
	}

	if !exist(kubeConfig.GetContextNames(), context) {
		return false
	}

	return true
}
