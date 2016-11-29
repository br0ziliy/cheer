// Package cheer
package cheer

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

var (
	// table is the table name.
	table = "cheer"
)

// Item defines the model.
type Item struct {
	ID        uint32         `db:"id"`
	Name      string         `db:"name"`
	CreatedAt mysql.NullTime `db:"created_at"`
	UpdatedAt mysql.NullTime `db:"updated_at"`
	DeletedAt mysql.NullTime `db:"deleted_at"`
}

// Connection is an interface for making queries.
type Connection interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}

// All gets all items.
func All(db Connection, ircnick string) ([]Item, bool, error) {
	var result []Item
	// select * from member where group_id = (select group_id from member where ircnick = 'maufart') and ircnick != 'maufart';
	err := db.Get(&result, fmt.Sprintf(`
		SELECT id, name, created_at, updated_at, deleted_at
		FROM %v, 
		WHERE 
		  AND deleted_at IS NULL
		`, table),
		ircnick)
	return result, err == sql.ErrNoRows, err
}

// Create adds an item.
func Create(db Connection, from_id string, to_id string, points string, cheer string, referer string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(points, cheer, from_id, to_id, referer)
		VALUES
		(?,?,?,?,?)
		`, table),
		points, cheer, from_id, to_id, referer)
	return result, err
}

