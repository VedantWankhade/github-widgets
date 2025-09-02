package widgets

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
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

	start := time.Now()

	comms := getWeeklyCommits(username, ghToken)
	// weeklyCommits := [52]int{5, 10, 15, 30, 20}
	numWeeks := 20
	weeklyCommits := comms[52-numWeeks : 52]
	app.Info(fmt.Sprintf("weekly commits retrieved for %s", username))

	height := 400
	maxBarHeight := height / 2
	width := 860
	graphStartX := 40
	dx := 40

	var ticks bytes.Buffer
	for r := range 10 {
		ticks.WriteString(fmt.Sprintf(`
  <line class="grid" x1="%d" y1="%d" x2="%d" y2="%d"/>
  <text class="label" x="%d" y="%d" text-anchor="end">%d</text>
			`, graphStartX+20, maxBarHeight-r*20, width, maxBarHeight-r*20, graphStartX+10, maxBarHeight-r*20, 5*r))
	}

	var bars bytes.Buffer
	for w, c := range weeklyCommits {
		bars.WriteString(fmt.Sprintf(`
  <rect class="bar" x="%d" y="%d" width="20" height="%d" rx="4" />
			`, graphStartX+20+dx*w, 200-c*4, c*4))
	}

	var axes bytes.Buffer
	axes.WriteString(fmt.Sprintf(`
  <line class="axis" x1="%d" y1="0"  x2="%d"  y2="%d"/>   <!-- Y axis -->
  <line class="axis" x1="%d" y1="200" x2="%d" y2="%d"/>  <!-- X axis -->
		`, graphStartX+20, graphStartX+20, maxBarHeight, graphStartX+20, width, maxBarHeight))

	endDate := time.Now()
	dateRanges := make([]struct {
		startDate time.Time
		endDate   time.Time
	}, numWeeks)

	var xLabels bytes.Buffer
	for w := numWeeks - 1; w >= 0; w-- {
		weekEnd := endDate.AddDate(0, 0, -7*(numWeeks-1-w))
		weekStart := weekEnd.AddDate(0, 0, -6)

		dateRanges[w] = struct {
			startDate time.Time
			endDate   time.Time
		}{
			weekStart,
			weekEnd,
		}
	}

	for w, r := range dateRanges {
		xLabels.WriteString(fmt.Sprintf(`
		  <text class="label" transform="translate(%d 260) rotate(-90)" text-anchor="middle">%s</text>
			`, graphStartX+40+dx*w, r.startDate.Format("02 Jan")+" - "+r.endDate.Format("02 Jan")))
	}

	res := bytes.NewReader([]byte(fmt.Sprintf(`

<svg fill="currentColor" width="%d" height="%d" xmlns="http://www.w3.org/2000/svg" role="img" aria-labelledby="title desc">
  <title id="title">Commitgraph</title>

  <!-- Styles -->
  <style>
    .axis { stroke: #2ea44f; stroke-width: 2; }
.grid { stroke: #2ea44f; stroke-opacity: 0.3; stroke-dasharray: 3 3; stroke-width: 1; }
.tick { stroke: #2ea44f; stroke-width: 1; }
.label { font: 12px system-ui, sans-serif; fill: #2ea44f; }
.title { font: 14px system-ui, sans-serif; font-weight: 600; fill: #2ea44f; }
    .bar { fill: #a6e9b6; stroke: #2ea44f; stroke-width: 2; }
  </style>

  <!-- Axes -->
%s
  <!-- Y ticks & grid (0..30 step 10) -->
	%s  
  
  <!-- Bars (example data: [12, 25, 9, 18, 30]) -->
		%s
  <!-- X labels -->
%s
  <!-- Axis titles -->
  <text class="title" x="420" y="320" text-anchor="middle">Weeks</text>
  <text class="title" transform="translate(20 140) rotate(-90)" text-anchor="middle">Commits</text>
</svg>
		`, width, height, axes.String(), ticks.String(), bars.String(), xLabels.String())))
	app.Info(fmt.Sprintf("[Generating svg took %dms]", time.Since(start).Milliseconds()))
	return res, nil
}

func perRepoCommitsWorker(user, token string, jobId int, repoName <-chan string, res chan<- []int) {
	app := config.GetApp()

	for j := range repoName {
		app.Info(fmt.Sprintf("worker %d started for repo %s", jobId, j))
		weeklyCommits := utils.GetWeeklyCommitCount(j, token)
		app.Info(fmt.Sprintf("worker %d ended for repo %s", jobId, j))
		res <- weeklyCommits
	}
}

func getWeeklyCommits(user string, token string) [52]int {
	app := config.GetApp()
	var repos []GHUserRepo
	var totalWeeklyCommits [52]int

	utils.GetGHUserRepos[GHUserRepo](user, token, &repos)

	// fire workers
	numWorkers := len(repos)
	app.Info(fmt.Sprintf("Firing off %d workers for %d repos", numWorkers, len(repos)))

	repoName := make(chan string, len(repos))
	res := make(chan []int, len(repos))

	for w := range numWorkers {
		go perRepoCommitsWorker(user, token, w, repoName, res)
	}

	// send to workers
	for _, r := range repos {
		repoName <- r.RepoName
	}

	close(repoName)

	// collect
	for range repos {
		weeklyCommits := <-res
		for i := 0; i < len(weeklyCommits) && i < 52; i++ {
			totalWeeklyCommits[i] += weeklyCommits[i]
		}
	}

	return totalWeeklyCommits
}
