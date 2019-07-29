package configreader

import (
	"encoding/json"
	"github.com/fabianflu/nxc/filefetcher"
)

var fetchedConfig = NxConfig{}

type Netbox struct {
	Url    string `json:"url"`
	ApiKey string `json:"api_key"`
}

type DnsMasterZone struct {
	Name       string   `json:"name"`
	IpAddress  string   `json:"ip"`
	DottedMail string   `json:"dotted_mail"`
	Zones      []string `json:"zones"`
}

type Namespace struct {
	DnsZones DnsZones `json:"dns"`
}

type DnsZones struct {
	MasterZones []DnsMasterZone `json:"masters"`
}

type NxConfig struct {
	Netbox    Netbox    `json:"netbox"`
	Namespace Namespace `json:"namespaces"`
}

func ReadNxConfig(content []byte) NxConfig {
	err := json.Unmarshal(content, &fetchedConfig)
	if err != nil {
		panic(err)
	}
	return fetchedConfig
}

func FetchNxConfigurationFromRemote(clientConfig NxcConfig) NxConfig {
	if fetchedConfig.Netbox.Url == "" {
		bytes, e := filefetcher.FetchFile(clientConfig.RemoteConfigPath, clientConfig.NXToken, clientConfig.BaseUrl)
		if e != nil {
			panic(e)
		}
		ReadNxConfig(bytes)
	}
	return fetchedConfig
}
