package metrics

import (
	"context"
	"fmt"
	"github.com/CanobbioE/gf-metrics/pkg/client"
	"time"
)

const (
	LeadTimeForChangesCommand = "ltc"

	// repositories
	gopkgRepo            = "gopkg"
	automationRepo       = "glofox-automation-tests"
	appInfraRepo         = "gf-app-infra"
	cloudWatchAlarmsRepo = "gf-cloudwatch-alarms"
	tsMonorepo           = "glofox"
	coreRepo             = "api"
)

/*
	LTC(gopkg) = 29mins
	LTC(automation) = 20mins
	LTC(gf-app-infra) = 10mins
	LTC(gf-cloudwatch-alarms) = 10mins
	LTC(ts-monorepo) = 38mins
	LTC(coreRepo) = 23mins
*/

type LtcArgs struct {
	Ctx   context.Context
	Repos []string
}

// LTC lead time for changes: how long does it take to release changes to prod
type LTC struct {
	name   string
	cli    *client.AuthClient
	ltcMap map[string]time.Duration
}

func NewLeadTimeForChangesCmd(cli *client.AuthClient) Command {
	return &LTC{
		name: LeadTimeForChangesCommand,
		cli:  cli,
		ltcMap: map[string]time.Duration{
			gopkgRepo:            29 * time.Minute,
			automationRepo:       20 * time.Minute,
			appInfraRepo:         10 * time.Minute,
			cloudWatchAlarmsRepo: 10 * time.Minute,
			tsMonorepo:           38 * time.Minute,
			coreRepo:             23 * time.Minute,
		},
	}
}

func (l *LTC) Run(args CommandArgs) CommandReturnValues {
	a, _ := args.(*LtcArgs)
	var sum float64
	for _, repo := range a.Repos {
		sum += l.ltcMap[repo].Minutes()
	}

	ltc := sum / float64(len(a.Repos))
	fmt.Printf("Lead time for changes is: %00f\n", ltc)
	return nil
}

func (l *LTC) Name() string {
	return l.name
}
