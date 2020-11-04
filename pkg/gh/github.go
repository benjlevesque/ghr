package gh

import (
	"context"

	"github.com/google/go-github/v32/github"
)

// GetReposForOrg returns a list of repos for a given owner
func GetReposForOrg(owner string) ([]string, error) {
	gh := github.NewClient(nil)

	repos, _, err := gh.Repositories.List(context.TODO(), owner, &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: 100,
		},
	})
	if err != nil {
		return nil, err
	}

	result := make([]string, len(repos))
	for i, repo := range repos {
		result[i] = owner + "/" + *repo.Name
	}
	return result, nil
}

// GetReleaseByTag returns a Release for a given repository tag
// if tag is left empty, it will default to latest
func GetReleaseByTag(owner, repo, tag string) (*github.RepositoryRelease, error) {
	gh := github.NewClient(nil)

	var release *github.RepositoryRelease
	var err error
	if tag == "" || tag == "latest" {
		release, _, err = gh.Repositories.GetLatestRelease(context.TODO(), owner, repo)
	} else {
		release, _, err = gh.Repositories.GetReleaseByTag(context.TODO(), owner, repo, tag)
	}
	if err != nil {
		return nil, err
	}
	return release, nil
}
