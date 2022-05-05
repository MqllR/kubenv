package sync_test

const (
	kubeconfig1 = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: FAKEVALUE
    server: https://fakeurl.com
  name: fakecluster1
contexts:
- context:
    cluster: fakecluster1
    namespace: fakens1
    user: fakeuser1
  name: fakecontext1
kind: Config
preferences: {}
users:
- name: fakeuser1
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      args:
      - token
      - -i
      - fakecluster
      command: aws-iam-authenticator
`
	kubeconfig2 = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: FAKEVALUE
    server: https://fakeurl.com
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
      - fakecluster
      command: aws-iam-authenticator
`
)
