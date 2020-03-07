package domain

import "time"

type Program struct {
	ID       string
	Name     string
	HostName string
	Duration time.Duration
}
