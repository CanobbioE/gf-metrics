package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/CanobbioE/gf-metrics/cfg"
	"github.com/CanobbioE/gf-metrics/cmd/metrics"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	tokenUsage     = `a GitHub token with repo permissions`
	ownerUsage     = `the owner of all the repositories, usually an organisation`
	configUsage    = `the path to a JSON file defining all or part of the options`
	teamUsage      = `a space separated list of GitHub handlers`
	timestampUsage = `the timestamp (in seconds) from which the metrics will start, defaults to one week ago`
	helpUsage      = `print this helpful message`
)

func getUsage() string {
	s := `gf-metrics (BETA) is a tool that partially automates the job of retrieving a team metrics.

	Usage:
		gf-metrics [--token <token> --owner <repository owner> --team <team members>]  [--settings <path to JSON file>]

	Options:
		--token <token>:
			%s
		--owner <owner>:
			%s
		--config <path to  config>:
			%s (see below)
		--team <team members>:
			%s
		--timestamp <timestamp>:
			%s
		--help, -h:
			%s

	Example:
		gf-metrics --config ./settings.json --token SecretToken123

	Settings:
		the JSON file can define one or more of the following fields:
			- "token": string
			- "team": array of string
			- "owner": string
			- "timestamp": string

		Each required field left blank must be provided through the CLI flags.
		Each value specified as a CLI flag will take priority over the configuration file.
		You can find a configuration file example at https://github.com/CanobbioE/gf-metrics
`
	return fmt.Sprintf(s, tokenUsage, ownerUsage, configUsage, teamUsage, timestampUsage, helpUsage)
}

func main() {

	ctx := context.Background()
	oneWeekAgo := strconv.Itoa(int(time.Now().Add(-1 * 7 * 24 * time.Hour).Unix()))

	var githubOauthToken, repoOwner, cfgPath, teamMembers, startDateTimestamp string
	var help bool
	flag.StringVar(&githubOauthToken, "token", "", tokenUsage)
	flag.StringVar(&repoOwner, "owner", "", ownerUsage)
	flag.StringVar(&cfgPath, "config", "", configUsage)
	flag.StringVar(&teamMembers, "team", "", teamUsage)
	flag.StringVar(&startDateTimestamp, "timestamp", "", timestampUsage)
	flag.BoolVar(&help, "help", false, helpUsage)

	flag.Usage = func() { fmt.Println(getUsage()) }

	flag.Parse()

	// handle help
	if help {
		flag.Usage()
		os.Exit(0)
	}

	// handle team members
	var team []string
	if !empty(teamMembers) {
		team = strings.Split(teamMembers, " ")
	}

	// handle config file
	if !empty(cfgPath) {
		config, err := cfg.FromFile(cfgPath)
		if err != nil {
			log.Fatalf("couldn't read file %s: %v", cfgPath, err)
		}

		config.FillEmpty(&githubOauthToken, &repoOwner, &team, &startDateTimestamp)
	}

	// handle required
	if (empty(githubOauthToken) || empty(repoOwner)) && empty(cfgPath) {
		log.Fatalln("expected GitHub oauth token and repository owner or a settings file")
	}

	// handle timestamp
	if empty(startDateTimestamp) {
		startDateTimestamp = oneWeekAgo
	}

	t, err := strconv.ParseInt(startDateTimestamp, 10, 64)
	if err != nil {
		log.Fatalf("couldn't parse timestamp %v: %v", startDateTimestamp, err)
	}

	metricsCmd, err := metrics.New(ctx, repoOwner, githubOauthToken)
	if err != nil {
		log.Fatalf("couldn't initialize the metrics command: %v", err)
	}

	metricsCmd.Run(metrics.Args{DeploymentFrequency: &metrics.DeploymentFrequencyArgs{
		Ctx:                ctx,
		Owner:              repoOwner,
		MetricStartingDate: time.Unix(t, 0),
		Team:               team,
	}})

}

func empty(s string) bool {
	return s == ""
}
