package dns

import (
	"github.com/fabianflu/nxc/filefetcher"
	"github.com/fabianflu/nxc/filehandler"
	"os/exec"
)

type LocalZone struct {
	ZoneName     string ""
	ZoneFileName string ""
	TempZonePath string ""
}

func (zone LocalZone) downloadZone() error {
	logger.Print("Downloading file:", zone.ZoneFileName)
	zoneUrl := filehandler.BuildFilePathFromParts(clientConfig.BaseUrl, clientConfig.DnsConfig.RemotePaths.ZonePath, zone.ZoneFileName)
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
	targetZonePath := filehandler.BuildFilePathFromParts(clientConfig.DnsConfig.LocalPaths.ZonePath, zone.ZoneFileName)
	logger.Printf("Applying zone %v by copying it to: %v",
		zone.ZoneName, clientConfig.DnsConfig.LocalPaths.ZonePath+zone.ZoneFileName)
	if filehandler.IsFileExistent(targetZonePath) {
		logger.Printf("Zone %v already exists comparing files to see if update is needed", zone.ZoneName)
		if filehandler.AreFilesEqualByHash(zone.TempZonePath, targetZonePath) {
			logger.Printf("Your zone file: %v is up to date, no modification needed", zone.ZoneFileName)
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
