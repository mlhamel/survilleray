package main

import (
	"fmt"

	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/server"
)

func main() {
	c := config.NewConfig()
	r := server.NewRouter()

	r.Run(buildConnexionString(c))
}

func buildConnexionString(c *config.Config) string {
	return fmt.Sprintf(":%s", c.HttpPort())
}
