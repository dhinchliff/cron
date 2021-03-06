package main

import "strings"

type CronParser struct{}

func (p *CronParser) Parse(expression string) (*Cron, error) {
	parts := strings.Split(expression, " ")

	minute, _ := p.parseMinute(parts[0])
	hour, _ := p.parseHour(parts[0])
	dayOfMonth, _ := p.parseDayOfMonth(parts[0])
	month, _ := p.parseMonth(parts[0])
	dayOfWeek, _ := p.parseDayOfWeek(parts[0])

	return &Cron{
		Minute:     minute,
		Hour:       hour,
		DayOfMonth: dayOfMonth,
		Month:      month,
		DayOfWeek:  dayOfWeek,
		Command:    parts[5],
	}, nil
}

func (p *CronParser) parseMinute(expression string) ([]int, error) {
	if expression == "*" {
		return numRange(0, 59), nil
	}

	return []int{1}, nil
}

func (p *CronParser) parseHour(expression string) ([]int, error) {
	if expression == "*" {
		return numRange(0, 23), nil
	}

	return []int{1}, nil
}

func (p *CronParser) parseDayOfMonth(expression string) ([]int, error) {
	if expression == "*" {
		return numRange(1, 31), nil
	}

	return []int{1}, nil
}

func (p *CronParser) parseMonth(expression string) ([]int, error) {
	if expression == "*" {
		return numRange(1, 12), nil
	}

	return []int{1}, nil
}

func (p *CronParser) parseDayOfWeek(expression string) ([]int, error) {
	if expression == "*" {
		return numRange(0, 6), nil
	}

	return []int{1}, nil
}
