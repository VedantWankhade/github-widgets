package widgets

import (
	"bytes"
	"fmt"
	"io"
	"net/url"

	"github.com/vedantwankhade/github-widgets/internal/config"
	"github.com/vedantwankhade/github-widgets/utils"
)

type GHUserRepo struct {
	RepoName string `json:"full_name"`
}

func CommitGraph(params url.Values) (io.ReadSeeker, error) {
	app := config.GetApp()

	ghToken := app.GetGHToken()
	if ghToken == "" {
		app.Warn("no 'GH_AUTH_TOKEN' set")
	}

	username := params.Get("user")
	if username == "" {
		app.Error("parameter 'user' not found in url")
		return nil, fmt.Errorf("please provide parameter 'user' for widget 'commitgraph'")
	}

	var repos []GHUserRepo

	utils.GetGHUserRepos[GHUserRepo]("octocat", app.GetGHToken(), &repos)

	for i, r := range repos {
		fmt.Printf("%d. %s\n", i, r.RepoName)
	}

	res := bytes.NewReader([]byte(`

	<svg width="200" height="250" version="1.1" xmlns="http://www.w3.org/2000/svg">
		<circle cx="100" cy="100" r="50"/>
		</svg>

		`))
	return res, nil
}
