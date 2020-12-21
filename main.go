package main

import (
	"context"
	"github.com/CanobbioE/gf-metrics/client"
	"log"
	"time"
)

const (
	GlofoxOrgOwner = "glofoxinc"
	FromBranch     = "develop"
	ToBranch       = "master"
)

func getPlatformSquadTeam() []string {
	return []string{"CanobbioE", "cristianpontes", "Gjergj"}
}

func main() {
	ctx := context.Background()
	oneWeekAgo := time.Now().Add(-7 * time.Hour * 24)
	deployPerRepo := make(map[string]int)

	c := client.NewWithOauth(ctx, "<access token>")

	log.Printf("Fetching repos from %s\n", GlofoxOrgOwner)
	repos, err := c.GetOrgRepos(ctx, GlofoxOrgOwner)
	if err != nil {
		panic(err)
	}

	log.Printf("Repo fetched. Counting deploys to %v starting from %v...\n", ToBranch, oneWeekAgo)
	for _, repo := range repos {
		count, err := c.CountTeamDeploys(ctx, &client.CountTeamPROpts{
			Team:                getPlatformSquadTeam(),
			ToBranch:            ToBranch,
			Repo:                repo,
			Organisation:        GlofoxOrgOwner,
			StartSearchFromDate: oneWeekAgo,
		})
		if err != nil {
			panic(err)
		}
		deployPerRepo[repo] = count
	}

	log.Println("Deploys per repo:")
	for k, v := range deployPerRepo {
		if v > 0 {
			log.Printf("%s: %d\n", k, v)
		}

	}

}
