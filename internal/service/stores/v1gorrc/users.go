package v1gorrc

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gocraft/dbr"
	"github.com/google/uuid"
)

const TableUsers = "users"

var _ UsersAPI = (*Store)(nil) // Used just to show store implements interface

type UsersAPI interface {
	CreateUser(ctx context.Context, user User) error

	GetUserByID(context.Context, uuid.UUID) (*User, error)
	GetUserByUserName(context.Context, string) (*User, error)

	DeleteUser(context.Context, uuid.UUID) error
	HardDeleteUser(context.Context, uuid.UUID) error
}

type User struct {
	Id        uuid.UUID    `db:"id"`
	UserName  string       `db:"user_name"`
	Password  string       `db:"password"`
	RealName  string       `db:"real_name"`
	Email     string       `db:"email"`
	Rank      string       `db:"rank"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

func (s *Store) CreateUser(ctx context.Context, user User) error {
	sr := s.Connection.NewSession(nil)

	stmt := sr.InsertInto(TableUsers).
		Columns(
			"id",
			"user_name",
			"password",
			"real_name",
			"email",
			"rank",
			// "created_at",
			// "updated_at",
			// "deleted_at",
		).Record(user)

	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	sr := s.Connection.NewSession(nil)

	stmt := sr.Select("*").
		From(TableUsers).
		Where("id = ?", id)

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var result User
	err := stmt.LoadOneContext(ctx, &result)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}
	return &result, nil
}

func (s *Store) GetUserByUserName(ctx context.Context, name string) (*User, error) {
	sr := s.Connection.NewSession(nil)

	stmt := sr.Select("*").
		From(TableUsers).
		Where("user_name = ?", name)

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var result User
	err := stmt.LoadOneContext(ctx, &result)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}
	return &result, nil
}

func (s *Store) DeleteUser(ctx context.Context, id uuid.UUID) error {
	sr := s.Connection.NewSession(nil)

	stmt := sr.Update(TableUsers).
		Set("deleted_at", time.Now()).
		Where("id = ?", id)

	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := stmt.ExecContext(ctx)
	return err
}

func (s *Store) HardDeleteUser(ctx context.Context, id uuid.UUID) error {
	sr := s.Connection.NewSession(nil)

	stmt := sr.DeleteFrom(TableUsers).
		Where("id = ?", id)

	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := stmt.ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}
