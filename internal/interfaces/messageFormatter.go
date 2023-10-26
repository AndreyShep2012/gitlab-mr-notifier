package interfaces

import "gitlab-mr-notifier/internal/models"

type MessageFormatter interface {
	GetBody(models.MergeRequest) string
	GetIntroText(mrs int) string
	GetPipelineFailedIntroText(mrs int) string
	GetPipelineFailedGetBody([]models.MergeRequest) string
}
