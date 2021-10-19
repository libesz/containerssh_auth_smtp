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
	for _, volume := range file.content.Volumes {
		for _, user := range volume.Users {
			if user == userName {
				return true
			}
		}
	}
	return false
}

func (file *mappingFileHandler) GetUserVolumeNames(userName string) ([]string, error) {
	volumes := []string{}
	for _, volume := range file.content.Volumes {
		for _, user := range volume.Users {
			if user == userName {
				volumes = append(volumes, file.content.VolumePrefix+volume.VolumeName)
			}
		}
	}
	if len(volumes) == 0 {
		return volumes, fmt.Errorf("User not found")
	}
	return volumes, nil
}

func (file *mappingFileHandler) GetVolumePrefix() string {
	return file.content.VolumePrefix
}
