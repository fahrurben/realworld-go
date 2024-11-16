package main

import (
	"github.com/fahrurben/realworld-go/cmd/server"
	"github.com/fahrurben/realworld-go/pkg/config"
)

func main() {

	// setup various configuration for app
	config.LoadAllConfigs(".env")

	server.Serve()
}
