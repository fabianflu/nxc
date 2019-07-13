package filefetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func FetchFile(filepath, nxtoken, baseurl string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", baseurl+filepath, nil)
	req.Header.Set("x-nx-token", nxtoken)
	res, _ := client.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

}
