package filesytem

import (
	"os"
	"path/filepath"
)

const PERMISSIONS = 0755

var workingDirectory string

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	pie(err)
	workingDirectory = dir + "/models"
}

func SetWorkingDirectory(dir string) {
	//TODO: Check for the existence of the directory
	if dir != "" {
		workingDirectory = dir
	}
}

func pie(e error) {
	if e != nil {
		panic(e)
	}
}

func createFile(fileName string, content string) {
	var e error
	f, e := os.Create(fileName)
	pie(e)
	//	fmt.Println("[YAO]creating file " + fileName)
	_, e = f.WriteString(content)
	pie(e)

	return
}

func CreateModelsFolder(output string) error {
	return os.Mkdir(workingDirectory, PERMISSIONS)
}

// Creates a folder in MODELSDIR with the given name
func CreateSchemaFolder(name string) error {
	return os.Mkdir(workingDirectory+"/"+name, PERMISSIONS)
}

// Creates the models/schemaname/query.go file or panics
func CreateQueryFile(schema string, content string) {
	createFile(workingDirectory+"/"+schema+"/query.go", content)
}

func CreateModelFile(schema string, name string, content string) {
	var fileName string
	if schema == "" { // mysql doesn't have schema as a namespace concept as postgres does
		os.Mkdir(workingDirectory+"/"+name, PERMISSIONS)
		fileName = workingDirectory + "/" + name + "/" + name + ".go"
	} else {
		os.Mkdir(workingDirectory+"/"+schema+"/"+name, PERMISSIONS)
		fileName = workingDirectory + "/" + schema + "/" + name + "/" + name + ".go"
	}
	createFile(fileName, content)
}

// Returns true if file exists, false otherwhise
func FileExists(fileName string) bool {
	_, ferr := os.Lstat(fileName)
	return ferr == nil
}
