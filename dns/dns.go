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

func ApplyDnsConfiguration(config configreader.NxcConfig) {
	clientConfig = config
	nxConfig := configreader.FetchNxConfigurationFromRemote(config)

	masterZone, e := defineZone(nxConfig)
	if e != nil {
		logger.Panicf("Failed to find matching zone for %v with rror: %v! Shutting down", config.DnsConfig.TargetServerName, e)
	}
	serviceReloadRequired := applyNameServerConfig(clientConfig.DnsConfig.TargetServerName)
	tempDirName := clientConfig.DnsConfig.LocalPaths.LocalTempPath
	filehandler.CreateDirIfNotExist(tempDirName)
	serviceReloadRequired = serviceReloadRequired || updateZones(masterZone, tempDirName)
	if serviceReloadRequired {
		e := exec.Command("systemctl", "reload", "bind9").Run()
		if e != nil {
			logger.Panic("Failed to reload service:", e)
		}
		logger.Printf("Done! Finished configuring %v! Zone update and bind reload successful!", clientConfig.DnsConfig.TargetServerName)
	} else {
		logger.Printf("Done! Dns configration already up to date!")
	}

}

func applyNameServerConfig(nameServer string) bool {
	confFileName := nameServer + ".conf"
	remoteConfFilePath := filehandler.BuildFilePathFromParts(clientConfig.BaseUrl, clientConfig.DnsConfig.RemotePaths.BindConfigPath, confFileName)
	tempConfFilePath := filehandler.BuildFilePathFromParts(clientConfig.DnsConfig.LocalPaths.LocalTempPath, confFileName)
	err := filefetcher.DownloadFileFromWeb(tempConfFilePath, remoteConfFilePath, true, clientConfig.NXToken)
	if err != nil {
		logger.Printf("Error: Failed to download nameserver conf file, SKIPPING configuration")
		return false
	}
	appliedConfPath := filehandler.BuildFilePathFromParts(clientConfig.DnsConfig.LocalPaths.BindConfigPath, confFileName)
	areFilesEqual := filehandler.AreFilesEqualByHash(tempConfFilePath, appliedConfPath)
	if areFilesEqual {
		logger.Printf("Conf file is already up to date")
		return false
	}
	output, commandError := exec.Command("named-checkconf", tempConfFilePath).CombinedOutput()
	result := string(output)
	if commandError != nil {
		logger.Printf("Conf-Check failed with error: %v, command result was: %v", commandError, result)
		return false
	}

	return true
}

func updateZones(masterZone configreader.DnsMasterZone, tempDirName string) bool {
	serviceReloadRequired := false
	for _, zoneName := range masterZone.Zones {
		localZone := LocalZone{
			ZoneName:     zoneName,
			ZoneFileName: zoneName + ".db",
			TempZonePath: filehandler.BuildFilePathFromParts(tempDirName, zoneName+".db"),
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
