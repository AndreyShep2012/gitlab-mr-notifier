package models

type MergeRequest struct {
	Title                         string
	Description                   string
	Author                        string
	UserNotesCount                int
	URL                           string
	HasConflicts                  bool
	IsBlockingDiscussionsResolved bool
}

type MergeRequests []MergeRequest
