package domain

import (
	"fmt"
	"strings"
)

type PreviewSpeed float64

const (
	SpeedHalf         PreviewSpeed = 0.5
	SpeedThreeQuarter PreviewSpeed = 0.75
	SpeedNormal       PreviewSpeed = 1
	SpeedDouble       PreviewSpeed = 2
)

func (s PreviewSpeed) Valid() bool {
	return s == SpeedHalf || s == SpeedThreeQuarter || s == SpeedNormal || s == SpeedDouble
}

func (s PreviewSpeed) Label() string {
	switch s {
	case SpeedHalf:
		return "0.5x"
	case SpeedThreeQuarter:
		return "0.75x"
	case SpeedDouble:
		return "2x"
	default:
		return "1x"
	}
}

type Resolution string

const (
	ResolutionLow    Resolution = "720p"
	ResolutionMedium Resolution = "1080p"
	ResolutionHigh   Resolution = "1440p"
)

func (r Resolution) Valid() bool {
	return r == ResolutionLow || r == ResolutionMedium || r == ResolutionHigh
}

func (r Resolution) Width() int {
	switch r {
	case ResolutionLow:
		return 1280
	case ResolutionHigh:
		return 2560
	default:
		return 1920
	}
}

func (r Resolution) Height() int {
	switch r {
	case ResolutionLow:
		return 720
	case ResolutionHigh:
		return 1440
	default:
		return 1080
	}
}

type CropWindow struct {
	StartMS int64 `json:"start_ms"`
	EndMS   int64 `json:"end_ms"`
}

func (c CropWindow) DurationMS() int64 {
	return c.EndMS - c.StartMS
}

func (c CropWindow) Contains(ms int64) bool {
	return ms >= c.StartMS && ms <= c.EndMS
}

type VideoAsset struct {
	ID          string `json:"id"`
	SourcePath  string `json:"source_path"`
	Title       string `json:"title"`
	DurationMS  int64  `json:"duration_ms"`
	FrameRate   int    `json:"frame_rate"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	CreatedStep int    `json:"created_step"`
}

func (v VideoAsset) Valid() bool {
	return strings.TrimSpace(v.ID) != "" && strings.TrimSpace(v.SourcePath) != "" && v.DurationMS > 0 && v.FrameRate > 0 && v.Width > 0 && v.Height > 0
}

func (v VideoAsset) AspectRatio() float64 {
	if v.Height == 0 {
		return 0
	}
	return float64(v.Width) / float64(v.Height)
}

type PreviewSpec struct {
	ID             string       `json:"id"`
	AssetID        string       `json:"asset_id"`
	Speed          PreviewSpeed `json:"speed"`
	Resolution     Resolution   `json:"resolution"`
	Crop           CropWindow   `json:"crop"`
	Interpolate    bool         `json:"interpolate"`
	OutputPath     string       `json:"output_path"`
	RequestedLabel string       `json:"requested_label"`
}

func (s PreviewSpec) DurationMS() int64 {
	if s.Speed <= 0 {
		return 0
	}
	return int64(float64(s.Crop.DurationMS()) / float64(s.Speed))
}

func (s PreviewSpec) Summary() string {
	interpolation := "off"
	if s.Interpolate {
		interpolation = "on"
	}
	return fmt.Sprintf("%s %s crop=%d-%dms interpolation=%s", s.Speed.Label(), s.Resolution, s.Crop.StartMS, s.Crop.EndMS, interpolation)
}

type RenderTask struct {
	ID             string       `json:"id"`
	PreviewID      string       `json:"preview_id"`
	Speed          PreviewSpeed `json:"speed"`
	RequestedLabel string       `json:"requested_label"`
	Resolution     Resolution   `json:"resolution"`
	Crop           CropWindow   `json:"crop"`
	Interpolate    bool         `json:"interpolate"`
	Command        string       `json:"command"`
	Status         string       `json:"status"`
}

func (t RenderTask) Finished() bool {
	return t.Status == "ready" || t.Status == "failed"
}

type PreviewRecord struct {
	Spec       PreviewSpec `json:"spec"`
	Task       RenderTask  `json:"task"`
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	CreatedSeq int64       `json:"created_seq"`
	UpdatedSeq int64       `json:"updated_seq"`
}

func (p PreviewRecord) IsTerminal() bool {
	return p.Status == "ready" || p.Status == "failed" || p.Status == "archived"
}

type ActivityEvent struct {
	Seq       int64  `json:"seq"`
	PreviewID string `json:"preview_id"`
	Kind      string `json:"kind"`
	Message   string `json:"message"`
}

type LibrarySnapshot struct {
	Assets   []VideoAsset    `json:"assets"`
	Previews []PreviewRecord `json:"previews"`
	Events   []ActivityEvent `json:"events"`
}
