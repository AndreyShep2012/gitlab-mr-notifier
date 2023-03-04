package gitlabapi_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"gitlab-mr-notifier/internal/gitlabapi"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestGetMRList(t *testing.T) {
	require.NoError(t, godotenv.Load("../../.env"))

	token := os.Getenv("GITLAB_TOKEN")
	require.NotEmpty(t, token)

	id := os.Getenv("GITLAB_GROUP_ID")
	require.NotEmpty(t, token)

	groupid, err := strconv.Atoi(id)
	require.NoError(t, err)

	api := gitlabapi.New()
	res, err := api.GetMRList(token, groupid)
	fmt.Println("GetMRList result: ", res)
	require.NoError(t, err)
}

func TestGetMRListEmptyCreds(t *testing.T) {
	require.NoError(t, godotenv.Load("../../.env"))

	token := os.Getenv("GITLAB_TOKEN")
	require.NotEmpty(t, token)

	api := gitlabapi.New()
	res, err := api.GetMRList(token, 0)
	require.Error(t, err)
	require.Empty(t, res)

	id := os.Getenv("GITLAB_GROUP_ID")
	require.NotEmpty(t, token)

	groupid, err := strconv.Atoi(id)
	require.NoError(t, err)

	res, err = api.GetMRList("", groupid)
	require.Error(t, err)
	require.Empty(t, res)
}
