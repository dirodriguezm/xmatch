package main

import (
	"github.com/dirodriguezm/xmatch/service/internal/di"
	httpservice "github.com/dirodriguezm/xmatch/service/internal/http_service"
)

func main() {
	ctr := di.BuildServiceContainer()
	var httpServer httpservice.HttpServer
	ctr.Resolve(&httpServer)
	httpServer.InitServer()
}
