package generate

import (
	"bytes"
	"database/sql"
	"fmt"
	fs "github.com/alfonsodev/yao/filesystem"
	_ "github.com/lib/pq"
	"os/exec"
	"strings"
	"text/template"
)

var db *sql.DB

type FieldInfo struct {
	Name     string
	Nullable string
	Datatype string
	KeyInfo  string
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

// Genreate one modle file per table
func Generate(schemaname string) {
	//TODO: make a flag to change folder name
	err := fs.CreateModelsFolder()
	if err != nil {
		fmt.Println("[YAO] ./models folder already exist")
	}

	schemas := yao.GetSchemas()
	if len(schemas) > 0 {
		for _, v := range schemas {
			generateModelFromSchema(v)
		}
		return
	}
}

func generateModelFromSchema(schemaname string) {
	info := yao.GetInformationSchema(schemaname)
	err := fs.CreateSchemaFolder(schemaname)
	if err != nil {
		fmt.Println("[YAO] can't create ./models/" + schemaname + " , folder already exist?")
	}

	for k, v := range info {
		fs.CreateModelFile(schemaname, strings.ToLower(k), PrintModel(k, v))
		fmt.Println("./models/" + schemaname + "/" + strings.ToLower(k) + ".go")
		cmd := exec.Command("go", "fmt", "./models/"+schemaname+"/"+UcFirst(k)+".go")
		err := cmd.Run()
		if err != nil {
			fmt.Println(string(err.Error()))
		}

	}
}

func PrintModel(name string, fields []FieldInfo) string {
	type TemplateData struct {
		Name         string
		StructFields string
		ScanFields   string
		Schema       string
		Table        string
		Keys         string
		Placeholders string
		SaveFields   string
		JsonFields   string
	}
	data := TemplateData{
		Name:   UcFirst(name),
		Schema: "usermanager",
		Table:  strings.ToLower(name),
	}
	// , "// struct fields\n", "//scan fields \n", "usermanager", "users", "keys", "place"}

	for i, v := range fields {
		data.StructFields += "   " + UcFirst(v.Name) + " " + v.Datatype + "\n"
		data.ScanFields += "        &u." + UcFirst(v.Name) + ", \n"
		data.SaveFields += "          util.GetValue(obj." + UcFirst(v.Name) + "), \n"
		if v.Name == "json" {
			data.JsonFields += "obj.Json.String = \"{}\" \n"
		}

		data.Placeholders += fmt.Sprintf(", $%v", i)
		fmt.Println(v.Name + " : keyinfo : " + v.KeyInfo + "\n")
		data.Keys += ", " + strings.ToLower(v.Name)
	}
	data.Placeholders = data.Placeholders[1:]
	data.Keys = data.Keys[1:]

	tmpl, err := template.ParseFiles("./template/model.tmpl")
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
