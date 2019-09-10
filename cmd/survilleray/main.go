package main

import (
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/jobs"
)

func main() {
	c := config.NewConfig()

	job := jobs.NewAcquisition(c)
	job.Run()
}
