package models

import "time"

type MergeRequest struct {
	Title               string
	Description         string
	Author              string
	URL                 string
	HasConflicts        bool
	UnresolvedThreads   int
	DetailedMergeStatus string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type MergeRequests []MergeRequest
