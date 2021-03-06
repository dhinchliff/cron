package main

import (
	"os"
)

func main() {
	args := os.Args[1:]
	parser := CronParser{}
	printer := CronPrinter{}

	cron, _ := parser.Parse(args[0])

	printer.Print(os.Stdout, cron)
}
