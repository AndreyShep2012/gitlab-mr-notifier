package config_test

import (
	"gitlab-mr-notifier/internal/config"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadFromFile(t *testing.T) {
	f, err := os.CreateTemp("", "tmpfile-")
	require.NoError(t, err)

	defer f.Close()
	defer os.Remove(f.Name())

	params := []string{
		"GITLAB_TOKEN=token",
		"GITLAB_GROUP_ID=64504965",
		"SLACK_WEBHOOK_URL=https://hooks.slack.com/services/T04RWJXV6KC",
		"CRON_PERIOD=10s",
	}

	for _, s := range params {
		_, err := f.WriteString(s + "\n")
		require.NoError(t, err)
	}

	require.NoError(t, os.Setenv("CONFIG_PATH", f.Name()))

	c, err := config.Load()
	require.NoError(t, err)
	require.Equal(t, "token", c.GitlabToken)
	require.Equal(t, 64504965, c.GitlabGroupID)
	require.Equal(t, "https://hooks.slack.com/services/T04RWJXV6KC", c.SlackWebhookURL)
	require.Equal(t, "10s", c.CronPeriod)
	require.Empty(t, c.CronTime)
	require.Equal(t, 500, c.MessageDescriptionLimit)
	require.Equal(t, "slack", c.Notifier)

	_, err = f.WriteString("CRON_TIME=10:30\n")
	require.NoError(t, err)

	_, err = f.WriteString("MESSAGE_DESCRIPTION_LIMIT=10\n")
	require.NoError(t, err)

	_, err = f.WriteString("NOTIFIER=log\n")
	require.NoError(t, err)

	c, err = config.Load()
	require.NoError(t, err)

	require.Equal(t, "10:30", c.CronTime)
	require.Equal(t, 10, c.MessageDescriptionLimit)
	require.Equal(t, "log", c.Notifier)

	c, err = config.Load()
	require.NoError(t, err)

	_, err = f.Seek(0, 0)
	require.NoError(t, err)

	_, err = f.WriteString("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	require.NoError(t, err)

	c, err = config.Load()
	require.Empty(t, c)
	require.ErrorContains(t, err, "load config file")

	require.NoError(t, os.Setenv("CONFIG_PATH", "unknown"))
	c, err = config.Load()
	require.Empty(t, c)
	require.ErrorContains(t, err, "load config file")
}

func TestFromEnv(t *testing.T) {
	require.NoError(t, os.Setenv("CONFIG_PATH", ""))
	require.NoError(t, os.Setenv("GITLAB_TOKEN", "token"))
	require.NoError(t, os.Setenv("GITLAB_GROUP_ID", "64504965"))
	require.NoError(t, os.Setenv("SLACK_WEBHOOK_URL", "https://hooks.slack.com/services/T04RWJXV6KC"))
	require.NoError(t, os.Setenv("CRON_PERIOD", "10s"))
	require.NoError(t, os.Setenv("CRON_TIME", "10:30"))
	require.NoError(t, os.Setenv("MESSAGE_DESCRIPTION_LIMIT", "10"))
	require.NoError(t, os.Setenv("NOTIFIER", "log"))

	c, err := config.Load()
	require.NoError(t, err)
	require.Equal(t, "token", c.GitlabToken)
	require.Equal(t, 64504965, c.GitlabGroupID)
	require.Equal(t, "https://hooks.slack.com/services/T04RWJXV6KC", c.SlackWebhookURL)
	require.Equal(t, "10s", c.CronPeriod)
	require.Equal(t, "10:30", c.CronTime)
	require.Equal(t, 10, c.MessageDescriptionLimit)
	require.Equal(t, "log", c.Notifier)

	require.NoError(t, os.Unsetenv("MESSAGE_DESCRIPTION_LIMIT"))
	require.NoError(t, os.Unsetenv("NOTIFIER"))
	c, err = config.Load()
	require.NoError(t, err)
	require.Equal(t, 500, c.MessageDescriptionLimit)
	require.Equal(t, "slack", c.Notifier)

	require.NoError(t, os.Unsetenv("GITLAB_TOKEN"))

	c, err = config.Load()
	require.ErrorContains(t, err, "required")
	require.Empty(t, c)
}

func TestRequired(t *testing.T) {
	require.NoError(t, os.Setenv("CONFIG_PATH", ""))
	require.NoError(t, os.Setenv("GITLAB_TOKEN", ""))
	require.NoError(t, os.Setenv("GITLAB_GROUP_ID", "0"))
	require.NoError(t, os.Setenv("SLACK_WEBHOOK_URL", ""))

	c, err := config.Load()
	require.ErrorContains(t, err, "GITLAB_TOKEN")
	require.ErrorContains(t, err, "GITLAB_GROUP_ID")
	require.ErrorContains(t, err, "SLACK_WEBHOOK_URL")
	require.Empty(t, c)
}
