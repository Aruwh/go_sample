package mongodb

import "errors"

var (
	// ErrConnection is returned if a connection to a mongo server / cluster cannot be established.
	ErrConnection = errors.New("connection error")
	// ErrPing is returned if Client cannot ping the mongo server / cluster.
	ErrPing = errors.New("ping error")
	// ErrInvalidDatabaseURI is returned if the provided uri is invalid.
	ErrInvalidDatabaseURI = errors.New("invalid database uri")
	// ErrInvalidDatabaseName is returned if the provided database name is invalid.
	ErrInvalidDatabaseName = errors.New("invalid database name")

	ErrIDConversion = errors.New("unable to convert id to object id")
	ErrDecodeResult = errors.New("unable to decode result")
	ErrMongo        = errors.New("mongo error")
)
