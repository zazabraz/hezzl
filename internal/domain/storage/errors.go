package storage

import (
	"errors"
)

var PgErrNoEffect = errors.New("error no effect")

var PgErrNoRows = errors.New("error no rows")
