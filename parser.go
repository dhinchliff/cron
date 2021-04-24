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
		rangeParts := strings.Split(part, "-")

		if len(rangeParts) == 2 {
			start, _ := strconv.Atoi(rangeParts[0])
			end, _ := strconv.Atoi(rangeParts[1])

			if start < min || start > max {
				return nil, fmt.Errorf("unexpected value %d, expected value between %d and %d", start, min, max)
			}

			if end < min || end > max {
				return nil, fmt.Errorf("unexpected value %d, expected value between %d and %d", end, min, max)
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
		} else {
			stepParts := strings.Split(part, "/")
			if len(stepParts) == 2 {
				i, _ := strconv.Atoi(stepParts[0])
				step, _ := strconv.Atoi(stepParts[1])

				if i < min || i > max {
					return nil, fmt.Errorf("unexpected value %d, expected value between %d and %d", i, min, max)
				}

				for ; i <= max; i += step {
					outMap[i] = struct{}{}
				}
			} else {
				i, err := strconv.Atoi(part)
				if err != nil {
					return nil, fmt.Errorf("invalid value %s, expected number", part)
				}

				if i < min || i > max {
					return nil, fmt.Errorf("unexpected value %d, expected value between %d and %d", i, min, max)
				}

				outMap[i] = struct{}{}
			}
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
