#YAO
 
## Yao's an ORM for Golang
Yao is an ORM for Golang, inspired in [Laravel eloquent ORM](http://laravel.com/docs/5.0/eloquent)
Yao generates a model for every table in yuor databse, making easier to deal with the most common CRUD boring operations.

### No Magic philosophy: 
The philosofy of Yao is to generate plain go code that has no other dependecies than the standard database/sql library.
Making easy to extend the models and understand what yao is doing. 
by using reflection during generation but not in run time. 

### Databases support:
Currently only Postgres is supported, but Yao implements an adapter pattern inspired by the core database/sql driver Go library ,
that make very easy to extend to different databses, have a look to [adapters/postgres/postgres.sql](http://github.com/alfonsodev/yao)
to have an idea what you need to implent. Basically your adapter needs to meet the Yao driver interface.

### Chainable calls:
Example: 
`.Where('email','=','me@mail.com').And('age','>',20).Or('name','=', 'Jhon').Get()`
## Install
`
go get github.com/alfonsodev/yao
`

## Generate your models 
`
yao gen -d dbname -u username -p password -H host 
`

Include
yao generates create a folder per each schema inside models, 
the default schema is called `public` so at least you'll have to import this one

`
import(
  publicSchema "github.com/youUser/yourpacker/models/public",
  userSchema "github.com/youUser/yourpacker/models/userschema",
)
`


## Insert 
  user = publicSchema.Users.New()
  user.Username.Scan("Username") 
  user.save() 

## Find all
  users = publicSchema.Users.All()
  
## Where
  users = publicSchema.Users.Where("email", "like", "*@gmail.com").Get()

## Where And Or 

	  users = publicSchema.Users
		.Where("email", "like", "*@gmail.com")
		.And("Age", ">", 23)
		.Get()

## Find by Primary Key 
  user = publicSchema.Users.Find(1)  

## Delete
  user = publicSchema.Users.Find(1)  
  user.Delete()
  // Or you can chain 
  publicSchema.Users.Find(1).Delete()

