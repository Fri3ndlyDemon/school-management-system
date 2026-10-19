package service

import (
	"database/sql"
	"fmt"
)

type ProxyEngine struct {
	db *sql.DB
}

func NewProxyEngine(db *sql.DB) *ProxyEngine {
	return &ProxyEngine{db: db}
}

// AssignProxiesForDay runs daily to re-route classes if primary teachers are absent
func (p *ProxyEngine) AssignProxiesForDay(date string) error {
	// 1. Find all classes scheduled for this date/day of the week
	// 2. Check if the primary assigned teacher has an approved leave entry
	// 3. If primary teacher is absent, fetch the course's proxy pool (array of teacher IDs)
	// 4. Find the first available teacher from the pool who has no conflicting classes today
	// 5. Save the proxy assignment to the database session table

	fmt.Println("Running automated proxy allocation for date:", date)
	return nil
}