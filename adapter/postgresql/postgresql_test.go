package postgresql

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "dbname=foo sslmode=disable")
	if err != nil {
		panic(err)
	}
}

func TestGetInformationSchema(t *testing.T) {
	fmt.Println("[postgresql] testing information schema")
	var p pq
	p.SetDb(db)
	info := p.GetInformationSchema("usermanager")
	fmt.Printf("\n[info]:%+v", info)
}

func TestGetPrimaryKey(t *testing.T) {
	// Test retreives correctly
	var p pq
	p.SetDb(db)
	expectedKeys := []string{"users_id", "orgs_id", "teams_id"}
	actualKeys := p.GetPrimaryKey("usermanager.users_orgs")
	if ok := reflect.DeepEqual(expectedKeys, actualKeys); !ok {
		t.Error("Expected keys doesnt match current keys")
	}

}
