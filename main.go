package main

import (
	"fmt"
	"github.com/fabianflu/nxc/configreader"
	"github.com/fabianflu/nxc/dns"
	"github.com/fabianflu/nxc/wireguard"
)

func main() {

	config := configreader.ReadConfig("config.json")
	for _, Mode := range config.Mode {
		switch Mode {
		case "dns":
			dns.DefineZones()
			dns.DownloadZones()
			dns.CheckZones()
			dns.ApplyZones()
		case "wireguard":
			for _, WireguardConfig := range config.WireguardConfig {
				fmt.Printf("Networkname: %s Peer: %s\n", WireguardConfig.NetworkName, WireguardConfig.Peer)
				wireguard.DownloadConfig()
				wireguard.CheckConfig()
				wireguard.ApplyConfig()
			}
		}

	}

	fmt.Printf("baseUrl: %s, token: %s, Wireguard-networks: %s\n", config.BaseUrl, config.NXToken, config.WireguardConfig[0].NetworkName)
}
