package config

type UserSpecs struct {
	UserName   string `yaml:"username"`
	VolumeName string `yaml:"volumename"`
}

type MappingFile struct {
	VolumePrefix  string      `yaml:"volumeprefix"`
	VolumePostfix string      `yaml:"volumepostfix"`
	Users         []UserSpecs `yaml:"users"`
}

type MappingFileHandler interface {
	Render() []byte
	Dump()
	Save(path string) error
	Load(path string) error
	Set(content MappingFile)
	UserExist(userName string) bool
	GetUserVolumeName(userName string) (string, error)
}
