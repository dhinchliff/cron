package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCronParser_Parse(t *testing.T) {
	tests := map[string]struct {
		expression string
		cron       *Cron
		err        error
	}{
		"basic expression": {
			expression: "1 2 3 4 5 /usr/bin/find",
			cron: &Cron{
				Minute:     []int{1},
				Hour:       []int{2},
				DayOfMonth: []int{3},
				Month:      []int{4},
				DayOfWeek:  []int{5},
				Command:    "/usr/bin/find",
			},
		},
		"all stars": {
			expression: "* * * * * /usr/bin/find",
			cron: &Cron{
				Minute:     numRange(0, 59),
				Hour:       numRange(0, 23),
				DayOfMonth: numRange(1, 31),
				Month:      numRange(1, 12),
				DayOfWeek:  numRange(0, 6),
				Command:    "/usr/bin/find",
			},
		},
		"lists": {
			expression: "1,2 1,2,3 4,5 6,7 3,4 /usr/bin/find",
			cron: &Cron{
				Minute:     []int{1,2},
				Hour:       []int{1,2,3},
				DayOfMonth: []int{4,5},
				Month:      []int{6,7},
				DayOfWeek:  []int{3,4},
				Command:    "/usr/bin/find",
			},
		},
		"list with duplicates": {
			expression: "1 1,2,1 4 6 3 /usr/bin/find",
			cron: &Cron{
				Minute:     []int{1},
				Hour:       []int{1,2},
				DayOfMonth: []int{4},
				Month:      []int{6},
				DayOfWeek:  []int{3},
				Command:    "/usr/bin/find",
			},
		},
		"expression with extra whitespace on ends": {
			expression: " 1 2 3 4 5 /usr/bin/find ",
			cron: &Cron{
				Minute:     []int{1},
				Hour:       []int{2},
				DayOfMonth: []int{3},
				Month:      []int{4},
				DayOfWeek:  []int{5},
				Command:    "/usr/bin/find",
			},
		},
		"expression with value out of range": {
			expression: "1 2 3 4 100 /usr/bin/find",
			err: fmt.Errorf("unexpected value 100, expected value between 0 and 6"),
		},
		"expression with missing fields": {
			expression: "1 2 3 4 /usr/bin/find",
			err: fmt.Errorf("invalid expression [1 2 3 4 /usr/bin/find]"),
		},
		"expression with too many fields": {
			expression: "1 2 3 4 5 6 /usr/bin/find",
			err: fmt.Errorf("invalid expression [1 2 3 4 5 6 /usr/bin/find]"),
		},
		"expression with invalid numbers fields": {
			expression: "1 2 3 4 a /usr/bin/find",
			err: fmt.Errorf("invalid value a, expected number"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := &CronParser{}
			cron, err := p.Parse(tc.expression)

			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.cron, cron)
		})
	}
}
