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
		"expression with ranges": {
			expression: "0-3,5 1-4 1-6 2-5 1-3 /usr/bin/find",
			cron: &Cron{
				Minute:     []int{0,1,2,3,5},
				Hour:       []int{1,2,3,4},
				DayOfMonth: []int{1,2,3,4,5,6},
				Month:      []int{2,3,4,5},
				DayOfWeek:  []int{1,2,3},
				Command:    "/usr/bin/find",
			},
		},
		"expression with wrap around ranges": {
			expression: "58-2 1 1 1 1 /usr/bin/find",
			cron: &Cron{
				Minute:     []int{0,1,2,58,59},
				Hour:       []int{1},
				DayOfMonth: []int{1},
				Month:      []int{1},
				DayOfWeek:  []int{1},
				Command:    "/usr/bin/find",
			},
		},
		"expression with steps": {
			expression: "50/2 1 1 1 1 /usr/bin/find",
			cron: &Cron{
				Minute:     []int{50,52,54,56,58},
				Hour:       []int{1},
				DayOfMonth: []int{1},
				Month:      []int{1},
				DayOfWeek:  []int{1},
				Command:    "/usr/bin/find",
			},
		},
		"expression with steps and ranges": {
			expression: "3,5-9,50/2 1 1 1 1 /usr/bin/find",
			cron: &Cron{
				Minute:     []int{3,5,6,7,8,9,50,52,54,56,58},
				Hour:       []int{1},
				DayOfMonth: []int{1},
				Month:      []int{1},
				DayOfWeek:  []int{1},
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
		"expression with command with spaces": {
			expression: "1 2 3 4 5 /usr/bin/find file",
			cron: &Cron{
				Minute:     []int{1},
				Hour:       []int{2},
				DayOfMonth: []int{3},
				Month:      []int{4},
				DayOfWeek:  []int{5},
				Command:    "/usr/bin/find file",
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

func TestCronParser_parseExpression(t *testing.T) {
	tests := map[string]struct {
		expression string
		min 	   int
		max 	   int
		out        []int
		err        error
	}{
		"steps with start on limit": {
			expression: "60/2",
			min: 0,
			max: 60,
			out: []int{60},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := &CronParser{}
			out, err := p.parseExpression(tc.expression, tc.min, tc.max)

			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.out, out)
		})
	}
}
