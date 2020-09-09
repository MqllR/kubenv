package aws

import (
	"os"

	"gopkg.in/ini.v1"
)

type IniFile struct {
	file string
	cfg  *ini.File
}

func NewIniFile(file string) (*IniFile, error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, err
	}

	cfg, err := ini.Load(file)
	if err != nil {
		return nil, err
	}

	return &IniFile{file, cfg}, nil
}

func (i *IniFile) EnsureSectionAndSave(section string, params map[string]string) error {

	s, _ := i.cfg.GetSection(section)

	if s == nil {
		var err error
		s, err = i.cfg.NewSection(section)
		if err != nil {
			return err
		}
	}

	for key, value := range params {
		s.Key(key).SetValue(value)
	}

	err := i.cfg.SaveTo(i.file)

	if err != nil {
		return err
	}

	return nil
}
