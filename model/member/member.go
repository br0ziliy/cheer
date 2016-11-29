// Package member
package member

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

var (
	// table is the table name.
	table = "member"
)

// Item defines the model.
type Item struct {
	ID        uint32         `db:"id"`
	IRC       string         `db:"ircnick"`
	Name      string         `db:"fullname"`
	GID       uint32         `db:"group_id"`
	Group     string         `db:"groupname"`
	CreatedAt mysql.NullTime `db:"created_at"`
	UpdatedAt mysql.NullTime `db:"updated_at"`
	DeletedAt mysql.NullTime `db:"deleted_at"`
}

type SanitizedItem struct {
	ID        uint32         `db:"id"`
	IRC       string         `db:"ircnick"`
	Name      string         `db:"fullname"`
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
		SELECT m.id, m.ircnick, m.fullname, m.group_id, g.name groupname, m.created_at, m.updated_at, m.deleted_at
		FROM %v m, ` + "`group` g" + `
		WHERE m.group_id = g.id
		  AND m.id = ?
		  AND m.deleted_at IS NULL
		LIMIT 1
		`, table),
		ID)
	return result, err == sql.ErrNoRows, err
}

// All gets all items.
func All(db Connection) ([]Item, bool, error) {
	var result []Item
	err := db.Select(&result, fmt.Sprintf(`
		SELECT m.id, m.ircnick, m.fullname, g.name groupname, m.created_at, m.updated_at, m.deleted_at
		FROM %v m, ` + "`group` g" + `
		WHERE m.group_id = g.id
		  AND m.deleted_at IS NULL
	`, table))
	return result, err == sql.ErrNoRows, err
}

// Get all irc nicknames
func GetNicks(db Connection) ([]SanitizedItem, bool, error) {
	var result []SanitizedItem
	err := db.Select(&result, fmt.Sprintf(`
		SELECT id, ircnick, fullname
		FROM %v
		WHERE deleted_at IS NULL
	`, table))
	return result, err == sql.ErrNoRows, err
}

// Get folks from your group
func GetMates(db Connection, ircnick string) ([]Item, bool, error) {
	var result []Item
	// select * from member where group_id = (select group_id from member where ircnick = 'maufart') and ircnick != 'maufart';
	err := db.Select(&result, fmt.Sprintf(`
	SELECT id, ircnick, fullname
	FROM %v
	WHERE group_id = (SELECT group_id
	                  FROM %v
			  WHERE ircnick = ?)
	  AND ircnick != ?
	`, table, table),
	ircnick, ircnick)
	return result, err == sql.ErrNoRows, err
}

func GetMyGroup(db Connection, ircnick string) (Item, bool, error) {
	result := Item{}
	err := db.Get(&result, fmt.Sprintf(`
		SELECT m.id, m.ircnick, m.fullname, g.name groupname, m.created_at, m.updated_at, m.deleted_at
		FROM %v m, ` + "`group` g" + `
		WHERE m.group_id = g.id
		  AND m.ircnick = ?
		  AND m.deleted_at IS NULL
		  LIMIT 1
	`, table),
	ircnick)
	return result, err == sql.ErrNoRows, err
}

// Create adds an item.
func Create(db Connection, ircnick string, fullname string, group_id string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(ircnick, fullname, group_id)
		VALUES
		(?,?,?)
		`, table),
		ircnick, fullname, group_id)
	return result, err
}

// Update makes changes to an existing item.
func Update(db Connection, ircnick string, fullname string, group_id string, ID string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		UPDATE %v
		SET ircnick = ?,
		    fullname = ?,
		    group_id = ?
		WHERE id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		ircnick, fullname, group_id, ID)
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
