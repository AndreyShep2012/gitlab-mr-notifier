package interfaces

import "gitlab-mr-notifier/models"

type Notifier interface {
	Notify(creds interface{}, mrs models.MergeRequests) error
}
