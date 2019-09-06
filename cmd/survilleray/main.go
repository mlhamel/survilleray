package main

import (
	"fmt"

	"github.com/mlhamel/survilleray/pkg/opensky"
)

const openskyURL = "https://opensky-network.org/api/states/all?lamin=%d&lamax=%d&lomin=%d&lomax=%d"

func main() {
	var r = opensky.NewRequest(openskyURL)

	vectors, err := r.GetPlanes()

	if err != nil {
		panic(err)
	}

	for i := 0; i < len(vectors); i++ {
		fmt.Println(vectors[i])
	}
}
