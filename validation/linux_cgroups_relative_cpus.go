package main

import (
	"github.com/mndrix/tap-go"
	rspec "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/opencontainers/runtime-tools/cgroups"
	"github.com/opencontainers/runtime-tools/validation/util"
)

func main() {
	var shares uint64 = 1024
	var period uint64 = 100000
	var quota int64 = 50000
	var cpus, mems string = "0-1", "0"
	g := util.GetDefaultGenerator()
	g.SetLinuxCgroupsPath(cgroups.RelCgroupPath)
	g.SetLinuxResourcesCPUShares(shares)
	g.SetLinuxResourcesCPUQuota(quota)
	g.SetLinuxResourcesCPUPeriod(period)
	g.SetLinuxResourcesCPUCpus(cpus)
	g.SetLinuxResourcesCPUMems(mems)
	err := util.RuntimeOutsideValidate(g, func(config *rspec.Spec, state *rspec.State) error {
		t := tap.New()
		t.Header(0)

		cg, err := cgroups.FindCgroup()
		t.Ok((err == nil), "find cpus cgroup")
		if err != nil {
			t.Diagnostic(err.Error())
			t.AutoPlan()
			return nil
		}

		lcd, err := cg.GetCPUData(state.Pid, config.Linux.CgroupsPath)
		t.Ok((err == nil), "get cpus cgroup data")
		if err != nil {
			t.Diagnostic(err.Error())
			t.AutoPlan()
			return nil
		}

		t.Ok(*lcd.Shares == shares, "cpus shares limit is set correctly")
		t.Diagnosticf("expect: %d, actual: %d", shares, lcd.Shares)

		t.Ok(*lcd.Quota == quota, "cpus quota is set correctly")
		t.Diagnosticf("expect: %d, actual: %d", quota, lcd.Quota)

		t.Ok(*lcd.Period == period, "cpus period is set correctly")
		t.Diagnosticf("expect: %d, actual: %d", period, lcd.Period)

		t.Ok(lcd.Cpus == cpus, "cpus cpus is set correctly")
		t.Diagnosticf("expect: %s, actual: %s", cpus, lcd.Cpus)

		t.Ok(lcd.Mems == mems, "cpus mems is set correctly")
		t.Diagnosticf("expect: %s, actual: %s", mems, lcd.Mems)

		t.AutoPlan()
		return nil
	})

	if err != nil {
		util.Fatal(err)
	}
}
