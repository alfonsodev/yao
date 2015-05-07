package main

import (
	"database/sql"
	"fmt"
	UM "github.com/alfonsodev/yao/models/usermanager"
	"testing"
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
	user := UM.NewUsers(db)
	user.Username.Scan("Albert")
	_, err := user.Save()
	if err != nil {
		t.Error("Can't save\n" + err.Error())
	}
}
