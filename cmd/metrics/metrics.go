package metrics

import (
	"context"
	"github.com/CanobbioE/gf-metrics/pkg/client"
	"log"
)

type Command interface {
	Run(CommandArgs) CommandReturnValues
	Name() string
}

type CommandReturnValues interface{}
type CommandArgs interface{}

type Args struct {
	DeploymentFrequency *DeploymentFrequencyArgs
}

type Metrics struct {
	token        string
	owner        string
	cli          *client.AuthClient
	commands     []Command
	commandsArgs map[string]CommandArgs
	repositories []string
}

func New(ctx context.Context, owner, token string) (*Metrics, error) {
	cli := client.NewWithOauth(ctx, token)

	log.Println("Fetching organisation repositories...")
	repos, err := cli.GetOrgRepos(ctx, owner)
	if err != nil {
		return nil, err
	}
	log.Println("Done fetching!")

	m := &Metrics{
		token:        token,
		owner:        owner,
		cli:          cli,
		repositories: repos,
	}

	return m, nil
}

func (m *Metrics) Run(args CommandArgs) CommandReturnValues {
	metricArgs := args.(Args)

	deployArgs := metricArgs.DeploymentFrequency
	deployArgs.Repositories = m.repositories

	repos := make([]string, 0)
	for repo, _ := range NewDeploymentFrequencyCmd(m.cli).Run(deployArgs).(map[string]int) {
		repos = append(repos, repo)
	}

	NewLeadTimeForChangesCmd(m.cli).Run(&LtcArgs{
		Ctx:   metricArgs.DeploymentFrequency.Ctx,
		Repos: repos,
	})

	return nil
}

func (m *Metrics) registerCommand(cmd Command, args CommandArgs) {
	m.commands = append(m.commands, cmd)

	m.commandsArgs[cmd.Name()] = args
}

// Mean Time to Detection: (How long does it take us to find bugs)

// Time to restore service (How long does it take us to fix issues)

// Change Failure Rate (What percentage of this weeks deployments had issues):

// Unplanned Work

// Average WIP During Week:
