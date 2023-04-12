package interfaces

import (
	"github.com/xanzy/go-gitlab"
)

type GitlabApi interface {
	GetMRList(groupid int) ([]*gitlab.MergeRequest, error)
	GetMRDiscussions(projectID, mergeRequestID, page, perPage int) ([]*gitlab.Discussion, error)
}
