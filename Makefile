DBNAME=yaotest
STRUCTFILE=fixtures/usermanager.sql
FIXTUREFILE=./fixtures/usermanager.users.data1.sql

structure:
	dropdb --if-exists ${DBNAME}
	createdb ${DBNAME}
	psql ${DBNAME} -f ${STRUCTFILE}

fixtures:
	psql $(DBNAME) -f ${FIXTUREFILE}

gen:
	-rm -rf ./models
	go run yao.go gen
test: structure fixtures 
	go test ./models_test.go -v 
all: structure fixtures gen
	go test ./... 

.PHONY: structure fixtures test gen all
