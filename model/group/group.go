// Package group
package group

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

var (
	// table is the table name.
	table = "`group`"
)

// Item defines the model.
type Item struct {
	ID          uint32         `db:"id"`
	Name        string         `db:"name"`
	Description string         `db:"description"`
	CreatedAt   mysql.NullTime `db:"created_at"`
	UpdatedAt   mysql.NullTime `db:"updated_at"`
	DeletedAt   mysql.NullTime `db:"deleted_at"`
}

type SanitizedItem struct {
	ID          uint32         `db:"id"`
	Name        string         `db:"name"`
}

// Connection is an interface for making queries.
type Connection interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}

// ByID gets an item by ID.
func ByID(db Connection, ID string) (Item, bool, error) {
	result := Item{}
	err := db.Get(&result, fmt.Sprintf(`
		SELECT id, name, description, created_at, updated_at, deleted_at
		FROM %v
		WHERE id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		ID)
	return result, err == sql.ErrNoRows, err
}

// All gets all items.
func All(db Connection) ([]Item, bool, error) {
	var result []Item
	err := db.Select(&result, fmt.Sprintf(`
		SELECT id, name, description, created_at, updated_at, deleted_at
		FROM %v
		WHERE deleted_at IS NULL
		`, table))
	return result, err == sql.ErrNoRows, err
}

func GetGroupIdName(db Connection) ([]SanitizedItem, bool, error) {
	var result []SanitizedItem
	err := db.Select(&result, fmt.Sprintf(`
		SELECT id, name
		FROM %v
		WHERE deleted_at IS NULL
		`, table))
	return result, err == sql.ErrNoRows, err
}

// Create adds an item.
func Create(db Connection, name string, description string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(name, description)
		VALUES
		(?,?)
		`, table),
		name, description)
	return result, err
}

// Update makes changes to an existing item.
func Update(db Connection, name string, description string, ID string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		UPDATE %v
		SET name = ?,
		    description = ?
		WHERE id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		name, description, ID)
	return result, err
}

// DeleteHard removes an item.
func DeleteHard(db Connection, ID string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		DELETE FROM %v
		WHERE id = ?
			AND deleted_at IS NULL
		`, table),
		ID)
	return result, err
}

// Delete marks an item as removed.
func DeleteSoft(db Connection, ID string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		UPDATE %v
		SET deleted_at = NOW()
		WHERE id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		ID)
	return result, err
}
