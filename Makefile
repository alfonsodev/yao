DBNAME=booktown
fixtures:
	psql $(DBNAME) -f fixtures/booktown.sql
test: 
	go test ./... 

gen:
	go run yao.go gen

.PHONY: fixtures, test
