package main

import (
	"time"
)

type HistoryRecord struct {
	Query     string `json:"query"`
	Timestamp string `json:"timestamp"`
}

func NewHistory() []HistoryRecord {
	return make([]HistoryRecord, 0)
}

func NewHistoryRecord(query string) HistoryRecord {
	return HistoryRecord{
		Query:     query,
		Timestamp: time.Now().String(),
	}
}
