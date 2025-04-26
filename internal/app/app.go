package app

import (
	grpcApp "grpc-sso/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcApp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// TODO: инициализировать приложение (storage)

	// TODO: init auth service (auth)

	grpcApp := grpcApp.New(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
