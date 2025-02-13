package app

import (
	"go/link_shortener/internal/app/grpcapp"
	"go/link_shortener/internal/service/link"
	inmemory "go/link_shortener/internal/storage/in_memory"
	"go/link_shortener/internal/storage/postgres"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(port int, storagePath string, isMemoryStorage bool) *App{
	var linkService *link.Link
	if isMemoryStorage{
		storage := inmemory.New()
		linkService = link.New(storage)
	} else{
		storage, err := postgres.New(storagePath)
		if err != nil{
			panic(err)
		}
		linkService = link.New(storage)
	}

	grpcApp := grpcapp.New(port, linkService)
	
	return &App{
		GRPCSrv: grpcApp,
	}
}