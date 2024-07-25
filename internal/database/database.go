package database

import (
	"context"
	"github/oxiginedev/gostarter/config"
	"net/url"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/pkg/errors"
)

type Connection struct {
	*pop.Connection
}

func Dial(cfg *config.DatabaseConfiguration) (*Connection, error) {
	dsn, err := url.Parse(cfg.BuildDSN())
	if err != nil {
		return nil, errors.Wrap(err, "parsing database url")
	}

	// var options map[string]string
	// if !util.IsStringEmpty(cfg.Op)

	conn, err := pop.NewConnection(&pop.ConnectionDetails{
		Dialect:         cfg.Driver,
		URL:             dsn.String(),
		Pool:            cfg.MaxPool,
		IdlePool:        cfg.MaxIdlePool,
		ConnMaxLifetime: time.Duration(cfg.MaxLifetime),
		ConnMaxIdleTime: time.Duration(cfg.MaxIdleTime),
	})
	if err != nil {
		return nil, errors.Wrap(err, "opening database connection")
	}

	if err := conn.Open(); err != nil {
		return nil, errors.Wrap(err, "checking database connection")
	}

	return &Connection{conn}, nil
}

type CommitWithError struct {
	Err error
}

func (e *CommitWithError) Error() string {
	return e.Err.Error()
}

func (e *CommitWithError) Cause() error {
	return e.Err
}

// NewCommitWithError creates an error that can be returned in a pop transaction
// without rolling back the transaction. This should only be used in cases where
// you want the transaction to commit but return an error message to the user.
func NewCommitWithError(err error) *CommitWithError {
	return &CommitWithError{Err: err}
}

func (c *Connection) Transaction(fn func(*Connection) error) error {
	if c.TX == nil {
		var rErr error
		if txErr := c.Connection.Transaction(func(tx *pop.Connection) error {
			err := fn(&Connection{tx})

			switch err.(type) {
			case *CommitWithError:
				rErr = err
				return nil
			default:
				return err
			}
		}); txErr != nil {
			return txErr
		}

		return rErr
	}

	return fn(c)
}

// WithContext returns a new connection with an updated context. This is
// typically used for tracing as the context contains trace span information.
func (c *Connection) WithContext(ctx context.Context) *Connection {
	return &Connection{c.Connection.WithContext(ctx)}
}
