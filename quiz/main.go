package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
)

type Problem struct {
	question string
	answer   string
}
type Problems []Problem

func main() {
	csvFileName := flag.String("csv", "./quiz/problems.csv", "a csv file in the format of 'question,answer'")
	flag.Parse()

	csvFile, err := os.Open(*csvFileName)
	if err != nil {
		log.Fatalf("open %s fail\n", *csvFileName)
	}
	reader := csv.NewReader(csvFile)
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("read problems fail")
	}

	problems := make(Problems, len(lines))
	for _, s := range lines {
		problems = append(problems, Problem{question: s[0], answer: s[1]})
		fmt.Println(s[0],s[1])
	}
}
