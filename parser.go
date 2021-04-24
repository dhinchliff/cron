package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type CronParser struct{}

func (p *CronParser) Parse(expression string) (*Cron, error) {
	expression = strings.TrimSpace(expression)
	parts := strings.Split(expression, " ")

	if len(parts) < 6 {
		return nil, fmt.Errorf("invalid expression [%s]", expression)
	}

	minute, err := p.parseMinute(parts[0])
	if err != nil {
		return nil, err
	}

	hour, err := p.parseHour(parts[1])
	if err != nil {
		return nil, err
	}

	dayOfMonth, err := p.parseDayOfMonth(parts[2])
	if err != nil {
		return nil, err
	}

	month, err := p.parseMonth(parts[3])
	if err != nil {
		return nil, err
	}

	dayOfWeek, err := p.parseDayOfWeek(parts[4])
	if err != nil {
		return nil, err
	}

	return &Cron{
		Minute:     minute,
		Hour:       hour,
		DayOfMonth: dayOfMonth,
		Month:      month,
		DayOfWeek:  dayOfWeek,
		Command:    strings.Join(parts[5:len(parts)], " "),
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

	var out []int
	outMap := make(map[int]struct{})
	parts := strings.Split(expression, ",")

	for _, part := range parts {
		rangeParts := strings.Split(part, "-")

		if len(rangeParts) == 2 {
			start, _ := strconv.Atoi(rangeParts[0])
			end, _ := strconv.Atoi(rangeParts[1])

			for _, i := range numRange(start, end) {
				outMap[i]= struct{}{}
			}
		} else {
			i, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid value %s, expected number", part)
			}

			if i < min || i > max {
				return nil, fmt.Errorf("unexpected value %d, expected value between %d and %d", i, min, max)
			}

			outMap[i]= struct{}{}
		}
	}

	for key := range outMap {
		out = append(out, key)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i] < out[j]
	})

	return out, nil
}