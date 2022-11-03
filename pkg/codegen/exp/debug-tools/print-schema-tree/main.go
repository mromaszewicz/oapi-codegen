package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	v2 "github.com/deepmap/oapi-codegen/v2/pkg/codegen/exp"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Please specify an OpenAPI spec path")
	}

	buf, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln("Error reading spec: ", err)
	}

	loader := openapi3.NewLoader()
	spec, err := loader.LoadFromData(buf)
	if err != nil {
		log.Fatalln("Error loading spec: ", err)
	}

	root := v2.BuildSchemaTree(spec)
	printTree(0, root)
}

func printTree(indentLevel int, node *v2.PathTreeNode) {
	padding := strings.Repeat("  ", indentLevel)
	if node.Ref != "" {
		fmt.Printf(padding+"%s (Ref=\"%s\")\n", node.PathElement, node.Ref)
	} else {
		fmt.Printf(padding+"%s\n", node.PathElement)
	}
	for _, c := range node.Children {
		printTree(indentLevel+1, c)
	}
}
