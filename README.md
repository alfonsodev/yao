#YAO
###Yet Another Orm  

[![Build Status](https://drone.io/github.com/alfonsodev/yao/status.png)](https://drone.io/github.com/alfonsodev/yao/latest)   

Yao is an ORM for Golang, inspired in [Laravel eloquent ORM](http://laravel.com/docs/5.0/eloquent)
Currently only compatible with postgresql.

## Usage

### Generate structs 
`
go get github.com/alfonsodev/yao
` 
`
yao gen -ddbname  -o$(pwd)/models
`

### Import and connect
Mysql:  
```
    import "github.com/youpackage/models/tablename"
```
Postgresql:  
```
    import "github.com/youpackage/models/schemaname/tablename"
```

### Save 
```
    var user Users
    user.Email.Scan("yao@yao-orm.org")
    user.Save()
```
### Select all users from Users table  

```
    users, err := Users.All()
    for i, user := range users {
      fmt.Println(user.Email.String)
    }
```
### Where, and, or, Get() to get filtered results
```
    users, err := Users.Where("Email", "LIKE", "%.edu").And("Location", "=", "Zurich").Get()
```
