package services

import (
	"fmt"
	github2 "github.com/google/go-github/v66/github"
	"github.com/xiaoxuan6/github-profile/github"
	"sort"
	"strings"
	"sync"
)

type Profile struct {
	Name        string
	HTMLUrl     string
	Description string
	Stars       int
	Forks       int
	Create      string
	Update      string
	Language    string
}

func (p Profile) FullName() string {
	return fmt.Sprintf("[%s](%s)", p.Name, p.HTMLUrl)
}

type PrRepository struct {
	Name        string
	Url         string
	Description string
	Language    string
	State       string
	Created     string
	Count       int
}

func (p PrRepository) FullName() string {
	reposName := strings.Split(p.Name, "/")[5]
	return fmt.Sprintf("[%s](%s)", reposName, p.Url)
}

func (p PrRepository) CountUrl(username string) string {
	return fmt.Sprintf(
		"[%d](%s)",
		p.Count,
		fmt.Sprintf("%s/pulls?q=is:pr+author:%s", p.Url, username),
	)
}

func GenerateProfile(username string) ([]Profile, []PrRepository) {
	var (
		wg    sync.WaitGroup
		items = make(chan map[string]interface{}, 2)
	)

	wg.Add(2)
	go func(items chan map[string]interface{}) {
		defer wg.Done()
		allRepository := github.FetchAllRepository(username)
		profiles := makeProfile(allRepository)

		items <- map[string]interface{}{"profiles": profiles}
	}(items)
	go func(items chan map[string]interface{}) {
		defer wg.Done()
		allIssues := github.FetchAllPrs(username)
		issues := makeIssues(allIssues, username)

		items <- map[string]interface{}{"issues": issues}
	}(items)

	wg.Wait()
	close(items)

	var (
		profiles []Profile
		issues   []PrRepository
	)
	for item := range items {
		if p, ok := item["profiles"]; ok {
			profiles = p.([]Profile)
		}
		if i, ok := item["issues"]; ok {
			issues = i.([]PrRepository)
		}
	}
	return profiles, issues
}

func makeProfile(repos []*github2.Repository) (profile []Profile) {

	var wg sync.WaitGroup
	for _, repo := range repos {
		if *repo.Fork == true {
			continue
		}
		repo := repo
		wg.Add(1)
		go func() {
			defer wg.Done()

			var description string
			if len(repo.GetDescription()) > 0 {
				description = Translate(repo.GetDescription())
			}

			profile = append(profile, Profile{
				Name:        repo.GetName(),
				HTMLUrl:     repo.GetHTMLURL(),
				Description: description,
				Stars:       repo.GetStargazersCount(),
				Forks:       repo.GetForksCount(),
				Create:      (*repo.CreatedAt).String()[:10],
				Update:      (*repo.UpdatedAt).String()[:10],
				Language:    repo.GetLanguage(),
			})
		}()
	}

	wg.Wait()

	sort.Slice(profile[:], func(i, j int) bool {
		return profile[j].Stars < profile[i].Stars
	})
	return
}

func getAllPrLinks(repositoryURL string) string {
	return strings.ReplaceAll(repositoryURL, "api.github.com/repos", "github.com")
}

func makeIssues(issues []*github2.Issue, username string) []PrRepository {
	var (
		wg           sync.WaitGroup
		prRepository []PrRepository
	)

	for _, issue := range issues {
		owner := strings.Split(*issue.RepositoryURL, "/")[4]
		if strings.Compare(owner, username) == 0 {
			continue
		}

		issue := issue
		wg.Add(1)
		go func() {
			defer wg.Done()

			var description, language string
			link := getAllPrLinks(*issue.RepositoryURL)
			repository, err := github.FetchRepository(link)
			if err == nil {
				description = Translate(repository.GetDescription())
				language = repository.GetLanguage()
			}

			prRepository = append(prRepository, PrRepository{
				Name:        *issue.RepositoryURL,
				Url:         link,
				Description: description,
				Language:    language,
				State:       *issue.State,
				Created:     (*issue.CreatedAt).String()[:10],
				Count:       1,
			})
		}()
	}

	wg.Wait()
	newPrRepository := make(map[string]PrRepository)
	for _, pr := range prRepository {
		if _, ok := newPrRepository[pr.Name]; !ok {
			newPrRepository[pr.Name] = pr
		} else {
			oldPrRepository := newPrRepository[pr.Name]
			oldPrRepository.Count += 1
			newPrRepository[pr.Name] = oldPrRepository
		}
	}

	prRepository = prRepository[:0]
	for _, pr := range newPrRepository {
		prRepository = append(prRepository, pr)
	}
	sort.Slice(prRepository[:], func(i, j int) bool {
		return prRepository[j].Count < prRepository[i].Count
	})

	return prRepository
}
