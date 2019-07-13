package main

import (
	"github.com/fabianflu/nxc/filefetcher"
)

func main() {
	filefetcher.FetchFile("nx.json", "", "https://nx-tokenservice.pegnu.workers.dev/")

}
