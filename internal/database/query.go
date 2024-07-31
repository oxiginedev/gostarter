package database

import (
	"database/sql"
	"errors"
)

// Find retrieves the first matching instance of a model based on provided query
func (c *Connection) First(dest interface{}, query string, args ...interface{}) error {
	if err := c.Eager().Q().Where(query, args...).First(dest); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRecordNotFound
		}

		return err
	}

	return nil
}
