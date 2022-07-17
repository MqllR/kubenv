package saver

type Saver interface {
	SaveConfig([]byte) error
}
