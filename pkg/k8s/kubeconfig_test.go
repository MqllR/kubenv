package k8s

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

var testingConfig = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: FAKEVALUE
    server: https://fakeurl.com
  name: fakecluster
contexts:
- context:
    cluster: fakecluster
    namespace: fakens
    user: fakeuser
  name: fakecontext
kind: Config
preferences: {}
users:
- name: fakeuser
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      args:
      - token
      - -i
      - fakecluster
      command: aws-iam-authenticator
`

func TestNewKubeConfig(t *testing.T) {
	kubeconfig := NewKubeConfig()

	if kubeconfig.APIVersion != "v1" {
		t.Errorf("Excepted value for ApiVersion %s but got %s", "v1", kubeconfig.APIVersion)
	}

	if kubeconfig.Kind != "Config" {
		t.Errorf("Excepted value for Kind %s but got %s", "Config", kubeconfig.Kind)
	}
}

func TestNewKubeConfigFromReader(t *testing.T) {
	r := strings.NewReader(testingConfig)
	kubeconfig, err := NewKubeConfigFromReader(r)

	if err != nil {
		t.Errorf("Unexpeted err: %s", err)
	}

	if kubeconfig.GetContextNames()[0] != "fakecontext" {
		t.Errorf("Got context %s but was expecting to get %s", kubeconfig.GetContextNames()[0], "fakecontext")
	}

	if kubeconfig.APIVersion != "v1" {
		t.Errorf("Excepted value for ApiVersion %s but got %s", "v1", kubeconfig.APIVersion)
	}

	if kubeconfig.Kind != "Config" {
		t.Errorf("Excepted value for Kind %s but got %s", "Config", kubeconfig.Kind)
	}
}

func TestUnmarshal(t *testing.T) {
	kubeconfig, err := loadKubeConfig()
	if err != nil {
		t.Error(err)
	}

	if len(kubeconfig.Clusters) != 1 {
		t.Errorf("Excepted array length of clusters %d but got %d", 1, len(kubeconfig.Clusters))
	}

	if len(kubeconfig.Users) != 1 {
		t.Errorf("Excepted array length of users %d but got %d", 1, len(kubeconfig.Users))
	}

	if len(kubeconfig.Contexts) != 1 {
		t.Errorf("Excepted array length of contexts %d but got %d", 1, len(kubeconfig.Contexts))
	}

	if kubeconfig.Contexts[0].Name != "fakecontext" {
		t.Errorf("Excepted context name %s but got %s", "fakecontext", kubeconfig.Contexts[0].Name)
	}

	if kubeconfig.Users[0].Name != "fakeuser" {
		t.Errorf("Excepted user name %s but got %s", "fakeuser", kubeconfig.Users[0].Name)
	}

	if kubeconfig.Clusters[0].Name != "fakecluster" {
		t.Errorf("Excepted cluster name %s but got %s", "fakecluster", kubeconfig.Clusters[0].Name)
	}
}

var (
	testSuitesGetByName = []struct {
		context     string
		errExpected bool
	}{
		{"foo", true},
		{"fakecontext", false},
	}
)

func TestGetContextByContextName(t *testing.T) {
	kubeconfig, err := loadKubeConfig()
	if err != nil {
		t.Error(err)
	}

	for _, test := range testSuitesGetByName {
		context, err := kubeconfig.GetContextByContextName(test.context)
		if err != nil && !test.errExpected {
			t.Errorf("Got an error %s, but this wasn't expected", err.Error())
		}

		if test.errExpected {
			continue
		}

		if !reflect.DeepEqual(context, kubeconfig.Contexts[0].Context) {
			t.Error("Context returned doesn't match with the original kubeconfig")
		}
	}
}

func TestGetUserByContextName(t *testing.T) {
	kubeconfig, err := loadKubeConfig()
	if err != nil {
		t.Error(err)
	}

	for _, test := range testSuitesGetByName {
		user, err := kubeconfig.GetUserByContextName(test.context)
		if err != nil && !test.errExpected {
			t.Errorf("Got an error %s, but this wasn't expected", err.Error())
		}

		if test.errExpected {
			continue
		}

		if !reflect.DeepEqual(user, kubeconfig.Users[0].User) {
			t.Error("User returned doesn't match with the original kubeconfig")
		}
	}
}

func TestGetClusterByContextName(t *testing.T) {
	kubeconfig, err := loadKubeConfig()
	if err != nil {
		t.Error(err)
	}

	for _, test := range testSuitesGetByName {
		cluster, err := kubeconfig.GetClusterByContextName(test.context)
		if err != nil && !test.errExpected {
			t.Errorf("Got an error %s, but this wasn't expected", err.Error())
		}

		if test.errExpected {
			continue
		}

		if !reflect.DeepEqual(cluster, kubeconfig.Clusters[0].Cluster) {
			t.Error("Cluster returned doesn't match with the original kubeconfig")
		}
	}
}

func TestWriteFile(t *testing.T) {
	kubeconfig, err := loadKubeConfig()
	if err != nil {
		t.Error(err)
	}

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	randString := func(n int) string {
		b := make([]byte, n)
		for i := range b {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		}
		return string(b)
	}

	filename := filepath.Join("/tmp", fmt.Sprintf("kubeconfig-test-%s", randString(5)))
	defer os.Remove(filename)

	err = kubeconfig.WriteFile(filename)
	if err != nil {
		t.Errorf("Writing test file error: %s", err)
	}

	info, err := os.Stat(filename)
	if err != nil {
		t.Errorf("Getting stat on test file error: %s", err)
	}

	if info.Mode() != 0600 {
		t.Errorf("Bad file mode: %v expected %v", info.Mode(), 0600)
	}
}

func loadKubeConfig() (*KubeConfig, error) {
	kubeconfig := NewKubeConfig()
	err := kubeconfig.Unmarshal([]byte(testingConfig))
	if err != nil {
		return nil, fmt.Errorf("Error when trying to unmarsh the test config: %s", err)
	}

	return kubeconfig, nil
}
