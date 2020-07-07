package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/spf13/viper"
)

var templateStr = `
# TencentCloud [MODIFY !!!] Module for Terraform

## terraform-tencentcloud-vpc

A terraform module used to create TencentCloud VPC, subnet and route entry[MODIFY !!!].

The following resources are included.

* [MODIFY !!!][VPC](https://www.terraform.io/docs/providers/tencentcloud/r/vpc.html)
* [MODIFY !!!][VPC Subnet](https://www.terraform.io/docs/providers/tencentcloud/r/subnet.html)
* [MODIFY !!!][VPC Route Entry](https://www.terraform.io/docs/providers/tencentcloud/r/route_table_entry.html)

## Usage


## Conditional Creation

This module can create VPC and VPC Subnet[MODIFY !!!].

## Inputs

%s

## Outputs

%s

## Authors

Created and maintained by [TencentCloud](https://github.com/terraform-providers/terraform-provider-tencentcloud)

## License

Mozilla Public License Version 2.0.
See LICENSE for full details.
`

var keywordList = []string{"number", "null", "string", "map(string)", "list(string)", "bool", "true", "false"}

func inputPreProcess(cfg string, desCfg string) {
	data, _ := ioutil.ReadFile(cfg)
	for _, keyword := range keywordList {
		data = bytes.Replace(data, []byte("= "+keyword), []byte("= \""+keyword+"\""), -1)
	}
	if err := ioutil.WriteFile(desCfg, data, 0644); err != nil {
		panic(err)
	}
}

func outputPreProcess(cfg string, desCfg string) {
	data, _ := ioutil.ReadFile(cfg)
	regex := regexp.MustCompile(`\s*value\s*=(.*)\n`)
	linesRegex := regex.FindAllStringSubmatch(string(data), -1)
	var outStr []string
	for _, item := range linesRegex {
		outStr = append(outStr, item[1])
	}

	for _, item := range outStr {
		data = bytes.Replace(data, []byte(item), []byte(" \"\""), -1)
	}
	if err := ioutil.WriteFile(desCfg, data, 0644); err != nil {
		panic(err)
	}
}

func parse(cfg string) map[string]interface{} {
	viper.SetConfigName(cfg)
	viper.AddConfigPath(".")
	viper.SetConfigType("hcl")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("ReadInConfig error: %s\n", err)
		panic(err)
	}

	return viper.AllSettings()
}

func inputGenReadmeStr(data jsonObj) string {
	var inputStr = "| Name | Description | Type | Default | Required |\n|------|-------------|:----:|:-----:|:-----:|\n"
	var inputLine = "| %s | %s | %s | %s | %s \n"

	for k, v := range data["variable"].(jsonObj) {
		of := reflect.ValueOf(v)
		typeName := of.Kind()
		if typeName == reflect.Map {
			keys := of.MapKeys()
			description, required, typeInterface := "", "no", "string"
			var defaultInterface interface{}
			for key := range keys {
				if keys[key].String() == "description" {
					description = of.MapIndex(keys[key]).Elem().String()
				}
				if keys[key].String() == "required" {
					requireds := of.MapIndex(keys[key])
					if !requireds.IsNil() {
						required = "yes"
					}
				}

				if keys[key].String() == "type" {
					spkind := of.MapIndex(keys[key])
					resultStr := spkind.Elem().String()
					resultStr = strings.ReplaceAll(resultStr, "${", "")
					resultStr = strings.ReplaceAll(resultStr, "}", "")
					typeInterface = resultStr
					if strings.Index(resultStr, "object") == 0 {
						typeInterface = "object"
					}
				}

				if keys[key].String() == "default" {
					ele := of.MapIndex(keys[key]).Elem()
					if ele.Kind() == reflect.String {
						result := strings.ReplaceAll(ele.String(), "${", "")
						result = strings.ReplaceAll(result, "}", "")
						defaultInterface = result
					} else if ele.Kind() == reflect.Int {
						defaultInterface = ele.Int()
					} else if ele.Kind() == reflect.Bool {
						defaultInterface = ele.Bool()

					} else if ele.Kind() == reflect.Map {
						if ele.IsValid() {
							defaultInterface, _ = json.Marshal(map[string]interface{}{})
						}
					} else if ele.Kind() == reflect.Slice {
						if ele.IsValid() {
							defaultInterface, _ = json.Marshal([]interface{}{})
						}
					} else {
						defaultInterface = "xxxxxxxx"
					}
				}
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

func generateReadmeStr(config string, desConfig string, preProcessFun func(string, string), genStrFun func(jsonObj) string) (readmeStr string) {
	//preProcessFun(config, desConfig)
	//data := parse(desConfig)'
	var bytes []byte
	var err error

	bytes, err = ioutil.ReadFile(config)
	if err != nil {
		fmt.Errorf("Failed to read file: %s\n", err)
	}

	var content interface{}
	content, err = getHclJSON(bytes, config)
	if err != nil {
		fmt.Errorf("Failed to read file: %s\n", err)
		return
	}

	readmeStr = genStrFun(content.(jsonObj))
	return
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
		/*	if v, ok := obj["variable"]; ok {
			mmp := map[string]interface{}{}
			values := v.(jsonObj)
			for kp, kv := range values {
				mmp[kp] = kv
			}
			return v, nil
		}*/
		return obj, nil
	}

	return nil, nil
}

func checkFileExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func main() {
	fmt.Println("Please input your variables.tf path, e.g. /Users/brick/workspace/terraform/terraform-tencentcloud-modules/terraform-tencentcloud-clb/")
	var path, inputStr, outputStr string
	fmt.Scanf("%s\n", &path)

	var inputCfg = "variables.tf"
	if checkFileExist(path + "/" + inputCfg) {
		desInputCfg := "tmp-" + inputCfg
		inputStr = generateReadmeStr(path+"/"+inputCfg, desInputCfg, inputPreProcess, inputGenReadmeStr)
	}

	var outputCfg = "outputs.tf"
	if checkFileExist(path + "/" + outputCfg) {
		desOutputCfg := "tmp-" + outputCfg
		outputStr = generateReadmeStr(path+"/"+outputCfg, desOutputCfg, outputPreProcess, outputGenReadmeStr)
	}

	readmeFile := "DEMO-README.md"
	readmeStr := fmt.Sprintf(templateStr, inputStr, outputStr)
	if err := ioutil.WriteFile(readmeFile, []byte(readmeStr), 0644); err != nil {
		panic(err)
	}
}
