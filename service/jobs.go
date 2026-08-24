package service

import (
	"fmt"

	"slowpreview/domain"
	"slowpreview/engine"
)

type JobView struct {
	ID          string
	Status      string
	SpeedLabel  string
	OutputPath  string
	WarningText string
}

func BuildJobView(record domain.PreviewRecord, warnings []string) JobView {
	warningText := "none"
	if len(warnings) > 0 {
		warningText = fmt.Sprintf("%d warning(s): %s", len(warnings), warnings[0])
	}
	return JobView{ID: record.Spec.ID, Status: record.Status, SpeedLabel: record.Spec.Speed.Label(), OutputPath: record.Spec.OutputPath, WarningText: warningText}
}

func PlanJob(asset domain.VideoAsset, spec domain.PreviewSpec) (engine.PreviewPlan, error) {
	if !asset.Valid() {
		return engine.PreviewPlan{}, fmt.Errorf("cannot plan invalid asset %s", asset.ID)
	}
	return engine.BuildPlan(asset, spec), nil
}

func Retryable(status string) bool {
	return status == domain.StatusFailed
}

func NextAction(status string) string {
	switch status {
	case domain.StatusDraft:
		return "queue"
	case domain.StatusQueued:
		return "generate"
	case domain.StatusReady:
		return "archive"
	case domain.StatusFailed:
		return "retry"
	default:
		return "inspect"
	}
}
