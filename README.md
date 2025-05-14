# gRPC App (SSO)

### Как запустить?

1. Переменуйте .env.example на .env
2. Дополните переменные, если их нет или они не верны
3. Напишите в терминале из корня: `go run cmd/sso/main.go`


### Архитектура

- Сервесный слой - internal/services/auth/auth.go

TODO: Добавить таблицу admin и перенести туда user-ов с isAdmin: 1