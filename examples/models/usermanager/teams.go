package usermanager

import (
	"database/sql"
	_ "github.com/lib/pq"
	"strings"
)

var TeamsDb *sql.DB

type Teams struct {
	Id          sql.NullInt64
	Orgs_id     sql.NullInt64
	Name        sql.NullString
	Description sql.NullString
	Permission  sql.NullInt64
}

func (obj *Teams) Save() (*sql.Rows, error) {
	query := "INSERT INTO usermanager.teams ( orgs_id, name, description, permission) "
	query += " VALUES( $1, $2, $3, $4) "

	return TeamsDb.Query(query,
		getValue(obj.Orgs_id),
		getValue(obj.Name),
		getValue(obj.Description),
		getValue(obj.Permission),
	)
}

func AllTeams() ([]Teams, error) {
	var results []Teams
	query := "SELECT * FROM usermanager.teams "
	rows, err := TeamsDb.Query(query)

	for rows.Next() {
		var u Teams
		rows.Scan(
			&u.Id,
			&u.Orgs_id,
			&u.Name,
			&u.Description,
			&u.Permission,
		)

		results = append(results, u)
	}

	return results, err
}

func NewTeams(db *sql.DB) Teams {
	TeamsDb = db
	var model Teams
	return model
}

func TeamsWhere(field string, condition string, value interface{}) *Query {
	q := Query{
		Schema: "usermanager",
		Table:  strings.ToLower("Teams"),
	}

	q.Where = Clause{
		Field:     field,
		Condition: condition,
		Value:     value,
	}

	return &q
}

func TeamsGet(rows *sql.Rows) []Modeler {
	var results []Modeler
	var row Teams
	for rows.Next() {
		rows.Scan(&row.Id, &row.Orgs_id, &row.Name, &row.Description, &row.Permission)
		results = append(results, &row)
	}

	return results
}
