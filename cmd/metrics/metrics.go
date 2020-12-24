package metrics

import (
	"context"
	"github.com/CanobbioE/gf-metrics/pkg/client"
	"log"
	"time"
)

type Command interface {
	Run(CommandArgs)
	Name() string
}

type CommandArgs interface{}

type Metrics struct {
	token        string
	owner        string
	cli          *client.AuthClient
	commands     map[Command]CommandArgs
	repositories []string
}

func New(ctx context.Context, owner, token string, team []string, metricStartingDate time.Time) (*Metrics, error) {
	cli := client.NewWithOauth(ctx, token)

	log.Println("Fetching organisation repositories...")
	repos, err := cli.GetOrgRepos(ctx, owner)
	if err != nil {
		return nil, err
	}
	log.Println("Done fetching!")

	// deployment frequency
	deployFreqCmd := NewDeploymentFrequencyCmd(cli)
	deployFreqArgs := &DeploymentFrequencyArgs{
		Ctx:                ctx,
		Repositories:       repos,
		Owner:              owner,
		MetricStartingDate: metricStartingDate,
		Team:               team,
	}

	commands := map[Command]CommandArgs{
		deployFreqCmd: deployFreqArgs,
	}

	return &Metrics{
		token:        token,
		owner:        owner,
		cli:          cli,
		commands:     commands,
		repositories: repos,
	}, nil
}

func (m *Metrics) Run() {
	for cmd, cmdArgs := range m.commands {
		cmd.Run(cmdArgs)
	}
}
