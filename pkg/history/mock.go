package history

type GeneratorMock struct{}

var _ IGenerator = &GeneratorMock{}

func (g *GeneratorMock) TimestampedFile() string {
	return "/home/foo/.kube/kubeconfig-1657725814415"
}
