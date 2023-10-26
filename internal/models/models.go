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
	ChangesCount        string
	Branch              string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	PipelineInfo        PipelineInfo
}

type MergeRequests []MergeRequest

type PipelineInfo struct {
	IsFailed bool
	URL      string
}
