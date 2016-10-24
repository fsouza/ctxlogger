package main

import (
	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	"github.com/fsouza/ctxlogger/examples/cats/service"
)

func main() {
	// showing 1 way of managing gizmo/config: importing from the environment
	var cfg service.Config
	config.LoadEnvConfig(&cfg)
	cfg.Server = &server.Config{}
	config.LoadEnvConfig(cfg.Server)

	server.Init("nyt-simple-proxy", cfg.Server)

	err := server.Register(service.NewSimpleService(&cfg, server.Log))
	if err != nil {
		server.Log.Fatal("unable to register service: ", err)
	}

	err = server.Run()
	if err != nil {
		server.Log.Fatal("server encountered a fatal error: ", err)
	}
}
