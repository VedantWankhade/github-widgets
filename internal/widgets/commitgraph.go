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
	weeklyCommits := comms[52-20 : 52]
	app.Info(fmt.Sprintf("weekly commits retrieved for %s", username))

	height := 200
	dx := 30

	var ticks bytes.Buffer
	for r := range 8 {
		ticks.WriteString(fmt.Sprintf(`
  <line class="grid" x1="20" y1="%d" x2="1000" y2="%d"/>
  <text class="label" x="20" y="%d" text-anchor="end">%d</text>
			`, height-r*20, height-r*20, height-r*20, 5*r))
	}

	var bars bytes.Buffer
	for w, c := range weeklyCommits {
		bars.WriteString(fmt.Sprintf(`
  <rect class="bar" x="%d" y="%d" width="20" height="%d" rx="4" />
			`, 20+dx*w, 200-c*4, c*4))
	}

	res := bytes.NewReader([]byte(fmt.Sprintf(`


	<svg width="1000" height="400" viewBox="0 0 500 320" xmlns="http://www.w3.org/2000/svg" role="img" aria-labelledby="title desc">
  <title id="title">SVG Bar Chart with Axis Titles</title>
  <desc id="desc">A simple bar chart showing five bars with labeled X and Y axes.</desc>

  <!-- Styles -->
  <style>
    .axis { stroke: #222; stroke-width: 2; }
    .grid { stroke: #bbb; stroke-dasharray: 3 3; stroke-width: 1; }
    .tick { stroke: #222; stroke-width: 1; }
    .label { font: 12px system-ui, sans-serif; fill: #222; }
    .title { font: 14px system-ui, sans-serif; font-weight: 600; fill: #111; }
    .bar { fill: #4E79A7; }
  </style>

  <!-- Axes -->
  <line class="axis" x1="20" y1="0"  x2="20"  y2="200"/>   <!-- Y axis -->
  <line class="axis" x1="20" y1="200" x2="220" y2="200"/>  <!-- X axis -->

  <!-- Y ticks & grid (0..30 step 10) -->
	%s  
  
  <!-- Bars (example data: [12, 25, 9, 18, 30]) -->
		%s
  <!-- X labels -->
  <!-- <text class="label" x="100" y="278" text-anchor="middle">A</text>
  <text class="label" x="180" y="278" text-anchor="middle">B</text>
  <text class="label" x="260" y="278" text-anchor="middle">C</text>
  <text class="label" x="340" y="278" text-anchor="middle">D</text>
  <text class="label" x="420" y="278" text-anchor="middle">E</text>
 -->
  <!-- Axis titles -->
  <!-- <text class="title" x="260" y="310" text-anchor="middle">Categories</text>
  <text class="title" transform="translate(15 140) rotate(-90)" text-anchor="middle">Value</text> -->
</svg>
	

		`, ticks.String(), bars.String())))
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
