package main

import (
	"flag"
	"fmt"
	"github.com/zjbztianya/gophercises/link"
	"os"
)

var htmlFileName *string

func init() {
	htmlFileName = flag.String("html", "link/ex4.html", "a html file")
	flag.Parse()
}

func main() {
	file, err := os.Open(*htmlFileName)
	if err != nil {
		panic(err)
	}

	links, err := link.Parse(file)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", links)
}
