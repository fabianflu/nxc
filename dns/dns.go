package dns

import (
	"errors"
	"fmt"
	"github.com/fabianflu/nxc/configreader"
	"github.com/fabianflu/nxc/filefetcher"
	"github.com/fabianflu/nxc/filehandler"
	"os/exec"
)

var clientConfig = configreader.NxcConfig{}

func ApplyDnsConfiguration(config configreader.NxcConfig) {
	clientConfig = config
	nxConfig := configreader.FetchNxConfigurationFromRemote(config)
	masterZone, e := defineZone(nxConfig)
	if e != nil {
		panic(e)
	}
	downloadZones(masterZone)
	checkZones(masterZone)
}

func defineZone(nxConfig configreader.NxConfig) (configreader.DnsMasterZone, error) {
	dnsZones := nxConfig.Namespace.DnsZones.MasterZones
	name := clientConfig.DnsConfig.TargetServerName
	for _, singleZone := range dnsZones {
		if singleZone.Name == name {
			return singleZone, nil
		}
	}
	return configreader.DnsMasterZone{}, errors.New("No matching zone found with name: " + name)
}

func downloadZones(zone configreader.DnsMasterZone) {
	tempDirName := clientConfig.DnsConfig.LocalTempDir
	filehandler.CreateDirIfNotExist(tempDirName)
	for _, zoneFileName := range zone.Zones {
		zoneFileName += ".db"
		tempZonePath := tempDirName + "/" + zoneFileName
		zoneUrl := clientConfig.BaseUrl + clientConfig.DnsConfig.RemoteZonePath + zoneFileName
		e := filefetcher.DownloadFileFromWeb(tempZonePath, zoneUrl, true, clientConfig.NXToken)
		if e != nil {
			panic(e)
		}
	}
}

func checkZones(zone configreader.DnsMasterZone) {
	tempDirName := clientConfig.DnsConfig.LocalTempDir
	for _, zoneName := range zone.Zones {
		zoneFileName := zoneName + ".db"
		tempZonePath := tempDirName + "/" + zoneFileName
		command := exec.Command("named-checkzone", zoneName, tempZonePath)
		fmt.Print(command)
	}

}
func applyZones() {

}
