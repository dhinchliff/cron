package main

import (
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
