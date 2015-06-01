package usermanager

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

type Modeler interface {
	Save() (*sql.Rows, error)
}

type Query struct {
	Schema  string
	Table   string
	Where   Clause
	Clauses []Clause
	Limit   int64
}

type Clause struct {
	Field       string
	Condition   string
	Value       interface{}
	Connector   string
	Parenthesis string
}

func panicIf(err error) {
	if err != nil {
		panic(err.Error())
	}
}
func getValue(v driver.Valuer) interface{} {
	o, err := v.Value()
	panicIf(err)

	return o
}

func (q *Query) And(field string, condition string, value interface{}) *Query {
	q.Clauses = append(q.Clauses, Clause{
		Field:     field,
		Condition: condition,
		Value:     value,
		Connector: "AND",
	})

	return q
}
func (q *Query) Or(field string, condition string, value interface{}) *Query {

	q.Clauses = append(q.Clauses, Clause{
		Field:     field,
		Condition: condition,
		Value:     value,
		Connector: "OR",
	})

	return q
}

func (q *Query) Get() ([]Modeler, error) {
	var values []interface{}
	sql := "SELECT * FROM " + q.Schema + "." + q.Table + " WHERE "
	// TODO: We are asuming there is always a Where clause
	sql += fmt.Sprintf("%s %s $1", q.Where.Field, q.Where.Condition)
	values = append(values, q.Where.Value)
	for i, v := range q.Clauses {
		if v.Connector != "" {
			sql += " " + v.Connector + " "
		}
		sql += fmt.Sprintf("%s %s $%v", v.Field, v.Condition, i+2)
		values = append(values, v.Value)
	}

	rows, err := UsersDb.Query(sql, values...)
	var results []Modeler
	switch q.Table {
	case "orgs":
		results = OrgsGet(rows)
	case "teams":
		results = TeamsGet(rows)
	case "users_orgs":
		results = Users_orgsGet(rows)
	case "users":
		results = UsersGet(rows)
	case "envs":
		results = EnvsGet(rows)

	}

	return results, err
}
