package usermanager

import (
	"database/sql"
	_ "github.com/lib/pq"
	"strings"
)

var UsersDb *sql.DB

type Users struct {
	Id          sql.NullInt64
	Username    sql.NullString
	Fullname    sql.NullString
	Image       sql.NullString
	Email       sql.NullString
	Location    sql.NullString
	Googleid    sql.NullString
	Googletoken sql.NullString
	Person      sql.NullString
	Joinedon    sql.NullInt64
}

func (obj *Users) Save() (*sql.Rows, error) {
	query := "INSERT INTO usermanager.users ( username, fullname, image, email, location, googleid, googletoken, person, joinedon) "
	query += " VALUES( $1, $2, $3, $4, $5, $6, $7, $8, $9) "

	return UsersDb.Query(query,
		getValue(obj.Username),
		getValue(obj.Fullname),
		getValue(obj.Image),
		getValue(obj.Email),
		getValue(obj.Location),
		getValue(obj.Googleid),
		getValue(obj.Googletoken),
		getValue(obj.Person),
		getValue(obj.Joinedon),
	)
}

func AllUsers() ([]Users, error) {
	var results []Users
	query := "SELECT * FROM usermanager.users "
	rows, err := UsersDb.Query(query)

	for rows.Next() {
		var u Users
		rows.Scan(
			&u.Id,
			&u.Username,
			&u.Fullname,
			&u.Image,
			&u.Email,
			&u.Location,
			&u.Googleid,
			&u.Googletoken,
			&u.Person,
			&u.Joinedon,
		)

		results = append(results, u)
	}

	return results, err
}

func NewUsers(db *sql.DB) Users {
	UsersDb = db
	var model Users
	return model
}

func UsersWhere(field string, condition string, value interface{}) *Query {
	q := Query{
		Schema: "usermanager",
		Table:  strings.ToLower("Users"),
	}

	q.Where = Clause{
		Field:     field,
		Condition: condition,
		Value:     value,
	}

	return &q
}

func UsersGet(rows *sql.Rows) []Modeler {
	var results []Modeler
	var row Users
	for rows.Next() {
		rows.Scan(&row.Id, &row.Username, &row.Fullname, &row.Image, &row.Email, &row.Location, &row.Googleid, &row.Googletoken, &row.Person, &row.Joinedon)
		results = append(results, &row)
	}

	return results
}
