package client

import (
	"context"
	"fmt"
	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
	"strings"
	"time"
)

type WithAuth struct {
	cli *github.Client
}

func NewWithOauth(ctx context.Context, token string) *WithAuth {
	return &WithAuth{
		cli: github.NewClient(
			oauth2.NewClient(
				ctx,
				oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token},
				),
			),
		),
	}
}

type CountTeamPROpts struct {
	Team                []string
	FromBranch          string
	ToBranch            string
	Repo                string
	Organisation        string
	StartSearchFromDate time.Time
}

func (c *WithAuth) CountTeamDeploys(ctx context.Context, opts *CountTeamPROpts) (int, error) {
	var count int
	prs, _, err := c.cli.PullRequests.List(ctx, opts.Organisation, opts.Repo, &github.PullRequestListOptions{
		State:       "closed",
		Head:        opts.FromBranch,
		Base:        opts.ToBranch,
	})

	if err != nil {
		return count, err
	}

	for _, pr := range prs {
		if isPrFromTeam(pr, opts.Team) && pr.ClosedAt.After(opts.StartSearchFromDate){
			count++
		}
	}
	return count, nil
}

func isPrFromTeam(pr *github.PullRequest, team []string) bool {

	teamIncludeUser := func() bool {
		for _, member := range team {
			if strings.EqualFold(member, *pr.User.Login) {
				return true
			}
		}
		return false
	}

	return pr.User != nil && pr.User.Login != nil && teamIncludeUser()
}

func (c *WithAuth) GetOrgRepos(ctx context.Context, organisation string) ([]string, error) {
	names := make([]string, 0)
	query := fmt.Sprintf("user:%s", organisation)

	repos, _, err := c.cli.Search.Repositories(ctx, query, nil)
	if err != nil {
		return names, err
	}

	for _, repo := range repos.Repositories {
		names = append(names, *repo.Name)
	}
	return names, nil
}
