package config

import (
	"testing"

	"github.com/containerssh/log"
)

func TestUserExist(t *testing.T) {
	file := MappingFile{
		VolumePrefix: "sites_",
		Volumes: []VolumeSpecs{
			{
				VolumeName: "bela-com",
				Users: []string{
					"bela",
				},
			},
		},
	}
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
	file := MappingFile{
		VolumePrefix: "sites_",
		Volumes: []VolumeSpecs{
			{
				VolumeName: "bela-com",
				Users: []string{
					"bela",
				},
			},
		},
	}
	handler := mappingFileHandler{log.NewTestLogger(t), file}
	volumeNames, err := handler.GetUserVolumeNames("bela")
	if err != nil {
		t.Fatal("User volume shall not error: bela")
	}
	if len(volumeNames) != 1 {
		t.Fatal("Volume list length must be 1")
	}
	if volumeNames[0] != "sites_bela-com" {
		t.Fatal("Volume name must be: sites-bela-com-1")
	}
}
