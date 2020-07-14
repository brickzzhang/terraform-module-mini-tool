package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func inputGenReadmeStr(data jsonObj) string {
	var inputStr = "| Name | Description | Type | Default | Required |\n|------|-------------|:----:|:-----:|:-----:|\n"
	var inputLine = "| %s | %s | %s | %s | %s \n"

	for k, v := range data["variable"].(jsonObj) {
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Map {
			description, required, typeInterface := "", "yes", "string"
			var defaultInterface interface{}
			for _, e := range val.MapKeys() {
				switch e.String() {
				case "description":
					description = val.MapIndex(e).Elem().String()
				case "type":
					typeInterface = val.MapIndex(e).Elem().String()
					typeInterface = strings.ReplaceAll(typeInterface, "${", "")
					typeInterface = strings.ReplaceAll(typeInterface, "}", "")
					typeInterface = strings.ReplaceAll(typeInterface, "\n", "<br>")
					if strings.Index(typeInterface, "object") == 0 {
						typeInterface = "object"
					}
					if strings.Index(typeInterface, "list") == 0 {
						typeInterface = "list"
					}
					if strings.Index(typeInterface, "map") == 0 {
						typeInterface = "map"
					}
				case "default":
					required = "no"
					ele := val.MapIndex(e).Elem()
					switch ele.Kind() {
					case reflect.String:
						defaultInterface = strings.ReplaceAll(ele.String(), "${", "")
						defaultInterface = strings.ReplaceAll(defaultInterface.(string), "}", "")
						if strings.ToLower(defaultInterface.(string)) == "true" || strings.ToLower(defaultInterface.(string)) == "false" {
							typeInterface = "bool"
						}
					case reflect.Int:
						defaultInterface = ele.Int()
						typeInterface = "number"
					case reflect.Bool:
						defaultInterface = ele.Bool()
						typeInterface = "bool"
					case reflect.Map:
						if ele.IsValid() {
							defaultInterface, _ = json.Marshal(map[string]interface{}{})
						}
						typeInterface = "map"
					case reflect.Slice:
						if ele.IsValid() {
							defaultInterface, _ = json.Marshal([]interface{}{})
						}
						typeInterface = "list"
					default:
						defaultInterface = "xxxxxxxx"
					}
				}
			}
			if defaultInterface == nil {
				defaultInterface = ""
			}
			str := fmt.Sprintf(inputLine, k, description, typeInterface, defaultInterface, required)
			inputStr += str
		}
	}

	return inputStr
}

func outputGenReadmeStr(data jsonObj) string {
	var outputStr = "| Name | Description |\n|------|-------------|\n"
	var outputLine = "| %s | %s |\n"

	for k, v := range data["output"].(jsonObj) {
		of := reflect.ValueOf(v)
		typeName := of.Kind()
		if typeName == reflect.Map {
			keys := of.MapKeys()
			description := ""

			for key := range keys {
				if keys[key].String() == "description" {
					description = of.MapIndex(keys[key]).Elem().String()
				}
			}
			str := fmt.Sprintf(outputLine, k, description)
			outputStr += str
		}
	}

	return outputStr
}

func getHclJSON(bytes []byte, filename string) (interface{}, error) {
	file, diags := hclsyntax.ParseConfig(bytes, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, diags
	}
	obj, err := convertFile(file)
	if err != nil {
		return nil, nil
	}

	if len(obj) > 0 {
		return obj, nil
	}

	return nil, nil
}

func generateReadmeStr(config string, genStrFun func(jsonObj) string) (readmeStr string, err error) {
	var data []byte

	data, err = ioutil.ReadFile(config)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	var content interface{}
	content, err = getHclJSON(data, config)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return genStrFun(content.(jsonObj)), nil
}

func demoReadmeGenerate(path string) {
	flag.Parse()

	inputStr, err := generateReadmeStr(path+"/"+"variables.tf", inputGenReadmeStr)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	outputStr, err := generateReadmeStr(path+"/"+"outputs.tf", outputGenReadmeStr)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	readmeStr := fmt.Sprintf(templateStr, inputStr, outputStr)
	if err := ioutil.WriteFile(path+"/"+"DEMO-README.md", []byte(readmeStr), 0644); err != nil {
		log.Fatalf("%+v", err)
	}
}
