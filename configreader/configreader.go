package configreader

import (
	"encoding/json"
	"io/ioutil"
)

type WireguardConfig struct {
	NetworkName string `json:"network-name"`
	Peer        string `json:"peer"`
}
type NxcConfig struct {
	NXToken         string            `json:"x-nx-token"`
	BaseUrl         string            `json:"baseurl"`
	Mode            []string          `json:"mode"`
	WireguardConfig []WireguardConfig `json:"wireguard"`
}

func ReadConfig(path string) NxcConfig {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	config := NxcConfig{}
	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		panic(err)
	}

	return config
}
