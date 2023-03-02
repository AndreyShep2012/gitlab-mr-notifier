package interfaces

import "gitlab-mr-notifier/models"

type GitApi interface {
	GetMRList(token string, groupid int) (models.MergeRequests, error)
}
