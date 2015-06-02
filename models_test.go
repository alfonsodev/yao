package main

import (
	"database/sql"
	"fmt"
	//	"fmt"
	"testing"

	Users "github.com/alfonsodev/yao/models/usermanager/users"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "dbname=yaotest sslmode=disable")
	if err != nil {
		panic("Can't connect to database. \n" + err.Error())
	}
	fmt.Println("Init tests...")
	// TODO: Check that the connection works
	// Get credentials from database
}

func TestUsers(t *testing.T) {
	// Can insert a new Row
	user := Users.New(db)
	user.Username.Scan("Albert")
	user.Fullname.Scan("Einstein")
	user.Image.Scan("einstein.jpg")
	user.Email.Scan("albert@uzh.edu")
	user.Location.Scan("Zurich")
	user.Googleid.Scan("312393lk3j1lkjl123")
	user.Googletoken.Scan("312393lk3j1lkjl123")
	user.Person.Scan("{}")
	user.Joinedon.Scan(695510502)

	_, err := user.Save()
	if err != nil {
		t.Error("Can't save\n" + err.Error())
	}

	// Can Retrieve all rows
	users, errAll := Users.All()

	if errAll != nil || len(users) != 11 {
		//		fmt.Println(errAll.Error())
		t.Error("Can't retrieve the users ")
	}

	// The data was inserted correctly, I can access the fields
	if users[10].Email.String != "albert@uzh.edu" {
		t.Error("Insert didnt work ")
	}

	// // I can apply WHERE cluase
	students, errStu := Users.Where("Email", "LIKE", "%.edu").Get()

	for _, v := range students {
		if v.Email.String == "" {
			t.Error("Emails should not be emtpy")
		}
	}

	if errStu != nil {
		t.Error("I can't filter by Email")
	}

	if len(students) != 2 {
		t.Error("Something is wrong with the number of results")
	}

	// // Chain Multiple WHERE
	subStudents, errSub := Users.Where("Email", "LIKE", "%.edu").And("Location", "=", "Zurich").Get()
	for _, ss := range subStudents {
		if ss.Email.String == "" {
			t.Error("Emails should not be emtpy")
		}
	}

	if errSub != nil {
		t.Error("Can't chain two Where cluases")
	}

	if len(subStudents) != 1 || subStudents[0].Username.String != "Albert" {
		t.Error("Something is wrong with filtering with two Where clauses " + subStudents[0].Username.String)
	}

	// Chain Where, multiple ANDS and OR
	_, errUsers2 := Users.Where("Email", "LIKE", "%.edu").Or("Email", "LIKE", "%.co.uk").And("Joinedon", "=", 695510502).Get()
	if errUsers2 != nil {
		t.Errorf("Can't filter with Where Or and And")
	}

}
