package main

import (
	"fmt"
	"github.com/zjbztianya/gophercises/phone/conf"
	"github.com/zjbztianya/gophercises/phone/models"
	"strings"
)

func normalize(number string) string {
	var builder strings.Builder
	for _, r := range number {
		if r >= '0' && r <= '9' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func main() {
	if err := conf.Init(); err != nil {
		fmt.Printf("conf.Init() error(%v)\n", err)
		panic(err)
	}

	if err := models.Init(); err != nil {
		fmt.Printf("models.Init() error(%v)\n", err)
		panic(err)
	}
	defer models.Close()

	if err := models.Seed(); err != nil {
		fmt.Printf("models.Seed() error(%v)\n", err)
		panic(err)
	}

	numbers, err := models.GetPhoneNumbers()
	if err != nil {
		fmt.Printf("models.GetPhoneNumbers() error(%v)\n", err)
	}

	numberMap := make(map[string]struct{})
	for _, p := range numbers {
		number := normalize(p.Number)
		if _, ok := numberMap[number]; ok {
			if err := models.DeletePhoneNumber(p.ID); err != nil {
				fmt.Printf("models.DeletePhoneNumber() error(%v)\n", err)
				panic(err)
			}
			continue
		}
		if number != p.Number {
			p.Number = number
			if err := models.UpdatePhoneNumber(p); err != nil {
				fmt.Printf("models.UpdatePhoneNumber() error(%v)\n", err)
				panic(err)
			}
		}
		numberMap[number] = struct{}{}
	}
}
