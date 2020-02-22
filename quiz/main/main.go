package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type Problem struct {
	text   string
	answer string
}
type Problems []Problem

func loadProblems(csvFileName *string) Problems {
	csvFile, err := os.Open(*csvFileName)
	if err != nil {
		log.Fatalf("open %s fail\n", *csvFileName)
	}
	reader := csv.NewReader(csvFile)
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("read problems fail")
	}

	problems := make(Problems, 0, len(lines))
	for _, s := range lines {
		problems = append(problems, Problem{text: s[0], answer: s[1]})
	}

	return problems
}

func starQuiz(problem Problem, correct chan int) {
	var answer string
	fmt.Println(problem.text + "?")
	fmt.Scanln(&answer)
	if answer == problem.answer {
		fmt.Println("correct!")
		correct <- 1
	} else {
		fmt.Println("wrong!")
		correct <- 0
	}
}

func run(problems Problems) (correct int) {
	quizTimer := time.NewTimer(30 * time.Second)
	correctChan := make(chan int, 1)
	for _, problem := range problems {
		go starQuiz(problem, correctChan)
		select {
		case <-quizTimer.C:
			fmt.Println("timeout!")
			return correct
		case result := <-correctChan:
			correct += result
		}
		quizTimer.Reset(30 * time.Second)
	}
	fmt.Println("quiz end!")
	return
}

func main() {
	csvFileName := flag.String("csv", "./quiz/problems.csv", "a csv file in the format of 'question,answer'")
	flag.Parse()

	fmt.Printf("right answers number:%d\n", run(loadProblems(csvFileName)))
}
