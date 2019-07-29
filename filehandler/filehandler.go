package filehandler

import (
	"log"
	"os"
	"time"
)

var logger = log.New(os.Stdout, "[file_handler] ", log.LstdFlags)

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func GetFileInfo(filepath string) (os.FileInfo, error) {
	if fileInfo, err := os.Stat(filepath); os.IsNotExist(err) {
		logger.Printf("File with path: %v does not exist. Error: %v", filepath, err)
		return nil, err
	} else {
		return fileInfo, nil
	}
}

func GetModifiedDate(filepath string) time.Time {
	info, err := GetFileInfo(filepath)
	if err != nil {
		logger.Printf("Failed to load FileInfo return 0-value of time")
		return time.Time{}
	}
	return info.ModTime()
}
