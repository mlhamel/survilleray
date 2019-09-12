package main

import (
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/server"
)

func main() {
	c := config.NewConfig()
	r := server.NewRouter()

	r.Run(c.HttpPort())
}