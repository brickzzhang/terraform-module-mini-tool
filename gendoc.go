package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
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
