package k8s_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/mqllr/kubenv/pkg/k8s"
	"github.com/mqllr/kubenv/pkg/saver"
)

var (
	testingConfig1 = `
`

	testingConfig2 = `
`

	testingBadConfig1 = `
`

	testingMergedConfig = `
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
	kubeconfig, err := loadKubeConfig(t, "kubeconfig1")

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
	kubeconfig, err := loadKubeConfig(t, "kubeconfig1")
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
	kubeconfig, err := loadKubeConfig(t, "kubeconfig1")
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
	kubeconfig, err := loadKubeConfig(t, "kubeconfig1")
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
	kubeconfig, err := loadKubeConfig(t, "kubeconfig_merge")
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
	testSuites := []struct {
		kubeconfig  string
		errExpected bool
	}{
		{"kubeconfig1", false},
		{"kubeconfig_bad", true},
	}

	for _, test := range testSuites {
		kubeconfig, err := loadKubeConfig(t, test.kubeconfig)
		if err != nil {
			t.Error(err)
		}

		err = kubeconfig.Save(saver.NewSaveMock())
		if err != nil && !test.errExpected {
			t.Errorf("Got an error %s, but this wasn't expected", err.Error())
		}
	}
}

func TestAppend(t *testing.T) {
	kubeconfig1, err := loadKubeConfig(t, "kubeconfig1")
	if err != nil {
		t.Error(err)
	}

	kubeconfig2, err := loadKubeConfig(t, "kubeconfig2")
	if err != nil {
		t.Error(err)
	}

	kubeconfigMerged, err := loadKubeConfig(t, "kubeconfig_merge")
	if err != nil {
		t.Error(err)
	}

	kubeconfig2.Append(kubeconfig1)
	if !reflect.DeepEqual(kubeconfig2, kubeconfigMerged) {
		t.Error("kubeconfig1 not equal the merged kubeconfig")
	}
}

func loadKubeConfig(t *testing.T, kubeConfig string) (*k8s.KubeConfig, error) {
	fh, err := os.Open(fmt.Sprintf("testdata/%s.yaml", kubeConfig))
	if err != nil {
		t.Fatal(err)
	}

	kubeconfig, err := k8s.NewKubeConfigFromReader(fh)
	if err != nil {
		return nil, fmt.Errorf("Error when trying to unmarsh the test config: %s", err)
	}

	return kubeconfig, nil
}
