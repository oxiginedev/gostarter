package database

import (
	"database/sql"

	"github.com/pkg/errors"
)

func (c *Connection) Find(dest interface{}, query string, args ...interface{}) error {
	if err := c.Eager().Q().Where(query, args...).First(dest); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return ErrRecordNotFound
		}

		return errors.Wrap(err, "error finding model")
	}

	return nil
}
