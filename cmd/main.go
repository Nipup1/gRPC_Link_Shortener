package main

import (
	"go/link_shortener/internal/app"
	"go/link_shortener/internal/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	application := app.New(cfg.GRPC.Port, cfg.StoragePath, cfg.InMemoryStorage)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCSrv.Stop()
}