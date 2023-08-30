package runner

import (
	"log"

	"gitlab-mr-notifier/internal/config"
	"gitlab-mr-notifier/internal/gitlabapi"
	lognotifier "gitlab-mr-notifier/internal/logNotifier"
	"gitlab-mr-notifier/internal/slack"
)

func check(config config.Config) {
	f := slack.NewSimpleMessageFormatter(config.MessageDescriptionLimit)
	notif := slack.New(f)
	if config.Notifier == "log" {
		notif = lognotifier.New(f)
	}
	gitapi := gitlabapi.New()

	log.Println("start checking")
	mrs, err := gitapi.GetMRList(config.GitlabToken, config.GitlabGroupID)
	if err != nil {
		log.Println("getting mr list error:", err)
	}
	log.Println("found ", len(mrs), " MRs")
	err = notif.Notify(config.SlackWebhookURL, mrs)
	if err != nil {
		log.Println("slack notification error:", err)
	}
}
