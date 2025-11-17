package slack

import (
	"errors"

	"gitlab-mr-notifier/internal/interfaces"
	"gitlab-mr-notifier/internal/models"

	slack_api "github.com/easonlin404/go-slack"
)

type slack struct {
	formatter interfaces.MessageFormatter
}

func New(formatter interfaces.MessageFormatter) interfaces.Notifier {
	return &slack{formatter: formatter}
}

func (s slack) Notify(creds any, mrs models.MergeRequests) error {
	webhookURL, ok := creds.(string)
	if !ok {
		return errors.New("wrong credentials format, please use string for webhook url")
	}

	api := slack_api.New().WebhookURL(webhookURL)
	successPipeline, failedPipeline := sortMRs(mrs)
	if err := s.sendSuccessPipelineMRs(api, successPipeline); err != nil {
		return err
	}

	if err := s.sendFailedPipelineMRs(api, failedPipeline); err != nil {
		return err
	}

	return nil
}

func (s slack) sendSuccessPipelineMRs(api *slack_api.Slack, mrs models.MergeRequests) error {
	_, err := api.SendMessage(slack_api.Message{Text: s.formatter.GetIntroText(len(mrs))})
	if err != nil {
		return err
	}

	for _, m := range mrs {
		_, err := api.SendMessage(slack_api.Message{Text: s.formatter.GetBody(m)})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s slack) sendFailedPipelineMRs(api *slack_api.Slack, mrs models.MergeRequests) error {
	_, err := api.SendMessage(slack_api.Message{Text: s.formatter.GetPipelineFailedIntroText(len(mrs))})
	if err != nil {
		return err
	}

	_, err = api.SendMessage(slack_api.Message{Text: s.formatter.GetPipelineFailedGetBody(mrs)})
	return err
}

func sortMRs(mrs models.MergeRequests) (successPipeline, failedPipeline models.MergeRequests) {
	for _, mr := range mrs {
		arr := &successPipeline
		if mr.PipelineInfo.IsFailed {
			arr = &failedPipeline
		}

		*arr = append(*arr, mr)
	}
	return
}
