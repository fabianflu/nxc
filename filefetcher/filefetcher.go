package filefetcher

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, "[file_fetcher] ", log.LstdFlags)

func FetchFile(filePath, nxToken, baseURL string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", baseURL, filePath), nil)
	if err != nil {
		logger.Printf("Could not create request: %v", err)
		return nil, err
	}
	req.Header.Set("x-nx-token", nxToken)

	res, err := client.Do(req)
	if err != nil {
		logger.Printf("Could not execute request: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Printf("Could not read request body: %v", err)
		return nil, err
	}

	return body, nil
}
