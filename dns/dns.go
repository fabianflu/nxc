package dns

import (
	"errors"
	"github.com/fabianflu/nxc/configreader"
	"github.com/fabianflu/nxc/filefetcher"
	"github.com/fabianflu/nxc/filehandler"
	"log"
	"os"
	"os/exec"
)

var logger = log.New(os.Stdout, "[dns] ", log.LstdFlags)

var clientConfig = configreader.NxcConfig{}

type LocalZone struct {
	ZoneName     string ""
	ZoneFileName string ""
	TempZonePath string ""
}

func ApplyDnsConfiguration(config configreader.NxcConfig) {
	clientConfig = config
	nxConfig := configreader.FetchNxConfigurationFromRemote(config)

	masterZone, e := defineZone(nxConfig)
	if e != nil {
		logger.Panicf("Failed to find matching zone for %v with rror: %v! Shutting down", config.DnsConfig.TargetServerName, e)
	}
	tempDirName := clientConfig.DnsConfig.LocalTempDir
	filehandler.CreateDirIfNotExist(tempDirName)
	for _, zoneName := range masterZone.Zones {
		localZone := LocalZone{
			ZoneName:     zoneName,
			ZoneFileName: zoneName + ".db",
			TempZonePath: tempDirName + "/" + zoneName + ".db",
		}

		err := localZone.downloadZone()
		if err != nil {
			logger.Printf("Failed to download zone: %v with error: %v! SKPPING", zoneName, err)
			continue
		}
		if !localZone.checkZone() {
			continue
		}
		localZone.applyZone()
	}
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

func (zone LocalZone) downloadZone() error {
	logger.Print("Downloading file:", zone.ZoneFileName)
	zoneUrl := clientConfig.BaseUrl + clientConfig.DnsConfig.RemoteZonePath + zone.ZoneFileName
	e := filefetcher.DownloadFileFromWeb(zone.TempZonePath, zoneUrl, true, clientConfig.NXToken)
	return e
}

func (zone LocalZone) checkZone() bool {
	out, runError := exec.Command("named-checkzone", zone.ZoneName, zone.TempZonePath).CombinedOutput()
	output := string(out)
	success := runError == nil
	if success {
		logger.Printf("Zone-check for zone %v and file %v successfull. %v", zone.ZoneName, zone.TempZonePath, output)
	} else {
		logger.Printf("Zone-check for zone %v and file %v failed. Command failed with output %v and error %v",
			zone.ZoneName, zone.TempZonePath, output, runError)
	}
	return success

}
func (zone LocalZone) applyZone() {
	logger.Printf("TODO: Applying zone %v by copying it to: %v",
		zone.ZoneName, clientConfig.DnsConfig.LocalZonePath+zone.ZoneFileName)
}
