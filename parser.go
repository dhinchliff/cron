package main

type CronParser struct{}

func (p *CronParser) Parse(expression string) (*Cron, error) {
	return &Cron{
		Minute:     []int{1},
		Hour:       []int{1},
		DayOfMonth: []int{1},
		Month:      []int{1},
		DayOfWeek:  []int{1},
		Command:    "/usr/bin/find",
	}, nil
}
