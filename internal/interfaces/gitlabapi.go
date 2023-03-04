package interfaces

import "gitlab-mr-notifier/internal/models"

type GitlabApi interface {
	GetMRList(token string, groupid int) (models.MergeRequests, error)
}
