package analysis

import (
	"fmt"
	"strings"

	"slowpreview/domain"
	"slowpreview/engine"
)

type Check struct {
	Name   string
	Passed bool
	Detail string
}

type Checklist struct {
	Checks []Check
}

func RunChecklist(asset domain.VideoAsset, spec domain.PreviewSpec) Checklist {
	checks := []Check{
		{Name: "asset dimensions", Passed: asset.Width >= 320 && asset.Height >= 180, Detail: fmt.Sprintf("%dx%d", asset.Width, asset.Height)},
		{Name: "crop bounds", Passed: spec.Crop.StartMS >= 0 && spec.Crop.EndMS <= asset.DurationMS && spec.Crop.DurationMS() > 0, Detail: fmt.Sprintf("%d-%dms", spec.Crop.StartMS, spec.Crop.EndMS)},
		{Name: "speed choice", Passed: spec.Speed.Valid(), Detail: spec.Speed.Label()},
		{Name: "resolution choice", Passed: spec.Resolution.Valid(), Detail: string(spec.Resolution)},
		{Name: "output path", Passed: strings.TrimSpace(spec.OutputPath) != "", Detail: spec.OutputPath},
	}
	plan := engine.BuildPlan(asset, spec)
	checks = append(checks, Check{Name: "quality warnings", Passed: len(plan.Warnings) == 0, Detail: strings.Join(plan.Warnings, "; ")})
	return Checklist{Checks: checks}
}

func (c Checklist) Passed() bool {
	for _, check := range c.Checks {
		if !check.Passed {
			return false
		}
	}
	return true
}

func (c Checklist) Failed() []Check {
	failed := make([]Check, 0)
	for _, check := range c.Checks {
		if !check.Passed {
			failed = append(failed, check)
		}
	}
	return failed
}

func (c Checklist) Summary() string {
	passed := 0
	for _, check := range c.Checks {
		if check.Passed {
			passed++
		}
	}
	return fmt.Sprintf("%d/%d checks passed", passed, len(c.Checks))
}

func FilterFailures(checks []Check) []string {
	issues := make([]string, 0)
	for _, check := range checks {
		if !check.Passed {
			issues = append(issues, check.Name+": "+check.Detail)
		}
	}
	return issues
}
