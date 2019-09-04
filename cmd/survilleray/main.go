package main

import (
	"github.com/mlhamel/survilleray/pkg/opensky"
)

const openskyURL = "https://opensky-network.org/api/states/all?lamin=%d&lamax=%d&lomin=%d&lomax=%d"

func main() {
	var r = opensky.NewRequest(openskyURL)

	r.Execute()
}
