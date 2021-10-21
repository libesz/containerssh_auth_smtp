package config

import (
	"fmt"
	"io/fs"
	"io/ioutil"

	"github.com/containerssh/log"
	"gopkg.in/yaml.v2"
)

type mappingFileHandler struct {
	logger       log.Logger
	contentToken chan struct{}
	content      MappingFile
}

type token struct{}

func NewMappingFileHandler(logger log.Logger) *mappingFileHandler {
	contentToken := make(chan struct{}, 1)
	contentToken <- token{}
	return &mappingFileHandler{
		logger:       logger,
		contentToken: contentToken,
	}
}

func (file *mappingFileHandler) Render() []byte {
	token := <-file.contentToken
	defer func() { file.contentToken <- token }()
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
	token := <-file.contentToken
	defer func() { file.contentToken <- token }()
	file.content = content
	return err
}

func (file *mappingFileHandler) Set(content MappingFile) {
	token := <-file.contentToken
	defer func() { file.contentToken <- token }()
	file.content = content
}

func (file *mappingFileHandler) UserExist(userName string) bool {
	token := <-file.contentToken
	defer func() { file.contentToken <- token }()
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
	token := <-file.contentToken
	defer func() { file.contentToken <- token }()
	volumes := []string{}
	for _, volume := range file.content.Volumes {
		for _, user := range volume.Users {
			if user == userName {
				volumes = append(volumes, file.content.VolumePrefix+volume.VolumeName)
			}
		}
	}
	if len(volumes) == 0 {
		return volumes, fmt.Errorf("user not found")
	}
	return volumes, nil
}

func (file *mappingFileHandler) GetVolumePrefix() string {
	token := <-file.contentToken
	defer func() { file.contentToken <- token }()
	return file.content.VolumePrefix
}
