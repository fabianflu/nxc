package main

import (
	"fmt"
	"github.com/fabianflu/nxc/configreader"
	"github.com/fabianflu/nxc/dns"
	"github.com/fabianflu/nxc/wireguard"
)

func main() {

	clientConfig := configreader.ReadConfig("config.json")
	configreader.FetchNxConfigurationFromRemote(clientConfig)
	for _, Mode := range clientConfig.Mode {
		switch Mode {
		case "dns":
			dns.ApplyDnsConfiguration(clientConfig)
		case "wireguard":
			for _, WireguardConfig := range clientConfig.WireguardConfig {
				fmt.Printf("Networkname: %s Peer: %s\n", WireguardConfig.NetworkName, WireguardConfig.Peer)
				wireguard.DownloadConfig()
				wireguard.CheckConfig()
				wireguard.ApplyConfig()
			}
		}

	}
}
