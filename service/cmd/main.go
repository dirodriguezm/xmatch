package main

import (
	"github.com/dirodriguezm/xmatch/service/internal/di"
	httpservice "github.com/dirodriguezm/xmatch/service/internal/http_service"

	"github.com/golobby/container/v3"
)

func main() {
	di.ContainerBuilder()
	var httpServer httpservice.HttpServer
	container.MustResolve(container.Global, &httpServer)
	httpServer.InitServer()
}
