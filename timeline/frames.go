package timeline

import "slowpreview/domain"

type FramePoint struct {
	SourceMS int64
	OutputMS int64
	Index    int
}

func FramePoints(spec domain.PreviewSpec, asset domain.VideoAsset) []FramePoint {
	if asset.FrameRate <= 0 || spec.Crop.DurationMS() <= 0 || spec.Speed <= 0 {
		return nil
	}
	frameMS := int64(1000 / asset.FrameRate)
	if frameMS <= 0 {
		frameMS = 1
	}
	points := make([]FramePoint, 0)
	for source := spec.Crop.StartMS; source <= spec.Crop.EndMS; source += frameMS {
		output := int64(float64(source-spec.Crop.StartMS) / float64(spec.Speed))
		points = append(points, FramePoint{SourceMS: source, OutputMS: output, Index: len(points)})
		if len(points) >= 2000 {
			break
		}
	}
	return points
}

func KeyFrames(points []FramePoint, every int) []FramePoint {
	if every < 1 {
		every = 1
	}
	keys := make([]FramePoint, 0)
	for i, point := range points {
		if i%every == 0 || i == len(points)-1 {
			keys = append(keys, point)
		}
	}
	return keys
}

func EstimateFrames(asset domain.VideoAsset, crop domain.CropWindow) int {
	if asset.FrameRate <= 0 || crop.DurationMS() <= 0 {
		return 0
	}
	return int(float64(crop.DurationMS()) * float64(asset.FrameRate) / 1000)
}
