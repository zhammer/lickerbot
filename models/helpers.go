package models

import "github.com/lib/pq"

// UniqueViolationErr returns true if an error is a postgres
// unique violation error.
func UniqueViolationErr(err error) bool {
	if err == nil {
		return false
	}
	pgErr, ok := err.(*pq.Error)
	if !ok {
		return false
	}

	return pgErr.Code.Name() == "unique_violation"
}
