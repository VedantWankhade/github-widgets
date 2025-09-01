package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vedantwankhade/github-widgets/internal/config"
)

const (
	GHAPIBaseUrl                      = "https://api.github.com"
	RepositoriesEndpointTemplate      = "/users/%s/repos"
	WeeklyCommitCountEndpoingTemplate = "/repos/%s/stats/participation"
)

var apiClient = &http.Client{
	Timeout: 2 * time.Second,
}

var app = config.GetApp()

func GetWeeklyCommitCount(repo string, token string) []int {

	type resType struct {
		Owner []int `json:"owner"`
	}

	url := fmt.Sprintf(GHAPIBaseUrl+WeeklyCommitCountEndpoingTemplate, repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		app.Error(err.Error())
		return nil
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	res, err := apiClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		app.Error(fmt.Sprintf("commit retrivel status: %d", res.StatusCode))
		return nil
	}
	defer res.Body.Close()
	var weeklyCommits resType

	if err = json.NewDecoder(res.Body).Decode(&weeklyCommits); err != nil {
		app.Error(err.Error())
		return nil
	}

	return weeklyCommits.Owner
}

func GetGHUserRepos[T any](user, token string, target *[]T) {
	perPage := 100
	page := 0

	for {
		done, err := getGHUserReposPaginate[T](user, perPage, page+1, token, target)
		if err != nil {
			app.Error(err.Error())
			break
		}
		if done {
			app.Info("repos collected")
			break
		}
		page++
	}
}

func getGHUserReposPaginate[T any](user string, perPage, pageNo int, token string, target *[]T) (bool, error) {
	queryParams := fmt.Sprintf("?sort=pushed&direction=desc&per_page=%d&page=%d", perPage, pageNo)
	url := fmt.Sprintf(GHAPIBaseUrl+RepositoriesEndpointTemplate+queryParams, user)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	res, err := apiClient.Do(req)
	if err != nil {
		return false, err
	}
	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("errror: request status: %d", res.StatusCode)
	}

	defer res.Body.Close()

	var auxTarget []T
	if err = json.NewDecoder(res.Body).Decode(&auxTarget); err != nil {
		return false, err
	}

	*target = append(*target, auxTarget...)
	if len(auxTarget) < perPage {
		return true, nil
	}

	return false, nil
}
