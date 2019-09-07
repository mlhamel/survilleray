package main

import (
	"fmt"

	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/opensky"
)

const openskyURL = "https://opensky-network.org/api/states/all?lamin=%d&lamax=%d&lomin=%d&lomax=%d"

func main() {
	c := config.NewConfig()
	db := c.DB()

	var r = opensky.NewRequest(openskyURL)

	vectors, err := r.GetPlanes()

	if err != nil {
		panic(err)
	}

	for i := 0; i < len(vectors); i++ {
		fmt.Println(vectors[i])
		db.NewRecord(vectors[i])
		db.Create(&vectors[i])
	}

	defer db.Close()
}
