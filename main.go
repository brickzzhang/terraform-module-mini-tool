package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	variablesFile := flag.String("variables", "variables.tf", "variables.tf file path")
	outputsFile := flag.String("outputs", "outputs.tf", "outputs.tf file path")
	readmeFile := flag.String("readme", "DEMO-README.md", "README.md file path")

	flag.Parse()

	inputStr, err := generateReadmeStr(*variablesFile, inputGenReadmeStr)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	outputStr, err := generateReadmeStr(*outputsFile, outputGenReadmeStr)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	readmeStr := fmt.Sprintf(templateStr, inputStr, outputStr)
	if err := ioutil.WriteFile(*readmeFile, []byte(readmeStr), 0644); err != nil {
		log.Fatalf("%+v", err)
	}
}
