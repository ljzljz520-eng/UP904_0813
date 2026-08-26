package analysis

import (
	"fmt"
	"sort"

	"slowpreview/domain"
	"slowpreview/engine"
	"slowpreview/timeline"
)

type Recommendation struct {
	Speed       domain.PreviewSpeed
	Resolution  domain.Resolution
	Interpolate bool
	Reason      string
	Confidence  int
}

type Assessment struct {
	Score       int
	Label       string
	Warnings    []string
	Suggestions []string
}

func Recommend(asset domain.VideoAsset, crop domain.CropWindow) []Recommendation {
	recommendations := make([]Recommendation, 0)
	base := domain.PreviewSpec{ID: "recommendation", AssetID: asset.ID, Speed: domain.SpeedHalf, Resolution: domain.ResolutionMedium, Crop: crop}
	if asset.FrameRate >= 60 {
		recommendations = append(recommendations, Recommendation{Speed: domain.SpeedHalf, Resolution: domain.ResolutionMedium, Interpolate: false, Reason: "high frame rate supports clean half-speed motion", Confidence: 92})
	} else {
		recommendations = append(recommendations, Recommendation{Speed: domain.SpeedThreeQuarter, Resolution: domain.ResolutionMedium, Interpolate: true, Reason: "interpolation fills motion at a moderate slowdown", Confidence: 78})
	}
	if asset.Width >= 2560 && asset.Height >= 1440 {
		recommendations = append(recommendations, Recommendation{Speed: domain.SpeedHalf, Resolution: domain.ResolutionHigh, Interpolate: asset.FrameRate < 60, Reason: "source supports detailed crop review", Confidence: 84})
	}
	if crop.DurationMS() <= 5000 {
		recommendations = append(recommendations, Recommendation{Speed: domain.SpeedDouble, Resolution: domain.ResolutionLow, Interpolate: false, Reason: "short crop is useful for a quick motion check", Confidence: 65})
	}
	_ = base
	sort.SliceStable(recommendations, func(i, j int) bool { return recommendations[i].Confidence > recommendations[j].Confidence })
	return recommendations
}

func Assess(asset domain.VideoAsset, spec domain.PreviewSpec) Assessment {
	plan := engine.BuildPlan(asset, spec)
	score := plan.QualityScore()
	assessment := Assessment{Score: score, Label: scoreLabel(score), Warnings: append([]string(nil), plan.Warnings...), Suggestions: make([]string, 0)}
	if spec.Speed == domain.SpeedHalf && asset.FrameRate < 60 {
		assessment.Suggestions = append(assessment.Suggestions, "enable interpolation for smoother half-speed motion")
	}
	if spec.Resolution == domain.ResolutionHigh && asset.Width < 2560 {
		assessment.Suggestions = append(assessment.Suggestions, "use medium resolution to avoid upscaling")
	}
	if spec.Crop.DurationMS() > 20000 {
		assessment.Suggestions = append(assessment.Suggestions, "shorten the crop for a focused preview")
	}
	if len(assessment.Warnings) == 0 {
		assessment.Suggestions = append(assessment.Suggestions, "configuration is ready for review")
	}
	return assessment
}

func scoreLabel(score int) string {
	switch {
	case score >= 85:
		return "excellent"
	case score >= 70:
		return "strong"
	case score >= 50:
		return "balanced"
	default:
		return "needs attention"
	}
}

func ExplainRecommendation(recommendation Recommendation) string {
	return fmt.Sprintf("%s %s %s (%d%% confidence)", recommendation.Speed.Label(), recommendation.Resolution, recommendation.Reason, recommendation.Confidence)
}

func CompareSpeeds(asset domain.VideoAsset, crop domain.CropWindow) map[domain.PreviewSpeed]timeline.QualityMetrics {
	comparison := make(map[domain.PreviewSpeed]timeline.QualityMetrics)
	for _, speed := range []domain.PreviewSpeed{domain.SpeedHalf, domain.SpeedThreeQuarter, domain.SpeedNormal, domain.SpeedDouble} {
		spec := domain.PreviewSpec{ID: "comparison", AssetID: asset.ID, Speed: speed, Resolution: domain.ResolutionMedium, Crop: crop}
		comparison[speed] = timeline.CalculateMetrics(asset, spec)
	}
	return comparison
}
