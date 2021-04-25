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
				Minute:     []int{1, 2},
				Hour:       []int{1, 2, 3},
				DayOfMonth: []int{4, 5},
				Month:      []int{6, 7},
				DayOfWeek:  []int{3, 4},
				Command:    "/usr/bin/find",
			},
		},
		"list with duplicates": {
			expression: "1 1,2,1 4 6 3 /usr/bin/find",
			cron: &Cron{
				Minute:     []int{1},
				Hour:       []int{1, 2},
				DayOfMonth: []int{4},
				Month:      []int{6},
				DayOfWeek:  []int{3},
				Command:    "/usr/bin/find",
			},
		},
		"expression with ranges": {
			expression: "0-3,5 1-4 1-6 2-5 1-3 /usr/bin/find",
			cron: &Cron{
				Minute:     []int{0, 1, 2, 3, 5},
				Hour:       []int{1, 2, 3, 4},
				DayOfMonth: []int{1, 2, 3, 4, 5, 6},
				Month:      []int{2, 3, 4, 5},
				DayOfWeek:  []int{1, 2, 3},
				Command:    "/usr/bin/find",
			},
		},
		"expression with wrap around ranges": {
			expression: "58-2 1 1 1 1 /usr/bin/find",
			cron: &Cron{
				Minute:     []int{0, 1, 2, 58, 59},
				Hour:       []int{1},
				DayOfMonth: []int{1},
				Month:      []int{1},
				DayOfWeek:  []int{1},
				Command:    "/usr/bin/find",
			},
		},
		"expression with steps": {
			expression: "50-59/2 1 1 1 1 /usr/bin/find",
			cron: &Cron{
				Minute:     []int{50, 52, 54, 56, 58},
				Hour:       []int{1},
				DayOfMonth: []int{1},
				Month:      []int{1},
				DayOfWeek:  []int{1},
				Command:    "/usr/bin/find",
			},
		},
		"expression with steps and ranges": {
			expression: "3,5-9,50-59/2 1 1 1 1 /usr/bin/find",
			cron: &Cron{
				Minute:     []int{3, 5, 6, 7, 8, 9, 50, 52, 54, 56, 58},
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
		"expression without command": {
			expression: "1 2 3 4 5 ",
			cron: &Cron{
				Minute:     []int{1},
				Hour:       []int{2},
				DayOfMonth: []int{3},
				Month:      []int{4},
				DayOfWeek:  []int{5},
			},
		},
		"expression with value out of range": {
			expression: "1 2 3 4 100 /usr/bin/find",
			err:        fmt.Errorf("invalid field 100: unexpected value 100, expected value between 0 and 6"),
		},
		"expression with missing fields": {
			expression: "1 2 3 4 /usr/bin/find",
			err:        fmt.Errorf("invalid field /usr/bin/find: too many slashes"),
		},
		"expression with invalid numbers": {
			expression: "1 2 3 4 a /usr/bin/find",
			err:        fmt.Errorf("invalid field a: invalid value a, expected number"),
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

func TestCronParser_parseField(t *testing.T) {
	tests := map[string]struct {
		expression string
		min        int
		max        int
		out        []int
		err        error
	}{
		"steps with empty start": {
			expression: "/5",
			min:        0,
			max:        59,
			out:        []int{0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55},
		},
		"not zero based steps with empty start": {
			expression: "/2",
			min:        1,
			max:        12,
			out:        []int{1, 3, 5, 7, 9, 11},
		},
		"steps with * as start": {
			expression: "*/5",
			min:        0,
			max:        59,
			out:        []int{0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55},
		},
		"steps on single digit": {
			expression: "60/2",
			min:        0,
			max:        60,
			out:        []int{60},
		},
		"range with steps": {
			expression: "20-35/2",
			min:        0,
			max:        59,
			out:        []int{20, 22, 24, 26, 28, 30, 32, 34},
		},
		"steps with start above valid range": {
			expression: "61/2",
			min:        0,
			max:        60,
			err:        fmt.Errorf("invalid field 61/2: unexpected value 61, expected value between 0 and 60"),
		},
		"steps with start below valid range": {
			expression: "0/2",
			min:        1,
			max:        5,
			err:        fmt.Errorf("invalid field 0/2: unexpected value 0, expected value between 1 and 5"),
		},
		"steps with invalid start number": {
			expression: "a/2",
			min:        0,
			max:        60,
			err:        fmt.Errorf("invalid field a/2: invalid value a, expected number"),
		},
		"steps with invalid increment number": {
			expression: "1/x",
			min:        1,
			max:        5,
			err:        fmt.Errorf("invalid field 1/x: invalid value x, expected number"),
		},
		"steps with empty step number": {
			expression: "/",
			min:        1,
			max:        12,
			err:        fmt.Errorf("invalid field /: empty value, expected number"),
		},
		"steps with zero step number": {
			expression: "1/0",
			min:        1,
			max:        12,
			err:        fmt.Errorf("invalid field 1/0: unexpected step 0, expected value greater than zero"),
		},
		"range wraps valid range": {
			expression: "0-60",
			min:        1,
			max:        5,
			err:        fmt.Errorf("invalid field 0-60: unexpected value 0, expected value between 1 and 5"),
		},
		"range overlaps start of valid range": {
			expression: "0-1",
			min:        1,
			max:        5,
			err:        fmt.Errorf("invalid field 0-1: unexpected value 0, expected value between 1 and 5"),
		},
		"range overlaps end of valid range": {
			expression: "5-60",
			min:        1,
			max:        5,
			err:        fmt.Errorf("invalid field 5-60: unexpected value 60, expected value between 1 and 5"),
		},
		"range with invalid start number": {
			expression: "!-60",
			min:        1,
			max:        5,
			err:        fmt.Errorf("invalid field !-60: invalid value !, expected number"),
		},
		"range with invalid end number": {
			expression: "1-!",
			min:        1,
			max:        5,
			err:        fmt.Errorf("invalid field 1-!: invalid value !, expected number"),
		},
		"range with empty start number": {
			expression: "-60",
			min:        1,
			max:        5,
			err:        fmt.Errorf("invalid field -60: empty value, expected number"),
		},
		"range with empty end number": {
			expression: "1-",
			min:        1,
			max:        5,
			err:        fmt.Errorf("invalid field 1-: empty value, expected number"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := &CronParser{}
			out, err := p.parseField(tc.expression, tc.min, tc.max)

			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.out, out)
		})
	}
}
