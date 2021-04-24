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

	minute, err := p.parseExpression(parts[0], 0, 59)
	if err != nil {
		return nil, err
	}

	hour, err := p.parseExpression(parts[1], 0, 23)
	if err != nil {
		return nil, err
	}

	dayOfMonth, err := p.parseExpression(parts[2], 1, 31)
	if err != nil {
		return nil, err
	}

	month, err := p.parseExpression(parts[3], 1, 12)
	if err != nil {
		return nil, err
	}

	dayOfWeek, err := p.parseExpression(parts[4], 0, 6)
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

func (p *CronParser) parseExpression(expression string, min int, max int) ([]int, error) {
	if expression == "*" {
		return numRange(min, max), nil
	}

	var out []int
	outMap := make(map[int]struct{})
	parts := strings.Split(expression, ",")

	for _, part := range parts {
		if rangeParts := strings.Split(part, "-"); len(rangeParts) == 2 {
			err := p.getRange(rangeParts[0], rangeParts[1], min, max, outMap)
			if err != nil {
				return nil, err
			}
		} else if stepParts := strings.Split(part, "/"); len(stepParts) == 2 {
			err := p.getSteps(stepParts[0], stepParts[1], min, max, outMap)
			if err != nil {
				return nil, err
			}
		} else {
			i, err := p.parseIntInRange(part, min, max)
			if err != nil {
				return nil, err
			}

			outMap[i] = struct{}{}
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

func (p *CronParser) getRange(startString string, endString string, min int, max int, outMap map[int]struct{}) error {
	start, err := p.parseIntInRange(startString, min, max)
	if err != nil {
		return err
	}
	end, err := p.parseIntInRange(endString, min, max)
	if err != nil {
		return err
	}

	if start <= end {
		for _, i := range numRange(start, end) {
			outMap[i] = struct{}{}
		}
	} else {
		for _, i := range numRange(min, end) {
			outMap[i] = struct{}{}
		}
		for _, i := range numRange(start, max) {
			outMap[i] = struct{}{}
		}
	}

	return nil
}

func (p *CronParser) getSteps(startString string, stepString string, min int, max int, outMap map[int]struct{}) error {
	i, err := p.parseIntInRange(startString, min, max)
	if err != nil {
		return err
	}
	step, err := p.parseInt(stepString)
	if err != nil {
		return err
	}

	for ; i <= max; i += step {
		outMap[i] = struct{}{}
	}

	return nil
}

func (p *CronParser) parseIntInRange(number string, min, max int) (int, error) {
	i, err := p.parseInt(number)
	if err != nil {
		return 0, err
	}

	if i < min || i > max {
		return 0, fmt.Errorf("unexpected value %d, expected value between %d and %d", i, min, max)
	}

	return i, nil
}

func (p *CronParser) parseInt(number string) (int, error) {
	i, err := strconv.Atoi(number)
	if err != nil {
		return 0, fmt.Errorf("invalid value %s, expected number", number)
	}

	return i, nil
}
