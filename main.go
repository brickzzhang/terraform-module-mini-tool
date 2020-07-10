package main

import (
	"fmt"
	"strings"
)

func main() {
	for {
	loop:
		fmt.Printf("%s", hint)
		var option string
		fmt.Scanf("%s\n", &option)

		switch strings.ToLower(option) {
		case "0":
			fmt.Printf("%s", moduleInitHint)
			var isContinue string
			fmt.Scanf("%s", &isContinue)
			switch strings.ToLower(isContinue) {
			case "y":
				var path string
				fmt.Println("path of your module: ")
				fmt.Scanf("%s", &path)
				moduleInit(path)
				fmt.Println("Module initiation done!")
			case "n":
				goto loop
			default:
				fmt.Println("Invalid type, retype please!")
			}
		case "1":
			var path string
			fmt.Println("path of your module :")
			fmt.Scanf("%s", &path)
			demoReadmeGenerate(path)
			fmt.Println("DEMO-README.md generation done!")
		case "q":
			fmt.Println("See you next time, goodbye!")
			return
		default:
			fmt.Println("Invalid type, retype please!")
		}
	}
}
