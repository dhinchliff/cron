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
	fields := strings.Split(expression, " ")

	if len(fields) < 5 {
		return nil, fmt.Errorf("invalid expression [%s]", expression)
	}

	minute, err := parseField(fields[0], 0, 59)
	if err != nil {
		return nil, err
	}

	hour, err := parseField(fields[1], 0, 23)
	if err != nil {
		return nil, err
	}

	dayOfMonth, err := parseField(fields[2], 1, 31)
	if err != nil {
		return nil, err
	}

	month, err := parseField(fields[3], 1, 12)
	if err != nil {
		return nil, err
	}

	dayOfWeek, err := parseField(fields[4], 0, 6)
	if err != nil {
		return nil, err
	}

	return &Cron{
		Minute:     minute,
		Hour:       hour,
		DayOfMonth: dayOfMonth,
		Month:      month,
		DayOfWeek:  dayOfWeek,
		Command:    strings.Join(fields[5:], " "),
	}, nil
}

func parseField(field string, min int, max int) ([]int, error) {
	var out []int
	outMap := make(map[int]struct{})
	expressions := strings.Split(field, ",")

	for _, expression := range expressions {
		if expression == "" {
			return nil, fmt.Errorf("empty field, check for double spaces")
		}

		err := parseExpression(expression, min, max, outMap)
		if err != nil {
			return nil, fmt.Errorf("invalid field %s: %s", field, err)
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

func parseExpression(expression string, min int, max int, outMap map[int]struct{}) (err error) {
	rangeAndStep := strings.Split(expression, "/")
	startAndEnd := strings.Split(rangeAndStep[0], "-")
	start := min
	end := max
	step := 1

	if len(rangeAndStep) > 2 {
		return fmt.Errorf("too many slashes")
	}

	if len(startAndEnd) > 2 {
		return fmt.Errorf("too many hyphens")
	}

	if len(rangeAndStep) == 2 {
		step, err = parseInt(rangeAndStep[1])
		if err != nil {
			return err
		}
	}

	if rangeAndStep[0] != "*" && rangeAndStep[0] != "" {
		start, err = parseIntInRange(startAndEnd[0], min, max)
		if err != nil {
			return err
		}

		if len(startAndEnd) == 2 {
			end, err = parseIntInRange(startAndEnd[1], min, max)
			if err != nil {
				return err
			}
		} else if len(rangeAndStep) == 1 {
			end = start
		}
	}

	err = getRangeStep(start, end, step, min, max, outMap)
	if err != nil {
		return err
	}

	return nil
}

func getRangeStep(start int, end int, step int, min int, max int, outMap map[int]struct{}) error {
	if step < 1 {
		return fmt.Errorf("unexpected step %d, expected value greater than zero", step)
	}

	if start <= end {
		for i := start; i <= end; i += step {
			outMap[i] = struct{}{}
		}
	} else {
		var i int
		for i = start; i <= max; i += step {
			outMap[i] = struct{}{}
		}
		for i = i - max - 1; i <= end; i += step {
			outMap[i] = struct{}{}
		}
	}

	return nil
}

func parseIntInRange(number string, min, max int) (int, error) {
	i, err := parseInt(number)
	if err != nil {
		return 0, err
	}

	if i < min || i > max {
		return 0, fmt.Errorf("unexpected value %d, expected value between %d and %d", i, min, max)
	}

	return i, nil
}

func parseInt(number string) (int, error) {
	if strings.TrimSpace(number) == "" {
		return 0, fmt.Errorf("empty value, expected number")
	}

	i, err := strconv.Atoi(number)
	if err != nil {
		return 0, fmt.Errorf("invalid value %s, expected number", number)
	}

	return i, nil
}
