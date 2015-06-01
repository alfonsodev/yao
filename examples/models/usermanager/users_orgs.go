package usermanager

import (
	"database/sql"
	_ "github.com/lib/pq"
	"strings"
)

var Users_orgsDb *sql.DB

type Users_orgs struct {
	Users_id sql.NullInt64
	Orgs_id  sql.NullInt64
	Teams_id sql.NullInt64
}

func (obj *Users_orgs) Save() (*sql.Rows, error) {
	query := "INSERT INTO usermanager.users_orgs () "
	query += " VALUES() "

	return Users_orgsDb.Query(query)
}

func AllUsers_orgs() ([]Users_orgs, error) {
	var results []Users_orgs
	query := "SELECT * FROM usermanager.users_orgs "
	rows, err := Users_orgsDb.Query(query)

	for rows.Next() {
		var u Users_orgs
		rows.Scan(
			&u.Users_id,
			&u.Orgs_id,
			&u.Teams_id,
		)

		results = append(results, u)
	}

	return results, err
}

func NewUsers_orgs(db *sql.DB) Users_orgs {
	Users_orgsDb = db
	var model Users_orgs
	return model
}

func Users_orgsWhere(field string, condition string, value interface{}) *Query {
	q := Query{
		Schema: "usermanager",
		Table:  strings.ToLower("Users_orgs"),
	}

	q.Where = Clause{
		Field:     field,
		Condition: condition,
		Value:     value,
	}

	return &q
}

func Users_orgsGet(rows *sql.Rows) []Modeler {
	var results []Modeler
	var row Users_orgs
	for rows.Next() {
		rows.Scan(&row.Users_id, &row.Orgs_id, &row.Teams_id)
		results = append(results, &row)
	}

	return results
}
