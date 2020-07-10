package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
)

func cmdRun(command string, argv ... string)  {
	cmd := exec.Command(command, argv...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("%s run failed with %s", command, err)
	}
}

func moduleInit(path string) {
	files := map[string]string{
		"variables.tf": variables,
		"outputs.tf":outputs,
		"version.tf": version,
		"README.md": readme,
		".gitignore": gitignore,
		"main.tf": mainTem,
		"LICENSE": license,
	}
	for k, v := range files {
		if err := ioutil.WriteFile(path+"/"+k, []byte(v), 0644); err != nil {
			log.Fatalf("%+v", err)
		}
	}
	cmdRun("mkdir", "-p", path+"/"+"examples")
}