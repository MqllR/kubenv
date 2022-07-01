package k8s_test

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/mqllr/kubenv/pkg/k8s"
)

var (
	testingConfig1 = `
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

	testingConfig2 = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: FAKEVALUE2
    server: https://fakeurl2.com
  name: fakecluster2
contexts:
- context:
    cluster: fakecluster2
    namespace: fakens2
    user: fakeuser2
  name: fakecontext2
kind: Config
preferences: {}
users:
- name: fakeuser2
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      args:
      - token
      - -i
      - fakecluster2
      command: aws-iam-authenticator
`

	testingMergedConfig = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: FAKEVALUE2
    server: https://fakeurl2.com
  name: fakecluster2
- cluster:
    certificate-authority-data: FAKEVALUE
    server: https://fakeurl.com
  name: fakecluster
contexts:
- context:
    cluster: fakecluster2
    namespace: fakens2
    user: fakeuser2
  name: fakecontext2
- context:
    cluster: fakecluster
    namespace: fakens
    user: fakeuser
  name: fakecontext
kind: Config
preferences: {}
users:
- name: fakeuser2
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      args:
      - token
      - -i
      - fakecluster2
      command: aws-iam-authenticator
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

	testingConfigStruct = &k8s.KubeConfig{
		APIVersion:  "v1",
		Kind:        "Config",
		Preferences: map[string]interface{}{},
		Clusters: []*k8s.ClusterWithName{
			{
				Name: "fakecluster",
				Cluster: &k8s.Cluster{
					CertificatAuthorityData: "FAKEVALUE",
					Server:                  "https://fakeurl.com",
				},
			},
		},
		Contexts: []*k8s.ContextWithName{
			{
				Name: "fakecontext",
				Context: &k8s.Context{
					Cluster:   "fakecluster",
					Namespace: "fakens",
					User:      "fakeuser",
				},
			},
		},
		Users: []*k8s.UserWithName{
			{
				Name: "fakeuser",
				User: &k8s.User{
					Exec: &k8s.Exec{
						APIVersion: "client.authentication.k8s.io/v1alpha1",
						Args:       []string{"token", "-i", "fakecluster"},
						Command:    "aws-iam-authenticator",
					},
				},
			},
		},
	}
)

func TestNewKubeConfig(t *testing.T) {
	kubeconfig := k8s.NewKubeConfig()

	if kubeconfig.APIVersion != "v1" {
		t.Errorf("Excepted value for ApiVersion %s but got %s", "v1", kubeconfig.APIVersion)
	}

	if kubeconfig.Kind != "Config" {
		t.Errorf("Excepted value for Kind %s but got %s", "Config", kubeconfig.Kind)
	}
}

func TestNewKubeConfigFromReader(t *testing.T) {
	kubeconfig, err := loadKubeConfig(testingConfig1)

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
	kubeconfig, err := loadKubeConfig(testingConfig1)
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
	kubeconfig, err := loadKubeConfig(testingConfig1)
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
	kubeconfig, err := loadKubeConfig(testingConfig1)
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

func TestGetKubeConfigByContextName(t *testing.T) {
	kubeconfig, err := loadKubeConfig(testingMergedConfig)
	if err != nil {
		t.Error(err)
	}

	for _, test := range testSuitesGetByName {
		k, err := kubeconfig.GetKubeConfigByContextName(test.context)
		if err != nil && !test.errExpected {
			t.Errorf("Got an error %s, but this wasn't expected", err.Error())
		}

		if test.errExpected {
			continue
		}

		if !reflect.DeepEqual(k, testingConfigStruct) {
			t.Error("Cluster returned doesn't match with the original kubeconfig")
		}
	}
}

func TestSaveFile(t *testing.T) {
	kubeconfig, err := loadKubeConfig(testingConfig1)
	if err != nil {
		t.Error(err)
	}

	buf := new(bytes.Buffer)

	err = kubeconfig.Save(buf)
	if err != nil {
		t.Errorf("Writing test file error: %s", err)
	}
}

func TestAppend(t *testing.T) {
	kubeconfig1, err := loadKubeConfig(testingConfig1)
	if err != nil {
		t.Error(err)
	}

	kubeconfig2, err := loadKubeConfig(testingConfig2)
	if err != nil {
		t.Error(err)
	}

	kubeconfigMerged, err := loadKubeConfig(testingMergedConfig)
	if err != nil {
		t.Error(err)
	}

	kubeconfig2.Append(kubeconfig1)
	if !reflect.DeepEqual(kubeconfig2, kubeconfigMerged) {
		t.Error("kubeconfig1 not equal the merged kubeconfig")
	}
}

func loadKubeConfig(kubeConfig string) (*k8s.KubeConfig, error) {
	r := strings.NewReader(kubeConfig)
	kubeconfig, err := k8s.NewKubeConfigFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("Error when trying to unmarsh the test config: %s", err)
	}

	return kubeconfig, nil
}
