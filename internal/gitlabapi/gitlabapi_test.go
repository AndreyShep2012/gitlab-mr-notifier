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

func TestGitlabapi(t *testing.T) {
	require.NoError(t, godotenv.Load("../../.env"))

	token := os.Getenv("GITLAB_TOKEN")
	require.NotEmpty(t, token)

	id := os.Getenv("GITLAB_GROUP_ID")
	require.NotEmpty(t, token)

	groupid, err := strconv.Atoi(id)
	require.NoError(t, err)

	api, err := gitlabapi.New(token)
	require.NoError(t, err)

	res, err := api.GetMRList(groupid)
	require.NoError(t, err)
	fmt.Println("GetMRList result: ", res)

	for _, mr := range res {
		discussions, err := api.GetMRDiscussions(mr.ProjectID, mr.IID, 0, 10)
		require.NoError(t, err)
		fmt.Println("MR discussions: ", discussions)
	}
}

func TestGetMRListEmptyCreds(t *testing.T) {
	require.NoError(t, godotenv.Load("../../.env"))

	token := os.Getenv("GITLAB_TOKEN")
	require.NotEmpty(t, token)

	api, err := gitlabapi.New(token)
	require.NoError(t, err)

	res, err := api.GetMRList(0)
	require.Error(t, err)
	require.Empty(t, res)

	id := os.Getenv("GITLAB_GROUP_ID")
	require.NotEmpty(t, token)

	groupid, err := strconv.Atoi(id)
	require.NoError(t, err)

	res, err = api.GetMRList(groupid)
	require.Error(t, err)
	require.Empty(t, res)
}
