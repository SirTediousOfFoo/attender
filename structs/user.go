// Package structs provides structures for configuration
package structs
import "time"

// User holds the data for our current user
type User struct {
	ID            int
	Name          string
	Surname       string
	Username      string
	Email         string
	Authenticated bool
	Stats         Stats
	Admin         bool
}

// Stats holds attendance data for our current user, probably will be nil for most cases except for the stats page
type Stats struct {
	AttendedTotal   int
	AttendedMonthly int
	AttendedYearly  int
}

// Member holds the data for a member of the group
type Member struct {
	ID       int
	Name     string
	Surname  string
	Dates		[]time.Time
}
