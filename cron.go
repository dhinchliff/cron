package main

type Cron struct {
	// Todo: replace slices with bit fields
	Minute     []int
	Hour       []int
	DayOfMonth []int
	Month      []int
	DayOfWeek  []int
	Command    string
}
