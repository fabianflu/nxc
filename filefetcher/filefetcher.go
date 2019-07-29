package filefetcher

import (
	"fmt"
	"github.com/fabianflu/nxc/filehandler"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var logger = log.New(os.Stdout, "[file_fetcher] ", log.LstdFlags)

func FetchFile(filePath, nxToken, baseURL string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", baseURL, filePath)
	var headers = make(map[string]string)
	headers["x-nx-token"] = nxToken
	response, err := GetContentFromWeb(url, headers)
	if err != nil {
		logger.Printf("Failed to fetch Web-COntent %v", err)
		return nil, err
	}
	return response, nil
}

func GetContentFromWeb(sourceUrl string, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	request, requestBuildErr := http.NewRequest(http.MethodGet, sourceUrl, nil)
	if requestBuildErr != nil {
		logger.Printf("Could not create request: %v", requestBuildErr)
		return nil, requestBuildErr
	}
	if headers != nil {
		for key, value := range headers {
			request.Header.Set(key, value)
		}
	}
	response, requestErr := client.Do(request)
	if requestErr != nil {
		logger.Printf("Execution of request: %v, resulted in error: %v", request, requestErr)
		return nil, requestErr
	}
	defer response.Body.Close()
	logger.Printf("Response code is %v", response.StatusCode)
	result, ioError := ioutil.ReadAll(response.Body)
	if ioError != nil {
		logger.Printf("Could not read response body: %v", ioError)
		return nil, ioError
	}
	return result, nil
}

func DownloadFileFromWeb(targetPath string, sourceUrl string, checkIfModified bool, token string) error {
	var headers = make(map[string]string)
	headers["x-nx-token"] = token
	if checkIfModified {
		if modDate := filehandler.GetModifiedDate(targetPath); modDate != (time.Time{}) {
			headers["If-Modified-Since"] = modDate.Format(time.ANSIC)
			logger.Printf("Header set to: %v", modDate.Format(time.ANSIC))
		}
	}
	fileContent, err := GetContentFromWeb(sourceUrl, headers)
	if err != nil {
		logger.Printf("Failed to download file. Error: %v", err)
	}
	err = ioutil.WriteFile(targetPath, fileContent, 0644)
	return err
}
