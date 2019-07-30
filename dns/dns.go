package dns

import (
	"errors"
	"github.com/fabianflu/nxc/configreader"
	"github.com/fabianflu/nxc/filefetcher"
	"github.com/fabianflu/nxc/filehandler"
	"log"
	"os"
	"os/exec"
	"strings"
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
	serviceReloadRequired := false
	tempDirName := clientConfig.DnsConfig.LocalTempDir
	filehandler.CreateDirIfNotExist(tempDirName)
	serviceReloadRequired = updateZones(masterZone, tempDirName)
	if serviceReloadRequired {
		e := exec.Command("systemctl", "reload", "bind9").Run()
		if e != nil {
			logger.Panic("Failed to reload service:", e)
		}
		logger.Printf("Done! Finished configuring %v! Zone update and bind reload successful!", clientConfig.DnsConfig.TargetServerName)
	} else {
		logger.Printf("DNS configration COMPLETED")
	}

}

func applyNameServerConfig(nameServer string) bool {
	//TODO check if name-server-configuration is up to date and act accordingly
	return false
}

func updateZones(masterZone configreader.DnsMasterZone, tempDirName string) bool {
	serviceReloadRequired := false
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
			logger.Printf("Error: Zone %v is not valid! SKIPPING", zoneName)
			continue
		}
		serviceReloadRequired = localZone.refreshZoneIfNeeded() || serviceReloadRequired
	}
	return serviceReloadRequired
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
func (zone LocalZone) refreshZoneIfNeeded() bool {
	logger.Printf("Applying zone %v by copying it to: %v",
		zone.ZoneName, clientConfig.DnsConfig.LocalZonePath+zone.ZoneFileName)

	localZonePath := clientConfig.DnsConfig.LocalZonePath
	zoneFilePathBuilder := strings.Builder{}
	zoneFilePathBuilder.WriteString(localZonePath)
	if !strings.HasSuffix(localZonePath, "/") {
		zoneFilePathBuilder.WriteString("/")
	}
	zoneFilePathBuilder.WriteString(zone.ZoneFileName)
	targetZonePath := zoneFilePathBuilder.String()
	if filehandler.IsFileExistent(targetZonePath) {
		logger.Printf("Zone %v already exists comparing files to see if update is needed")
		if filehandler.AreFilesEqualByHash(zone.TempZonePath, targetZonePath) {
			logger.Printf("Your zone file: %v is up to date, no modification needed")
			return false
		}
	}
	copyErr := filehandler.CopyOrOverwrite(zone.TempZonePath, targetZonePath)
	if copyErr != nil {
		logger.Print("Failed to update", targetZonePath)
		return false
	}
	logger.Print("Updated:", targetZonePath)
	return true
}
