package usermanager

import (
	"database/sql"
	_ "github.com/lib/pq"
	"strings"
)

var EnvsDb *sql.DB

type Envs struct {
	Id       sql.NullInt64
	Orgs_id  sql.NullInt64
	Users_id sql.NullInt64
	Name     sql.NullString
	Doc      sql.NullString
}

func (obj *Envs) Save() (*sql.Rows, error) {
	query := "INSERT INTO usermanager.envs ( orgs_id, users_id, name, doc) "
	query += " VALUES( $1, $2, $3, $4) "

	return EnvsDb.Query(query,
		getValue(obj.Orgs_id),
		getValue(obj.Users_id),
		getValue(obj.Name),
		getValue(obj.Doc),
	)
}

func AllEnvs() ([]Envs, error) {
	var results []Envs
	query := "SELECT * FROM usermanager.envs "
	rows, err := EnvsDb.Query(query)

	for rows.Next() {
		var u Envs
		rows.Scan(
			&u.Id,
			&u.Orgs_id,
			&u.Users_id,
			&u.Name,
			&u.Doc,
		)

		results = append(results, u)
	}

	return results, err
}

func NewEnvs(db *sql.DB) Envs {
	EnvsDb = db
	var model Envs
	return model
}

func EnvsWhere(field string, condition string, value interface{}) *Query {
	q := Query{
		Schema: "usermanager",
		Table:  strings.ToLower("Envs"),
	}

	q.Where = Clause{
		Field:     field,
		Condition: condition,
		Value:     value,
	}

	return &q
}

func EnvsGet(rows *sql.Rows) []Modeler {
	var results []Modeler
	var row Envs
	for rows.Next() {
		rows.Scan(&row.Id, &row.Orgs_id, &row.Users_id, &row.Name, &row.Doc)
		results = append(results, &row)
	}

	return results
}
