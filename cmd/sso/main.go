package main

import (
	"fmt"
	"grpc-sso/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO: инициализировать объект логгера

	// TODO: инициализировать приложение (app)

	// TODO: запустить gRPC-сервер приложения
}
