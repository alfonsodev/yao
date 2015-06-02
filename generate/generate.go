package generate

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	fs "github.com/alfonsodev/yao/filesystem"
	_ "github.com/lib/pq"
)

var db *sql.DB

type FieldInfo struct {
	Name     string
	Nullable string
	Datatype string
	KeyInfo  string
}

type TemplateData struct {
	Name           string
	StructFields   string
	ScanFields     string
	Schema         string
	Table          string
	Keys           string
	Placeholders   string
	SaveFields     string
	JsonFields     string
	SwitchForGet   string
	AllFieldsByRef string
}

type YaoDriver interface {
	SetDb(db *sql.DB)
	GetInformationSchema(schemaname string) map[string][]FieldInfo
	GetSchemas() []string
}

// posible drivers
var drivers = make(map[string]YaoDriver)

// yao is the currect driver in use
var yao YaoDriver

// Register makes a yao driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, driver YaoDriver) {
	if driver == nil {
		panic("yao: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("yao: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

// Return a slice of strings with the registered driver names
func Drivers() []string {
	var output []string
	for k, _ := range drivers {
		output = append(output, k)
	}

	return output
}

func Open(driverName string, params string) (*sql.DB, error) {
	var err error
	db, err = sql.Open(driverName, params)
	yao = drivers[driverName]
	yao.SetDb(db)
	return db, err
}

func UcFirst(s string) string {
	s = strings.ToLower(s)
	return strings.ToUpper(s[:1]) + s[1:]
}

// Genreate one model file per table
func Generate(schemaname string, output string) {
	fs.SetWorkingDirectory(output)
	err := fs.CreateModelsFolder(output)
	if err != nil {
		fmt.Println("[YAO] ./models folder already exist")
		fmt.Println(err.Error())
	}

	schemas := yao.GetSchemas()
	if len(schemas) > 0 {
		for _, schename := range schemas {
			info := yao.GetInformationSchema(schename)
			generateModelFromSchema(schename, info)
		}
		return
	}
}

func panicIfErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func genererateQueryFile(schemaname string, info map[string][]FieldInfo) {
	data := TemplateData{Schema: schemaname}

	for k, _ := range info {
		data.SwitchForGet += fmt.Sprintf("		case \"%v\": results = %vGet(rows) \n", strings.ToLower(k), UcFirst(k))
	}

	// Create static functions file (query.go)
	tmpl, err := template.ParseFiles(os.Getenv("GOPATH") + "/src/github.com/alfonsodev/yao/template/.tmpl")
	panicIfErr(err)
	out := new(bytes.Buffer)
	err = tmpl.Execute(out, data)
	fs.CreateQueryFile(schemaname, string(out.Bytes()))
	panicIfErr(err)
}

func generateModelFromSchema(schemaname string, info map[string][]FieldInfo) {
	err := fs.CreateSchemaFolder(schemaname)
	if err != nil {
		fmt.Println("[YAO] can't create ./models/" + schemaname + " , folder already exist?")
	}

	for k, v := range info {
		fs.CreateModelFile(schemaname, strings.ToLower(k), PrintModel(k, v))
		cmd := exec.Command("go", "fmt", "./models/"+schemaname+"/"+k+"/"+UcFirst(k)+".go")
		err := cmd.Run()
		if err != nil {
			fmt.Println(string(err.Error()))
		}
	}
}

func PrintQuery(tables []string, fields []FieldInfo) string {
	var out string
	out += "switch q.Table { "
	for _, v := range tables {
		out += "case '" + v + "' :"
		out += ""
	}
	out += "}" // close switch

	return out

}

func PrintModel(name string, fields []FieldInfo) string {

	data := TemplateData{
		Name:   UcFirst(name),
		Schema: "usermanager",
		Table:  strings.ToLower(name),
	}

	for i, v := range fields {
		data.StructFields += "   " + UcFirst(v.Name) + " " + v.Datatype + "\n"
		data.ScanFields += "        &u." + UcFirst(v.Name) + ", \n"
		data.AllFieldsByRef += "&row." + UcFirst(v.Name) + ", "
		if v.KeyInfo != "pk" {
			data.SaveFields += "          getValue(obj." + UcFirst(v.Name) + "), \n"
			if v.Name == "json" {
				data.JsonFields += "obj.Json.String = \"{}\" \n"
			}

			// fmt.Println(v.Name + " : keyinfo : " + v.KeyInfo + "\n")
			data.Placeholders += fmt.Sprintf(", $%v", i)
			data.Keys += ", " + strings.ToLower(v.Name)
		}
	}

	data.AllFieldsByRef = strings.TrimRight(data.AllFieldsByRef, ", ")

	if len(data.Placeholders) > 0 {
		data.Placeholders = data.Placeholders[1:]
	}

	if len(data.Keys) > 0 {
		data.Keys = data.Keys[1:]
	}

	tmpl, err := template.ParseFiles(os.Getenv("GOPATH") + "/src/github.com/alfonsodev/yao/template/model.tmpl")
	if err != nil {
		panic(err)
	}

	out := new(bytes.Buffer)

	err = tmpl.Execute(out, data)
	if err != nil {
		panic(err)
	}

	return string(out.Bytes())
}
