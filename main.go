package main

import (
	"fmt"
	"github.com/fabianflu/nxc/configreader"
	"github.com/fabianflu/nxc/filefetcher"
)

func main() {
	config := configreader.ReadConfig("config.json")
	filefetcher.FetchFile("nx.json", config.NXToken, config.BaseUrl)
	fmt.Printf("baseUrl: %s, token: %s\n", config.BaseUrl, config.NXToken)
}
