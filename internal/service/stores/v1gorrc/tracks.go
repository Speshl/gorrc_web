package v1gorrc

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gocraft/dbr"
	"github.com/google/uuid"
)

const TableTracks = "tracks"

var _ TracksAPI = (*Store)(nil) // Used just to show store implements interface

type TracksAPI interface {
	CreateTrack(ctx context.Context, track Track) error

	GetTrackByID(context.Context, uuid.UUID) (*Track, error)
	GetTrackByShortName(context.Context, string) (*Track, error)

	GetTracks(context.Context) ([]Track, error)

	DeleteTrack(context.Context, uuid.UUID) error
	HardDeleteTrack(context.Context, uuid.UUID) error
}

type Track struct {
	Id          uuid.UUID    `db:"id"`
	Name        string       `db:"name"`
	ShortName   string       `db:"short_name"`
	Type        string       `db:"type"`
	Logo        string       `db:"logo"`
	Description string       `db:"description"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
}

func (s *Store) CreateTrack(ctx context.Context, track Track) error {
	sr := s.Connection.NewSession(nil)

	stmt := sr.InsertInto(TableTracks).
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
		).Record(track)

	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetTrackByID(ctx context.Context, id uuid.UUID) (*Track, error) {
	sr := s.Connection.NewSession(nil)

	stmt := sr.Select("*").
		From(TableTracks).
		Where("id = ?", id)

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var result Track
	err := stmt.LoadOneContext(ctx, &result)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}
	return &result, nil
}

func (s *Store) GetTrackByShortName(ctx context.Context, name string) (*Track, error) {
	sr := s.Connection.NewSession(nil)

	stmt := sr.Select("*").
		From(TableTracks).
		Where("short_name = ?", name)

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var result Track
	err := stmt.LoadOneContext(ctx, &result)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}
	return &result, nil
}

func (s *Store) GetTracks(ctx context.Context) ([]Track, error) {
	sr := s.Connection.NewSession(nil)

	stmt := sr.Select("*").
		From(TableTracks).
		Where("deleted_at IS null")

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var result []Track
	_, err := stmt.LoadContext(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Store) DeleteTrack(ctx context.Context, id uuid.UUID) error {
	sr := s.Connection.NewSession(nil)

	stmt := sr.Update(TableTracks).
		Set("deleted_at", time.Now()).
		Where("id = ?", id)

	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := stmt.ExecContext(ctx)
	return err
}

func (s *Store) HardDeleteTrack(ctx context.Context, id uuid.UUID) error {
	sr := s.Connection.NewSession(nil)

	stmt := sr.DeleteFrom(TableTracks).
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
