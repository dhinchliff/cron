package main

import (
	"fmt"
	"io"
)

type CronPrinter struct{}

func (p *CronPrinter) Print(w io.Writer, cron *Cron) {
	fmt.Fprintf(w, "%-14s", "minute")
	p.printNumbers(w, cron.Minute)

	fmt.Fprintf(w, "%-14s", "hour")
	p.printNumbers(w, cron.Hour)

	fmt.Fprintf(w, "%-14s", "day of month")
	p.printNumbers(w, cron.DayOfMonth)

	fmt.Fprintf(w, "%-14s", "month")
	p.printNumbers(w, cron.Month)

	fmt.Fprintf(w, "%-14s", "day of week")
	p.printNumbers(w, cron.DayOfWeek)

	if cron.Command != "" {
		fmt.Fprintf(w, "%-14s%s\n", "command", cron.Command)
	}
}

func (p *CronPrinter) printNumbers(w io.Writer, numbers []int) {
	for _, n := range numbers {
		fmt.Fprintf(w, "%d ", n)
	}

	fmt.Fprintln(w)
}
