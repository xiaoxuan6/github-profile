package github

import (
	"context"
	"fmt"
	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v66/github"
	"log"
	"os"
	"strings"
)

var Client *github.Client

func Init() {
	rateLimiter, err := github_ratelimit.NewRateLimitWaiterClient(nil)
	if err != nil {
		panic(err)
	}

	Client = github.NewClient(rateLimiter).WithAuthToken(os.Getenv("GITHUB_TOKEN"))
}

// FetchRepository repository => https://github.com/xiaoxuan6/github-profile
func FetchRepository(repository string) (*github.Repository, error) {
	sl := strings.Split(repository, "/")
	repos, _, err := Client.Repositories.Get(context.Background(), sl[3], sl[4])
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func FetchAllRepository(username string) []*github.Repository {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	count := 0
	var allRepository []*github.Repository
	for {
		allRepos, response, err := Client.Repositories.List(context.Background(), username, opt)
		if err != nil {
			if count == 3 {
				break
			}
			count += 1
			log.Println(fmt.Sprintf("Repository list num: %d err：%s", count, err.Error()))
			continue
		}

		allRepository = append(allRepository, allRepos...)

		if response.NextPage == 0 {
			break
		}
		opt.Page = response.NextPage
	}

	return allRepository
}

func FetchAllPrs(username string) []*github.Issue {
	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	count := 0
	var allIssue []*github.Issue
	for {
		allIssues, response, err := Client.Search.Issues(context.Background(), fmt.Sprintf("is:pr author:%s", username), opt)
		if err != nil {
			if count == 3 {
				break
			}
			count += 1
			log.Println(fmt.Sprintf("Pr repository num: %d err：%s", count, err.Error()))
			continue
		}

		allIssue = append(allIssue, allIssues.Issues...)

		if len(allIssue) >= 1000 {
			break
		}

		if response.NextPage == 0 {
			break
		}
		opt.Page = response.NextPage
	}

	return allIssue
}
