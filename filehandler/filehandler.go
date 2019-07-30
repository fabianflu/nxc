package filehandler

import (
	"crypto/sha1"
	"io"
	"io/ioutil"
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
	if fileInfo, err := os.Stat(filepath); err != nil {
		return nil, err
	} else {
		return fileInfo, nil
	}
}

func GetModifiedDate(filepath string) time.Time {
	info, err := GetFileInfo(filepath)
	if err != nil {
		logger.Printf("Failed to load FileInfo. Returning 0-value of time as modifiedDate")
		return time.Time{}
	}
	return info.ModTime()
}
func IsFileExistent(file string) bool {
	_, e := os.Open(file)
	return !os.IsNotExist(e)
}
func getFileHash(filepath string) string {
	file, e := os.Open(filepath)
	if e != nil {
		logger.Printf("Failed to hash file %v returning empty hash", filepath)
		return ""
	}
	fileContent, e := ioutil.ReadAll(file)
	sha1HashBuilder := sha1.New()
	return string(sha1HashBuilder.Sum(fileContent))
}

func AreFilesEqualByHash(file string, file2 string) bool {
	leftHash := getFileHash(file)
	return leftHash != "" && leftHash == getFileHash(file2)
}

func CopyOrOverwrite(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}
