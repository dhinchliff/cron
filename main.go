package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	parser := CronParser{}
	printer := CronPrinter{}

	cron, err := parser.Parse(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	printer.Print(os.Stdout, cron)
}
