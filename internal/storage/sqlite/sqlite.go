package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"grpc-sso/internal/domain/models"
	"grpc-sso/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New creates new sqlite storage.
func New(storagePath string) (*Storage, error) {
	const fn = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

// SaveUser saves user to db.
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// TODO: Переписать на транзакцию
	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// User returns user by email.
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const fn = "storage.sqlite.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", fn, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash) // Получаем результат запроса
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", fn, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", fn, err)
	}

	return user, nil
}

//func (s *Storage) SavePermission(ctx context.Context, userID int64, permission models.Permission, appID string) error {
//	const op = "storage.sqlite.SavePermission"
//
//	stmt, err := s.db.Prepare("INSERT INTO permissions(user_id, permission, app_id) VALUES(?, ?, ?)")
//	if err != nil {
//		return fmt.Errorf("%s: %w", op, err)
//	}
//
//	_, err = stmt.ExecContext(ctx, userID, permission, appID)
//	if err != nil {
//		return fmt.Errorf("%s: %w", op, err)
//	}
//
//	return nil
//}

// App returns app by id.
func (s *Storage) App(ctx context.Context, id int) (models.App, error) {
	const op = "storage.sqlite.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

	var app models.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.sqlite.IsAdmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool

	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
