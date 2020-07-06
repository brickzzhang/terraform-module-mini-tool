package main

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"

	"regexp"
)

var templateStr =
`
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
	if err :=ioutil.WriteFile(desCfg, data, 0644); err != nil {
		panic(err)
	}
}

func outputPreProcess(cfg string, desCfg string) {
	data, _ := ioutil.ReadFile(cfg)
	regex :=  regexp.MustCompile(`\s*value\s*=(.*)\n`)
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

func inputGenReadmeStr(data map[string]interface{}) string {
	var inputStr = "| Name | Description | Type | Default | Required |\n|------|-------------|:----:|:-----:|:-----:|\n"
	var inputLine = "| %s | %s | %s | %s | %s \n"

	for _, v := range data["variable"].([]map[string]interface{}) {
		for kk, vv := range v {
			for _, vvv := range vv.([]map[string]interface{}) {
				name := kk
				description := ""
				if vvv["description"] != nil {
					description = vvv["description"].(string)
				}
				typeInterface := reflect.ValueOf(vvv["type"])
				defaultInterface := reflect.ValueOf(vvv["default"])
				required := "no"
				if vvv["required"] != nil {
					required = "yes"
				}
				str := fmt.Sprintf(inputLine, name, description, typeInterface, defaultInterface, required)
				inputStr += str
			}
		}
	}

	return inputStr
}

func outputGenReadmeStr(data map[string]interface{}) string {
	var outputStr = "| Name | Description |\n|------|-------------|\n"
	var outputLine = "| %s | %s |\n"

	for _, v := range data["output"].([]map[string]interface{}) {
		for kk, vv := range v {
			for _, vvv := range vv.([]map[string]interface{}) {
				name := kk
				description := ""
				if vvv["description"] != nil {
					description = vvv["description"].(string)
				}
				str := fmt.Sprintf(outputLine, name, description)
				outputStr += str
			}
		}
	}

	return outputStr
}

func generateReadmeStr(config string, desConfig string, preProcessFun func(string, string), genStrFun func(map[string]interface{}) string ) (readmeStr string) {
	preProcessFun(config, desConfig)
	data := parse(desConfig)
	readmeStr = genStrFun(data)
	return
}

func cleanTmp(file string) {
	cmd := exec.Command("rm", file)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
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
	if checkFileExist(path+"/"+inputCfg) {
		desInputCfg := "tmp-"+inputCfg
		inputStr = generateReadmeStr(path+"/"+inputCfg, desInputCfg, inputPreProcess, inputGenReadmeStr)
		cleanTmp(desInputCfg)
	}

	var outputCfg = "outputs.tf"
	if checkFileExist(path+"/"+outputCfg) {
		desOutputCfg := "tmp-"+outputCfg
		outputStr = generateReadmeStr(path+"/"+outputCfg, desOutputCfg, outputPreProcess, outputGenReadmeStr)
		cleanTmp(desOutputCfg)
	}

	readmeFile := "DEMO-README.md"
	readmeStr := fmt.Sprintf(templateStr, inputStr, outputStr)
	if err := ioutil.WriteFile(readmeFile, []byte(readmeStr), 0644); err != nil {
		panic(err)
	}
}