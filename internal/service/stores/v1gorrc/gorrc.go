package v1gorrc

import (
	"errors"

	"github.com/gocraft/dbr"
)

var ErrNotFound = errors.New("not found")
var ErrDuplicateEntry = errors.New("duplicate entry for key")

// Error Codes for MySql
const MySqlDuplicateEntryNumber = 1062

type StoreAPI interface {
	UsersAPI
	TracksAPI
	CarsAPI
}

type Store struct {
	*base
}

type base struct {
	Connection *dbr.Connection
}

func NewStore(connection *dbr.Connection) *Store {
	s := Store{
		base: &base{
			Connection: connection,
		},
	}

	return &s
}
