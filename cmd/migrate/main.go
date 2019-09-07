package main

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/mlhamel/survilleray/pkg/config"

	"github.com/mlhamel/survilleray/pkg/models"
)

func main() {
	c := config.NewConfig()

	fmt.Printf("Migrating %s\n", c.DSN())

	c.DB().Debug().AutoMigrate(&models.Vector{}) //Database migration
}
