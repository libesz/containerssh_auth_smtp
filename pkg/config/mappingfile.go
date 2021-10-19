package config

import (
	"fmt"
	"io/fs"
	"io/ioutil"

	"github.com/containerssh/log"
	"gopkg.in/yaml.v2"
)

type mappingFileHandler struct {
	logger  log.Logger
	content MappingFile
}

func NewMappingFileHandler(logger log.Logger) *mappingFileHandler {
	return &mappingFileHandler{
		logger: logger,
	}
}

func (file *mappingFileHandler) Render() []byte {
	d, err := yaml.Marshal(&file.content)
	if err != nil {
		file.logger.Error("error: %v", err)
		return []byte{}
	}
	return d
}

func (file *mappingFileHandler) Dump() {
	file.logger.Info(string(file.Render()))
}

func (file *mappingFileHandler) Save(path string) error {
	err := ioutil.WriteFile(path, file.Render(), fs.ModePerm)
	return err
}

func (file *mappingFileHandler) Load(path string) error {
	rawContent, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	content := MappingFile{}
	err = yaml.Unmarshal(rawContent, &content)
	file.content = content
	return err
}

func (file *mappingFileHandler) Set(content MappingFile) {
	file.content = content
}

func (file *mappingFileHandler) UserExist(userName string) bool {
	for _, user := range file.content.Users {
		if user.UserName == userName {
			return true
		}
	}
	return false
}

func (file *mappingFileHandler) GetUserVolumeName(userName string) (string, error) {
	for _, user := range file.content.Users {
		if user.UserName == userName {
			return file.content.VolumePrefix + user.VolumeName + file.content.VolumePostfix, nil
		}
	}
	return "", fmt.Errorf("user not found")
}
