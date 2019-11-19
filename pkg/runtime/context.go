package runtime

import (
	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/config"
)

type Context struct {
	cfg      *config.Config
	database *gorm.DB
}

func NewContext(cfg *config.Config, database *gorm.DB) *Context {
	return &Context{cfg: cfg, database: database}
}

func (c *Context) Config() *config.Config {
	return c.cfg
}

func (c *Context) Database() *gorm.DB {
	if c.database != nil {
		return c.database
	}

	database, err := gorm.Open("postgres", c.cfg.DSN())

	if err != nil {
		panic(err)
	}

	c.database = database

	return c.database
}
