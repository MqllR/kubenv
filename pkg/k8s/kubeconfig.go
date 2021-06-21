package k8s

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// KubeConfig represents a kubernetes client configuration
type KubeConfig struct {
	APIVersion     string                 `yaml:"apiVersion"`
	Clusters       []*ClusterWithName     `yaml:"clusters"`
	Contexts       []*ContextWithName     `yaml:"contexts"`
	CurrentContext string                 `yaml:"current-context,omitempty"`
	Kind           string                 `yaml:"kind"`
	Preferences    map[string]interface{} `yaml:"preferences"`
	Users          []*UserWithName        `yaml:"users"`
}

// NewKubeConfig creates a new struct KubeConfig
func NewKubeConfig() *KubeConfig {
	return &KubeConfig{
		APIVersion:  "v1", // Initiale values
		Preferences: map[string]interface{}{},
		Kind:        "Config",
	}
}

// NewKubeConfigFromFile creates a new struct KubeConfig from a file
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

// Unmarshal fills a kubeConfig struct with yaml.Unmarshal
func (kubeConfig *KubeConfig) Unmarshal(config []byte) error {
	err := yaml.Unmarshal(config, kubeConfig)
	if err != nil {
		return err
	}

	return nil
}

// Marshal converts to []byte a KubeConfig
func (kubeConfig *KubeConfig) Marshal() ([]byte, error) {
	return yaml.Marshal(&kubeConfig)
}

// WriteFile writes the kubeconfig in the given file
func (kubeConfig *KubeConfig) WriteFile(file string) error {
	config, err := kubeConfig.Marshal()
	if err != nil {
		return nil
	}

	return ioutil.WriteFile(file, config, 0644)
}

// WriteTempFile writes the kubeconfig in a temporary file
// returns the temporary file path
func (kubeConfig *KubeConfig) WriteTempFile() (string, error) {
	temp, err := ioutil.TempFile("/tmp", "kubeconfig-*")
	if err != nil {
		return "", fmt.Errorf("Cannot create a temporary file %s", err)
	}

	tempKubeConfig := temp.Name()
	defer temp.Close()

	data, err := kubeConfig.Marshal()
	if err != nil {
		return "", fmt.Errorf("Unable to marshal kubeconfig: %s", err)
	}

	_, err = temp.Write(data)
	if err != nil {
		return "", fmt.Errorf("Error when writting the temporary kubeconfig: %s", err)
	}

	return tempKubeConfig, nil
}

// Append merges 2 KubeConfig struct in one.
func (kubeConfig *KubeConfig) Append(config *KubeConfig) {
	kubeConfig.Clusters = append(kubeConfig.Clusters, config.Clusters...)
	kubeConfig.Contexts = append(kubeConfig.Contexts, config.Contexts...)
	kubeConfig.Users = append(kubeConfig.Users, config.Users...)
}

// GetContextNames returns a list of all the context names
func (kubeConfig *KubeConfig) GetContextNames() []string {
	var contexts []string
	for _, context := range kubeConfig.Contexts {
		contexts = append(contexts, context.Name)
	}

	return contexts
}

// ToString converts a KubeConfig in a string
func (kubeConfig *KubeConfig) ToString() (string, error) {
	config, err := kubeConfig.Marshal()
	return string(config), err
}

// IsContextExist checks if a given context exist in the KubeConfig
func (kubeConfig *KubeConfig) IsContextExist(context string) bool {
	exist := func(slice []string, item string) bool {
		for _, s := range slice {
			if item == s {
				return true
			}
		}
		return false
	}

	return exist(kubeConfig.GetContextNames(), context)
}

// SetCurrentContext just set the given context to CurrentContext
func (kubeConfig *KubeConfig) SetCurrentContext(context string) error {
	if !kubeConfig.IsContextExist(context) {
		return fmt.Errorf("Context %s doesn't exist in kubeconfig file", context)
	}

	kubeConfig.CurrentContext = context
	return nil
}
