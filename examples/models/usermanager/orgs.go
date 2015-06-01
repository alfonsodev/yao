package usermanager

import (
	"database/sql"
	_ "github.com/lib/pq"
	"strings"
)

var OrgsDb *sql.DB

type Orgs struct {
	Id       sql.NullInt64
	Name     sql.NullString
	Website  sql.NullString
	Location sql.NullString
}

func (obj *Orgs) Save() (*sql.Rows, error) {
	query := "INSERT INTO usermanager.orgs ( name, website, location) "
	query += " VALUES( $1, $2, $3) "

	return OrgsDb.Query(query,
		getValue(obj.Name),
		getValue(obj.Website),
		getValue(obj.Location),
	)
}

func AllOrgs() ([]Orgs, error) {
	var results []Orgs
	query := "SELECT * FROM usermanager.orgs "
	rows, err := OrgsDb.Query(query)

	for rows.Next() {
		var u Orgs
		rows.Scan(
			&u.Id,
			&u.Name,
			&u.Website,
			&u.Location,
		)

		results = append(results, u)
	}

	return results, err
}

func NewOrgs(db *sql.DB) Orgs {
	OrgsDb = db
	var model Orgs
	return model
}

func OrgsWhere(field string, condition string, value interface{}) *Query {
	q := Query{
		Schema: "usermanager",
		Table:  strings.ToLower("Orgs"),
	}

	q.Where = Clause{
		Field:     field,
		Condition: condition,
		Value:     value,
	}

	return &q
}

func OrgsGet(rows *sql.Rows) []Modeler {
	var results []Modeler
	var row Orgs
	for rows.Next() {
		rows.Scan(&row.Id, &row.Name, &row.Website, &row.Location)
		results = append(results, &row)
	}

	return results
}
