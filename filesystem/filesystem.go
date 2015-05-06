package filesytem

import (
	"fmt"
	"os"

)

const MODELSDIR = "./models" //TODO: make it configurable by flag
const PERMISSIONS = 0777

func CreateModelsFolder() error {
	return os.Mkdir(MODELSDIR, PERMISSIONS)
}

func CreateSchemaFolder(schema string) error {
	return os.Mkdir(MODELSDIR + "/" + schema, PERMISSIONS)
}

func CreateModelFile(schema string, name string, content string) error {
	var fileName string 
	if schema == "" {
		fileName = MODELSDIR + "/" + name + ".go"
	} else {
		fileName = MODELSDIR + "/"+schema + "/" + name + ".go"
	}
	// _, ferr := os.Lstat(fileName)

	// if ferr == nil {
	// 	fmt.Println("[YAO] Model file " + fileName + " already exist, use --force to overwrite")
	// 	return ferr
	// }

	fmt.Println("[YAO]creating file " + MODELSDIR + "/"+schema + "/" + ".go")
	f, e := os.Create(fileName)
	if e != nil {
		panic(e)
	}
	_, err := f.WriteString(content)
	if err != nil {
		panic(err)
	}

	return e
}

func CreateFileFromString(fileName string, content string) {

}

func main() {

}
