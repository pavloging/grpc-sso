package grpcApp

import (
	"fmt"
	authgrpc "grpc-sso/internal/grpc/auth"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New creates new gRPC server app.
func New(log *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun runs gRPC server and panics if any error occurs.
func (app *App) MustRun() error {
	if err := app.Run(); err != nil {
		panic(err)
	}

	return nil
}

// Run starts gRPC server app.
func (app *App) Run() error {
	const fn = "grpcApp.Run"

	log := app.log.With(slog.String("fn", fn), slog.Int("port", app.port))

	// Обрабатываем tcp пакеты
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	log.Info("gRPC server started")

	// Запускаем сервер в горутине
	log.Info("gRPC server starting", slog.String("address", l.Addr().String()))
	if err := app.gRPCServer.Serve(l); err != nil {
		log.Error("gRPC server failed", slog.String("error", err.Error()))
	}

	return nil
}

// Stop stops gRPC server app.
func (app *App) Stop() {
	const fn = "grpcApp.Stop"

	app.log.With(slog.String("fn", fn)).Info("stopping gRPC server", slog.Int("port", app.port))

	// Блокируем новые запросы и ждем завершения текущих => выключаем
	app.gRPCServer.GracefulStop()
}
