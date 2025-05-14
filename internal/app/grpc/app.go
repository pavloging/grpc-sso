package grpcapp

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
func New(log *slog.Logger, authService authgrpc.Auth, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun runs gRPC server and panics if any error occurs.
func (a *App) MustRun() error {
	if err := a.Run(); err != nil {
		panic(err)
	}

	return nil
}

// Run starts gRPC server app.
func (a *App) Run() error {
	const fn = "grpcApp.Run"

	log := a.log.With(slog.String("fn", fn), slog.Int("port", a.port))

	// Обрабатываем tcp пакеты
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	log.Info("gRPC server started", slog.String("address", l.Addr().String()))

	log.Info("gRPC server starting", slog.String("address", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		log.Error("gRPC server failed", slog.String("error", err.Error()))
	}

	return nil
}

// Stop stops gRPC server app.
func (a *App) Stop() {
	const fn = "grpcApp.Stop"

	a.log.With(slog.String("fn", fn)).Info("stopping gRPC server", slog.Int("port", a.port))

	// Блокируем новые запросы и ждем завершения текущих => выключаем
	a.gRPCServer.GracefulStop()
}
