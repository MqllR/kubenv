package saver

type GeneratorMock struct{}

func (g *GeneratorMock) GenerateHistoryFilename() string {
	return "/home/foo/.kube/kubeconfig-1657725814415"
}
