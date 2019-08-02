package configreader

import (
	"encoding/json"
	"io/ioutil"
)

type DnsConfig struct {
	TargetServerName string               `json:"dns-server-name"`
	RemotePaths      DnsRemotePathsConfig `json:"remote"`
	LocalPaths       DnsLocalPathsConfig  `json:"local"`
}

type DnsLocalPathsConfig struct {
	LocalTempPath  string `json:"temp-path"`
	ZonePath       string `json:"zone-path"`
	BindConfigPath string `json:"config-path"`
}

type DnsRemotePathsConfig struct {
	ZonePath       string `json:"zone-path"`
	BindConfigPath string `json:"config-path"`
}

type WireguardConfig struct {
	NetworkName string `json:"network-name"`
	Peer        string `json:"peer"`
}
type NxcConfig struct {
	NXToken          string            `json:"x-nx-token"`
	BaseUrl          string            `json:"baseurl"`
	RemoteConfigPath string            `json:"remoteConfigPath"`
	Mode             []string          `json:"mode"`
	WireguardConfig  []WireguardConfig `json:"wireguard"`
	DnsConfig        DnsConfig         `json:"dns"`
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
