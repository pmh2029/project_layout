package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

type JsonSchemaValidator struct {
	basePath string
	schemas  map[string]*gojsonschema.Schema // load all schemas file from the base path and store them in the schemas map
}

func NewJsonSchemaValidator() (*JsonSchemaValidator, error) {
	pwdPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	pwdPath = (filepath.Dir(pwdPath))                // need to make sure path is directly project location
	pwdPath = strings.ReplaceAll(pwdPath, "\\", "/") // replace "\" by "/" in windows os

	validator := &JsonSchemaValidator{
		basePath: pwdPath,
		schemas:  make(map[string]*gojsonschema.Schema),
	}
	err = validator.loadDirSchemas("")
	if err != nil {
		return nil, err
	}

	return validator, nil
}

func (validator *JsonSchemaValidator) loadDirSchemas(path string) error {
	err := filepath.Walk(validator.basePath+path, func(path string, f os.FileInfo, _ error) error {
		if f.IsDir() {
			return nil
		}
		if !strings.HasSuffix(f.Name(), ".json") {
			return nil
		}

		path = "file://" + strings.ReplaceAll(path, `\`, "/") // replace "\" by "/" in windows os
		schemaLoader := gojsonschema.NewReferenceLoader(path)
		schema, err := gojsonschema.NewSchema(schemaLoader)
		if err != nil {
			return err
		}
		validator.schemas[f.Name()] = schema
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
func (validator *JsonSchemaValidator) Validate(
	schemaFile string,
	data interface{},
) (*gojsonschema.Result, error) {
	schema, schemaExists := validator.schemas[schemaFile]
	if !schemaExists {
		return nil, fmt.Errorf("The schema '%v' was not found for json validation", schemaFile)
	}

	dataLoader := gojsonschema.NewGoLoader(data)

	result, err := schema.Validate(dataLoader)
	if err != nil {
		return nil, err
	}
	if len(result.Errors()) == 0 {
		return nil, nil
	}

	return result, nil
}
