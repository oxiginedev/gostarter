package migrations

import (
	"fmt"
	"github/oxiginedev/gostarter/internal/database"

	"gorm.io/gorm"
)

func AutoMigrate(db database.Database) error {
	err := db.GetDB().AutoMigrate(migrateTables()...)
	if err != nil {
		return err
	}

	for _, alterColumn := range migrateColumnChanges() {
		err := alterColumn.apply(db.GetDB())
		if err != nil {
			return fmt.Errorf("error migrating changes to %s for column [%s]", alterColumn.TableName, alterColumn.Column)
		}
	}

	return nil
}

// Register migrations here
func migrateTables() []interface{} {
	return []interface{}{
		&database.User{},
	}
}

func migrateColumnChanges() []AlterColumn {
	return []AlterColumn{
		//
	}
}

type AlterColumn struct {
	Model     interface{}
	TableName string
	Column    string
	Type      string
}

func (a *AlterColumn) apply(db *gorm.DB) error {
	query := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s USING %s::%s", a.TableName, a.Column, a.Type, a.Column, a.Type)
	if err := db.Exec(query).Error; err != nil {
		return err
	}

	if err := db.Migrator().AlterColumn(a.Model, a.Column); err != nil {
		return err
	}

	return nil
}
