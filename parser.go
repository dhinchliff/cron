package main

import (
	"strconv"
	"strings"
)

type CronParser struct{}

func (p *CronParser) Parse(expression string) (*Cron, error) {
	expression = strings.TrimSpace(expression)
	parts := strings.Split(expression, " ")

	minute, _ := p.parseMinute(parts[0])
	hour, _ := p.parseHour(parts[1])
	dayOfMonth, _ := p.parseDayOfMonth(parts[2])
	month, _ := p.parseMonth(parts[3])
	dayOfWeek, _ := p.parseDayOfWeek(parts[4])

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
	return p.parseExpression(expression, 0, 59)
}

func (p *CronParser) parseHour(expression string) ([]int, error) {
	return p.parseExpression(expression, 0, 23)
}

func (p *CronParser) parseDayOfMonth(expression string) ([]int, error) {
	return p.parseExpression(expression, 1, 31)
}

func (p *CronParser) parseMonth(expression string) ([]int, error) {
	return p.parseExpression(expression, 1, 12)
}

func (p *CronParser) parseDayOfWeek(expression string) ([]int, error) {
	return p.parseExpression(expression, 0, 6)
}

func (p *CronParser) parseExpression(expression string, min int, max int) ([]int, error) {
	if expression == "*" {
		return numRange(min, max), nil
	}

	i, err := strconv.Atoi(expression)
	if err != nil {
		return nil, err
	}

	return []int{i}, nil
}