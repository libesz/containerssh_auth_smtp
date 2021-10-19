package config

import (
	"testing"

	"github.com/containerssh/log"
)

func TestUserExist(t *testing.T) {
	file := MappingFile{VolumePrefix: "sites-", VolumePostfix: "-1", Users: []UserSpecs{
		{
			UserName:   "bela",
			VolumeName: "bela-com",
		},
	}}
	handler := mappingFileHandler{log.NewTestLogger(t), file}
	//handler.Dump()
	if !handler.UserExist("bela") {
		t.Fatal("User shall exist: bela")
	}
	if handler.UserExist("geza") {
		t.Fatal("User shall not exist: bela")
	}
}

func TestGetVolumeName(t *testing.T) {
	file := MappingFile{VolumePrefix: "sites-", VolumePostfix: "-1", Users: []UserSpecs{
		{
			UserName:   "bela",
			VolumeName: "bela-com",
		},
	}}
	handler := mappingFileHandler{log.NewTestLogger(t), file}
	volumeName, err := handler.GetUserVolumeName("bela")
	if err != nil {
		t.Fatal("User volume shall not error: bela")
	}
	if volumeName != "sites-bela-com-1" {
		t.Fatal("Volume name must be: sites-bela-com-1")
	}
}
