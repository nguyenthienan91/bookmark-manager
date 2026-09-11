package main

import (
	"github.com/nguyenthienan91/bookmark-manager/internal/api"
)	

// @title           Bookmark Manager API
// @version         1.0.0
// @description     API for Bookmark Manager application.
// @BasePath  			/

func main() {
	// create app config
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	app := api.NewEngine(cfg)
	err = app.Start()
	if err != nil {
		panic(err)
	}
}