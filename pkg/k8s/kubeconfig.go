package k8s

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

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

// NewKubeConfigFromReader creates a new struct KubeConfig from an io.Reader
func NewKubeConfigFromReader(r io.Reader) (*KubeConfig, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("Error when reading kubeconfig reader: %s", err)
	}

	k := NewKubeConfig()

	if err = k.unmarshal(content); err != nil {
		return nil, fmt.Errorf("Can't unmarshal the kubeconfig file: %s", err)
	}

	return k, nil
}

// Save writes the kubeconfig in the given file
func (kubeConfig *KubeConfig) Save(w io.Writer) error {
	config, err := kubeConfig.marshal()
	if err != nil {
		return nil
	}

	_, err = w.Write(config)
	return err
}

// WriteTempFile writes the kubeconfig in a temporary file
// returns the temporary file path
func (kubeConfig *KubeConfig) WriteTempFile() (string, error) {
	temp, err := os.CreateTemp("/tmp", "kubeconfig-*")
	if err != nil {
		return "", fmt.Errorf("Cannot create a temporary file %s", err)
	}

	tempKubeConfig := temp.Name()
	defer temp.Close()

	data, err := kubeConfig.marshal()
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
	config, err := kubeConfig.marshal()
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

// GetContextByContextName returns a Context from its name
func (kubeConfig *KubeConfig) GetContextByContextName(context string) (*Context, error) {
	for _, ctx := range kubeConfig.Contexts {
		if ctx.Name == context {
			return ctx.Context, nil
		}
	}

	return nil, fmt.Errorf("Context %s not found in the context", context)
}

// GetClusterByContextName returns a Cluster from the context name
func (kubeConfig *KubeConfig) GetClusterByContextName(context string) (*Cluster, error) {
	ctx, err := kubeConfig.GetContextByContextName(context)
	if err != nil {
		return nil, err
	}

	for _, cluster := range kubeConfig.Clusters {
		if cluster.Name == ctx.Cluster {
			return cluster.Cluster, nil
		}
	}

	return nil, fmt.Errorf("Cluster not found not found for the context %s", context)
}

// GetUserByContextName returns a Users from the context name
func (kubeConfig *KubeConfig) GetUserByContextName(context string) (*User, error) {
	ctx, err := kubeConfig.GetContextByContextName(context)
	if err != nil {
		return nil, err
	}

	for _, user := range kubeConfig.Users {
		if user.Name == ctx.User {
			return user.User, nil
		}
	}

	return nil, fmt.Errorf("User not found for the context %s", context)
}

// GetKubeConfigByContextName returns a KubeConfig
func (kubeConfig *KubeConfig) GetKubeConfigByContextName(context string) (*KubeConfig, error) {
	k := &KubeConfig{}

	ctx, err := kubeConfig.GetContextByContextName(context)
	if err != nil {
		return nil, err
	}

	k.Contexts = []*ContextWithName{
		{Context: ctx, Name: context},
	}

	user, err := kubeConfig.GetUserByContextName(context)
	if err != nil {
		return nil, err
	}

	k.Users = []*UserWithName{
		{User: user, Name: context},
	}

	cluster, err := kubeConfig.GetClusterByContextName(context)
	if err != nil {
		return nil, err
	}

	k.Clusters = []*ClusterWithName{
		{Cluster: cluster, Name: context},
	}

	return k, nil
}

// ExecCommand executes any kind of command for a selected context
// this will write a temporary kubeconfig file in /tmp.
func (kubeConfig *KubeConfig) ExecCommand(context string, cmd []string) error {
	err := kubeConfig.SetCurrentContext(context)
	if err != nil {
		return fmt.Errorf("Error when settings the context: %s", err)
	}

	tempKubeConfig, err := kubeConfig.WriteTempFile()
	if err != nil {
		return fmt.Errorf("Error when creating the temporary kubeconfig file: %s", err)
	}
	defer os.Remove(tempKubeConfig)

	exe, err := exec.LookPath(cmd[0])
	if err != nil {
		return fmt.Errorf("Binary not found: %s", err)
	}

	envs := os.Environ()
	isExist := func(envs []string, key string) (bool, int) {
		for i, env := range envs {
			if env == key {
				return true, i
			}
		}

		return false, 0
	}

	exist, i := isExist(envs, "KUBECONFIG")
	localKubeConfig := "KUBECONFIG=" + tempKubeConfig
	if exist {
		envs[i] = localKubeConfig
	} else {
		envs = append(envs, localKubeConfig)
	}

	c := exec.Cmd{
		Path:   exe,
		Args:   cmd[0:],
		Env:    envs,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err = c.Run()

	if err != nil {
		return err
	}

	return nil
}

// unmarshal fills a kubeConfig struct with yaml.Unmarshal
func (kubeConfig *KubeConfig) unmarshal(config []byte) error {
	err := yaml.Unmarshal(config, kubeConfig)
	if err != nil {
		return err
	}

	return nil
}

// marshal converts to []byte a KubeConfig
func (kubeConfig *KubeConfig) marshal() ([]byte, error) {
	return yaml.Marshal(&kubeConfig)
}
