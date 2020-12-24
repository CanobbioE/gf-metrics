package metrics

import (
	"context"
	"github.com/CanobbioE/gf-metrics/pkg/client"
	"log"
	"time"
)

const DeploymentFrequencyCommand = "deployment-frequency"

type DeploymentFrequencyArgs struct {
	Ctx                context.Context
	Repositories       []string
	Owner              string
	MetricStartingDate time.Time
	Team               []string
}

// DeploymentFrequency counts how many PR from the team (in the specified timeframe) reached master
type DeploymentFrequency struct {
	name string
	cli  *client.AuthClient
}

func NewDeploymentFrequencyCmd(client *client.AuthClient) Command {
	return &DeploymentFrequency{
		name: DeploymentFrequencyCommand,
		cli:  client,
	}
}

func (d *DeploymentFrequency) Run(args CommandArgs) {

	log.Println("Calculating deployment frequency...")
	a := args.(*DeploymentFrequencyArgs)
	deployPerRepo := make(map[string]int)

	for _, repo := range a.Repositories {

		count, err := d.cli.CountTeamDeploys(a.Ctx, &client.CountTeamPROpts{
			Team:                a.Team,
			State:               "closed",
			ToBranch:            "master",
			Repo:                repo,
			Organisation:        a.Owner,
			StartSearchFromDate: a.MetricStartingDate,
		})
		if err != nil {
			log.Fatalf("couldn't count deploys for team %v on repo %s", a.Team, repo)
		}

		deployPerRepo[repo] = count
	}

	printDeployFrequency(deployPerRepo)
}

func printDeployFrequency(deployPerRepo map[string]int) {
	log.Println("Deploys per repo:")
	var total int
	for k, v := range deployPerRepo {
		if v > 0 {
			total += v
			log.Printf("[%s]: %d\n", k, v)
		}
	}
	log.Printf("Total deployments: %d\n", total)
}

func (d *DeploymentFrequency) Name() string {
	return d.name
}
