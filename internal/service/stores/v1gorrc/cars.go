package v1gorrc

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gocraft/dbr"
	"github.com/google/uuid"
)

const TableCars = "cars"

var _ CarsAPI = (*Store)(nil) // Used just to show store implements interface

type CarsAPI interface {
	CreateCar(ctx context.Context, car Car) error

	GetCarByID(context.Context, uuid.UUID) (*Car, error)
	GetCarByShortName(context.Context, string) (*Car, error)

	GetCarsByTrack(context.Context, uuid.UUID) ([]Car, error)

	DeleteCar(context.Context, uuid.UUID) error
	HardDeleteCar(context.Context, uuid.UUID) error
}

type Car struct {
	Id          uuid.UUID      `db:"id"`
	Name        string         `db:"name"`
	ShortName   string         `db:"short_name"`
	Type        string         `db:"type"`
	Logo        string         `db:"logo"`
	Description string         `db:"description"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
	DeletedAt   sql.NullTime   `db:"deleted_at"`
	Track       uuid.UUID      `db:"track"`
	Password    sql.NullString `db:"password"`
}

func (s *Store) CreateCar(ctx context.Context, car Car) error {
	sr := s.Connection.NewSession(nil)

	stmt := sr.InsertInto(TableCars).
		Columns(
			"id",
			"name",
			"short_name",
			"type",
			"logo",
			"description",
			// "created_at",
			// "updated_at",
			// "deleted_at",
			"track",
			"password",
		).Record(car)

	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetCarByID(ctx context.Context, id uuid.UUID) (*Car, error) {
	sr := s.Connection.NewSession(nil)

	stmt := sr.Select("*").
		From(TableCars).
		Where("id = ?", id)

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var result Car
	err := stmt.LoadOneContext(ctx, &result)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}
	return &result, nil
}

func (s *Store) GetCarByShortName(ctx context.Context, name string) (*Car, error) {
	sr := s.Connection.NewSession(nil)

	stmt := sr.Select("*").
		From(TableCars).
		Where("short_name = ?", name)

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var result Car
	err := stmt.LoadOneContext(ctx, &result)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}
	return &result, nil
}

func (s *Store) GetCarsByTrack(ctx context.Context, track uuid.UUID) ([]Car, error) {
	sr := s.Connection.NewSession(nil)

	stmt := sr.Select("*").
		From(TableCars).
		Where("track = ?", track)

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var result []Car
	_, err := stmt.LoadContext(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Store) DeleteCar(ctx context.Context, id uuid.UUID) error {
	sr := s.Connection.NewSession(nil)

	stmt := sr.Update(TableCars).
		Set("deleted_at", time.Now()).
		Where("id = ?", id)

	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := stmt.ExecContext(ctx)
	return err
}

func (s *Store) HardDeleteCar(ctx context.Context, id uuid.UUID) error {
	sr := s.Connection.NewSession(nil)

	stmt := sr.DeleteFrom(TableCars).
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
