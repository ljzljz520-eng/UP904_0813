package timeline

import (
	"math"

	"slowpreview/domain"
)

type QualityMetrics struct {
	InputFrames  int
	OutputFrames int
	InputMS      int64
	OutputMS     int64
	Compression  float64
	Sharpness    string
}

func CalculateMetrics(asset domain.VideoAsset, spec domain.PreviewSpec) QualityMetrics {
	input := EstimateFrames(asset, spec.Crop)
	output := int(math.Round(float64(input) / float64(spec.Speed)))
	if spec.Interpolate {
		output *= 2
	}
	sharpness := "balanced"
	if spec.Resolution == domain.ResolutionLow {
		sharpness = "compact"
	} else if spec.Resolution == domain.ResolutionHigh {
		sharpness = "detailed"
	}
	return QualityMetrics{InputFrames: input, OutputFrames: output, InputMS: spec.Crop.DurationMS(), OutputMS: spec.DurationMS(), Compression: compressionRatio(input, output), Sharpness: sharpness}
}

func compressionRatio(input, output int) float64 {
	if output == 0 {
		return 0
	}
	return float64(input) / float64(output)
}

func Score(metrics QualityMetrics) int {
	score := 50
	if metrics.Sharpness == "detailed" {
		score += 20
	} else if metrics.Sharpness == "compact" {
		score -= 10
	}
	if metrics.OutputFrames > metrics.InputFrames*3 {
		score -= 15
	} else if metrics.OutputFrames < metrics.InputFrames/2 {
		score += 10
	}
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}

func FrameRateLabel(asset domain.VideoAsset, interpolate bool) string {
	rate := asset.FrameRate
	if interpolate {
		rate *= 2
	}
	return string(rune('0'+rate/10)) + string(rune('0'+rate%10)) + " fps"
}
