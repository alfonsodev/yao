package util

import (
	"fmt"
	"reflect"
	"strings"
)

//returns
// keys , string that contains field names comma separated "name, surname, address"
// placeholders, string containing list of placeholders"($1, $2, $3)"
// values, []interface{}, contains every field value of strc
func InsertHelper(strc interface{}) (string, string, []interface{}) {
	val := reflect.ValueOf(strc).Elem()
	keys := ""
	placeholders := ""
	var values []interface{}

	for i := 0; i < val.NumField(); i++ {
		fieldName := strings.ToLower(val.Type().Field(i).Name)
		//TODO: This is not a good way of ignoring fields, we can come up with a general rule that could be
		// - private fields for instance should be ignored
		// - Primary key should be ignore
		if fieldName == "yaodb" || fieldName == "id" || fieldName == "yaowhereclauses" {
			continue
		}
		// fmt.Println("[YAO]: " + fieldName + " >>")
		valueField := val.Field(i)
		values = append(values, valueField.Field(0).Addr().Interface())

		// if valueField.Field(0).IsValid() {
		//TODO: get sql tag if exists the fieldName  tag.Get("tag_name")
		// tag := val.Type().Field(i).Tag

		yaoTag := val.Type().Field(i).Tag.Get("yao")
		//TODO: this is wring , should be read from the tag,
		if yaoTag != "pk" {
			keys += fmt.Sprintf(", %s", fieldName)
		}
		// fmt.Printf(":::::%v ", valueField.Type().Name())
		// }
		// keys += fmt.Sprintf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"))
	}

	if string(keys[0]) == "," {
		keys = keys[1:]
	}

	for j := 1; j <= len(values); j++ {

		placeholders += fmt.Sprintf(", $%v", j)
	}

	placeholders = placeholders[1:]

	return keys, placeholders, values
}

// returns
// keys, string, of keys + placeholder as "name=$1, address=$2, email=$3"
// v, []interface{}, contains every field value of strc
func UpdateHelper(u interface{}) (string, []interface{}) {
	val := reflect.ValueOf(u).Elem()
	keys := ""
	v := make([]interface{}, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		if true {

			//			fmt.Printf("\n[VALUE]:%+v\n", valueField.String())
			// TODO: get sql tag if exists the fieldName  tag.Get("tag_name"
			// tag := val.Type().Field(i).Tag

			fieldName := strings.ToLower(val.Type().Field(i).Name)
			if i > 0 {
				keys += fmt.Sprintf(", %s = $%v", fieldName, i+1)
			}
			v[i] = valueField.Field(0).Addr().Interface()
			//			fmt.Printf("[-]%+v", v[i])
		}
	}

	if string(keys[0]) == "," {
		keys = keys[1:]
	}

	return keys, v
}

// func main() {
//  var b User
//  b.Name = "a"
//  keys, _ := Reflect(&b)
//  q := "UPDATE " + reflect.TypeOf(b).Name() + " SET " + keys
//  fmt.Printf(q)
// }
