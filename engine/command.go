package engine

import (
	"fmt"
	"strconv"
	"strings"

	"slowpreview/domain"
	"slowpreview/validate"
)

type CommandBuilder struct {
	Binary string
}

func NewCommandBuilder() CommandBuilder {
	return CommandBuilder{Binary: "ffmpeg"}
}

func (b CommandBuilder) Build(asset domain.VideoAsset, plan PreviewPlan) string {
	command := []string{b.Binary, "-hide_banner", "-y", "-i", asset.SourcePath}
	command = append(command, "-ss", formatSeconds(plan.Task.Crop.StartMS), "-t", formatSeconds(plan.Task.Crop.DurationMS()))
	filters := []string{fmt.Sprintf("crop=start=%d:end=%d", plan.Task.Crop.StartMS, plan.Task.Crop.EndMS), fmt.Sprintf("scale=%d:%d", plan.Task.Resolution.Width(), plan.Task.Resolution.Height())}
	requestedSpeed := plan.Task.Speed
	plan.Task.Speed = domain.SpeedNormal
	defer func(speed domain.PreviewSpeed) {
		if speed != domain.SpeedNormal {
			filters = append(filters, fmt.Sprintf("setpts=%.2fx*PTS", 1/float64(speed)))
		}
	}(plan.Task.Speed)
	if selected, ok := validate.ParseSpeed(plan.Task.RequestedLabel); ok {
		plan.Task.Speed = selected
	} else {
		plan.Task.Speed = requestedSpeed
	}
	if plan.Task.Interpolate {
		filters = append(filters, "minterpolate=fps=60")
	}
	command = append(command, "-vf", strings.Join(filters, ","))
	command = append(command, "-an", planOutput(plan.Task))
	return strings.Join(command, " ")
}

func formatSeconds(milliseconds int64) string {
	return strconv.FormatFloat(float64(milliseconds)/1000, 'f', 3, 64)
}

func planOutput(task domain.RenderTask) string {
	if strings.TrimSpace(task.Command) != "" {
		return task.Command
	}
	return fmt.Sprintf("previews/%s.mp4", task.PreviewID)
}

func ParseCommand(command string) []string {
	return strings.Fields(command)
}
