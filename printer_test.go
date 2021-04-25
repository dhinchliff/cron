package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCronPrinter_Print(t *testing.T) {
	tests := map[string]struct {
		cron *Cron
		out  []string
	}{
		"basic cron": {
			cron: &Cron{
				Minute:     []int{1},
				Hour:       []int{1},
				DayOfMonth: []int{1},
				Month:      []int{1},
				DayOfWeek:  []int{1},
				Command:    "/usr/bin/find",
			},
			out: []string{
				"minute        1 ",
				"hour          1 ",
				"day of month  1 ",
				"month         1 ",
				"day of week   1 ",
				"command       /usr/bin/find",
			},
		},
		"cron without command": {
			cron: &Cron{
				Minute:     []int{1},
				Hour:       []int{1},
				DayOfMonth: []int{1},
				Month:      []int{1},
				DayOfWeek:  []int{1},
				Command:    "",
			},
			out: []string{
				"minute        1 ",
				"hour          1 ",
				"day of month  1 ",
				"month         1 ",
				"day of week   1 ",
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			p := &CronPrinter{}
			p.Print(&buf, tc.cron)

			linesOut := strings.Split(string(buf.Bytes()), "\n")

			assert.Equal(t, tc.out, linesOut[:len(linesOut)-1])
			assert.Equal(t, "", linesOut[len(linesOut)-1])
		})
	}
}
