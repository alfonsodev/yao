package filesytem

import "os"

const MODELSDIR = "./models" //TODO: make it configurable by flag
const PERMISSIONS = 0777

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

func CreateModelsFolder() error {
	return os.Mkdir(MODELSDIR, PERMISSIONS)
}

// Creates a folder in MODELSDIR with the given name
func CreateSchemaFolder(name string) error {
	return os.Mkdir(MODELSDIR+"/"+name, PERMISSIONS)
}

// Creates the models/schemaname/query.go file or panics
func CreateQueryFile(schema string, content string) {
	createFile(MODELSDIR+"/"+schema+"/query.go", content)
}

func CreateModelFile(schema string, name string, content string) {
	var fileName string
	if schema == "" {
		fileName = MODELSDIR + "/" + name + ".go"
	} else {
		fileName = MODELSDIR + "/" + schema + "/" + name + ".go"
	}
	createFile(fileName, content)
}

// Returns true if file exists, false otherwhise
func FileExists(fileName string) bool {
	_, ferr := os.Lstat(fileName)
	return ferr == nil
}
