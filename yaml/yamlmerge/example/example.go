package main

import (
	"fmt"
	"os"

	"github.com/zusux/gokit/yaml/yamlmerge"
)

func main() {
	vars, _ := os.ReadFile("./vars.yaml")
	tmpl, _ := os.ReadFile("./template.yaml")

	out, err := yamlmerge.MergeWithComments(string(tmpl), vars)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}
