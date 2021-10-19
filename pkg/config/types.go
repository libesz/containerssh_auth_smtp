package config

type VolumeSpecs struct {
	VolumeName string   `yaml:volumename`
	Users      []string `yaml:users`
}

type MappingFile struct {
	VolumePrefix string        `yaml:"volumeprefix"`
	Volumes      []VolumeSpecs `yaml: "volumes"`
}

type MappingFileHandler interface {
	Render() []byte
	Dump()
	Save(path string) error
	Load(path string) error
	Set(content MappingFile)
	UserExist(userName string) bool
	GetUserVolumeNames(userName string) ([]string, error)
	GetVolumePrefix() string
}
