package engine

import (
	"fmt"
	"strings"

	"slowpreview/domain"
)

type RenderReceipt struct {
	TaskID     string
	Command    string
	OutputPath string
	Status     string
	Message    string
}

func SimulateRender(task domain.RenderTask, command string) RenderReceipt {
	if strings.TrimSpace(command) == "" {
		return RenderReceipt{TaskID: task.ID, Status: domain.StatusFailed, Message: "render command is empty"}
	}
	if task.Crop.DurationMS() <= 0 {
		return RenderReceipt{TaskID: task.ID, Command: command, Status: domain.StatusFailed, Message: "crop duration is empty"}
	}
	if !task.Speed.Valid() {
		return RenderReceipt{TaskID: task.ID, Command: command, Status: domain.StatusFailed, Message: "speed is unsupported"}
	}
	output := planOutput(task)
	return RenderReceipt{TaskID: task.ID, Command: command, OutputPath: output, Status: domain.StatusReady, Message: fmt.Sprintf("preview ready at %s", output)}
}

func StatusMessage(receipt RenderReceipt) string {
	if receipt.Status == domain.StatusReady {
		return receipt.Message
	}
	return "preview failed: " + receipt.Message
}

func CommandHasSpeed(command string, speed domain.PreviewSpeed) bool {
	needle := fmt.Sprintf("%.2fx*PTS", 1/float64(speed))
	return strings.Contains(command, needle)
}
