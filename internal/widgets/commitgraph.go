package widgets

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

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

	weeklyCommits := getWeeklyCommits(username, ghToken)

	app.Info(fmt.Sprintf("weekly commits retrieved for %s", username))

	var points [52]string
	dx := 16
	height := 200
	x := 0
	start := time.Now()
	for i := range points {
		points[i] = fmt.Sprintf("%d,%d", x+dx*i, height-weeklyCommits[i])
	}

	res := bytes.NewReader([]byte(fmt.Sprintf(`
	<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" role="img" aria-label="commitgraph">
	  <title>commitgraph</title>

			<polyline
				 fill="none"
				 stroke="#0074d9"
				 stroke-width="3"
				 points="%s"
			/>
		</svg>
		`, dx*52, height, strings.Join(points[:], "\n"))))
	app.Info(fmt.Sprintf("[Generating svg took %dms]", time.Since(start).Milliseconds()))
	return res, nil
}

func getWeeklyCommits(user string, token string) [52]int {
	app := config.GetApp()
	var repos []GHUserRepo
	var totalWeeklyCommits [52]int
	start := time.Now()
	utils.GetGHUserRepos[GHUserRepo](user, token, &repos)
	app.Info(fmt.Sprintf("[Collecting repos took %dms]", time.Since(start).Milliseconds()))
	for _, r := range repos {
		start := time.Now()
		weeklyCommits := utils.GetWeeklyCommitCount(r.RepoName, token)
		app.Info(fmt.Sprintf("[Getting commits for %s took %dms]", r.RepoName, time.Since(start).Milliseconds()))
		start = time.Now()
		for i := 0; i < len(weeklyCommits) && i < 52; i++ {
			totalWeeklyCommits[i] += weeklyCommits[i]
		}
		app.Info(fmt.Sprintf("[Adding weekly data took %dms]", time.Since(start).Milliseconds()))
	}

	return totalWeeklyCommits
}
