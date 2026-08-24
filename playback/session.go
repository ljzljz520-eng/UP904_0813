package playback

import (
	"fmt"
	"sort"

	"slowpreview/domain"
	"slowpreview/timeline"
)

type Session struct {
	Asset       domain.VideoAsset
	Spec        domain.PreviewSpec
	Points      []timeline.FramePoint
	CurrentMS   int64
	Playing     bool
	LoopEnabled bool
}

type Cue struct {
	AtMS   int64
	Label  string
	Detail string
}

func NewSession(asset domain.VideoAsset, spec domain.PreviewSpec) Session {
	points := timeline.FramePoints(spec, asset)
	return Session{Asset: asset, Spec: spec, Points: points, CurrentMS: 0, Playing: false, LoopEnabled: true}
}

func (s *Session) Play() string {
	if len(s.Points) == 0 {
		return "no frames available"
	}
	s.Playing = true
	return "playing"
}

func (s *Session) Pause() string {
	s.Playing = false
	return "paused"
}

func (s *Session) Seek(outputMS int64) int64 {
	if outputMS < 0 {
		outputMS = 0
	}
	duration := s.Spec.DurationMS()
	if outputMS > duration {
		outputMS = duration
	}
	s.CurrentMS = outputMS
	return s.CurrentMS
}

func (s Session) CurrentFrame() (timeline.FramePoint, bool) {
	if len(s.Points) == 0 {
		return timeline.FramePoint{}, false
	}
	best := s.Points[0]
	for _, point := range s.Points {
		if point.OutputMS > s.CurrentMS {
			break
		}
		best = point
	}
	return best, true
}

func (s Session) Progress() float64 {
	duration := s.Spec.DurationMS()
	if duration <= 0 {
		return 0
	}
	return float64(s.CurrentMS) / float64(duration)
}

func (s Session) Label() string {
	frame, ok := s.CurrentFrame()
	if !ok {
		return "empty preview"
	}
	return fmt.Sprintf("frame %d at %dms", frame.Index, frame.OutputMS)
}

func (s Session) Cues(interval int64) []Cue {
	if interval <= 0 {
		interval = 1000
	}
	cues := make([]Cue, 0)
	for at := int64(0); at <= s.Spec.DurationMS(); at += interval {
		cues = append(cues, Cue{AtMS: at, Label: formatTimestamp(at), Detail: cueDetail(at, s.Spec)})
	}
	if len(cues) == 0 || cues[len(cues)-1].AtMS != s.Spec.DurationMS() {
		at := s.Spec.DurationMS()
		cues = append(cues, Cue{AtMS: at, Label: formatTimestamp(at), Detail: cueDetail(at, s.Spec)})
	}
	return cues
}

func cueDetail(at int64, spec domain.PreviewSpec) string {
	if at == 0 {
		return "preview start"
	}
	if at >= spec.DurationMS() {
		return "preview end"
	}
	return fmt.Sprintf("slow-motion position %dms", at)
}

func formatTimestamp(value int64) string {
	seconds := value / 1000
	return fmt.Sprintf("%02d:%02d.%03d", seconds/60, seconds%60, value%1000)
}

func (s Session) SortedCues(interval int64) []Cue {
	cues := s.Cues(interval)
	sort.SliceStable(cues, func(i, j int) bool { return cues[i].AtMS < cues[j].AtMS })
	return cues
}
