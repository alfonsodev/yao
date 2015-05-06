package postgresql

import (
	"database/sql"
	"fmt"
	g "github.com/alfonsodev/yao/generate"
	_ "github.com/lib/pq"
)

//Maps postgres types with go sql package types
//From http://www.postgresql.org/docs/9.4/static/datatype-numeric.html
var typeMap map[string]string

// pq is a struct that meets generate.YaoDriver interface
type pq struct {
	init bool
	db   *sql.DB
}

func (p *pq) SetDb(db *sql.DB) {
	p.init = true
	p.db = db
}

// drivers have to register themselves
func init() {
	var driver pq

	g.Register("postgres", &driver)

	typeMap = map[string]string{
		"integer":           "sql.NullInt64",
		"smallint":          "sql.NullInt64",
		"bigint":            "sql.NullInt64",
		"decimal":           "sql.NullInt64",
		"numeric":           "sql.NullInt64",
		"real":              "sql.NullFloat64",
		"double precision":  "sql.NullFloat64",
		"smallserial":       "sql.NullInt64",
		"serial":            "sql.NullInt64",
		"bigserial":         "sql.NullInt64",
		"boolean":           "sql.NullBool",
		"character":         "sql.NullString",
		"character varying": "sql.NullString",
		"json":              "sql.NullString",
		"timestamp with time zone": "sql.NullString",
	}
}

func (p *pq) GetSchemas() []string {
	var schemas []string
	var schemaName string
	query := `SELECT schema_name FROM information_schema.schemata
		WHERE schema_name NOT LIKE $1 AND schema_name NOT LIKE $2`
	rows, err := p.db.Query(query, "pg%", "information_schema")
	PanicIf(err)
	for rows.Next() {
		PanicIf(rows.Scan(&schemaName))
		schemas = append(schemas, schemaName)
	}

	return schemas
}

func (p *pq) GetPrimaryKey(tableName string) []string {
	var output []string
	query := `SELECT a.attname FROM   pg_index i
JOIN   pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
WHERE  i.indrelid = '%s'::regclass
AND    i.indisprimary`
	var colName string
	rows, err := p.db.Query(fmt.Sprintf(query, tableName))
	PanicIf(err)
	for rows.Next() {
		PanicIf(rows.Scan(&colName))
		output = append(output, colName)
	}

	return output
}

func InSlice(x string, a []string) bool {
	for _, v := range a {
		if v == x {
			return true
		}
	}
	return false
}

func (p *pq) GetInformationSchema(schemaName string) map[string][]g.FieldInfo {
	informationSchema := make(map[string][]g.FieldInfo)
	fmt.Println("[postgresql] Getting information for schema: " + schemaName)
	if schemaName == "" {
		schemaName = "public" //default schema in postgresql
	}

	selectTables := "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = $1"
	selectColumns := `SELECT column_name, is_nullable, data_type
  FROM information_schema.columns
  WHERE table_schema = $1
  AND table_name = $2`

	tables, err1 := p.db.Query(selectTables, schemaName)
	PanicIf(err1)
	//each table
	for tables.Next() {
		var name string
		var colName, nullable, dataType string
		var fields []g.FieldInfo

		err := tables.Scan(&name)
		PanicIf(err)
		fmt.Printf("Table: %s\n", name)

		cols, e := p.db.Query(selectColumns, schemaName, name)
		PanicIf(e)
		keys := p.GetPrimaryKey(schemaName + "." + name)

		//for each column
		for cols.Next() {
			KeyInfo := ""
			if InSlice(colName, keys) {
				KeyInfo = "pk"
			}
			PanicIf(cols.Scan(&colName, &nullable, &dataType))
			var ok bool
			dataType, ok = typeMap[dataType]
			if !ok {
				panic("Type " + dataType + " not implemented.")
			}
			fields = append(fields, g.FieldInfo{colName, nullable, dataType, KeyInfo})
			fmt.Printf("%s-%s-%s-%s\n", colName, nullable, dataType, KeyInfo)
		}

		informationSchema[name] = fields
	}

	return informationSchema
}

func PanicIf(err error) {
	if err != nil {
		fmt.Println("[postgresql]: " + err.Error())
		panic(err)
	}
}
