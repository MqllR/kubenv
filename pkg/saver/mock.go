package saver

type SaveMock struct{}

var _ Saver = &SaveMock{}

func NewSaveMock() *SaveMock { return &SaveMock{} }

func (s *SaveMock) SaveConfig(data []byte) error {
	return nil
}
