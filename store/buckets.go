package store

const (
	assetsBucket = "VideoAsset"
	specsBucket  = "PreviewSpec"
	tasksBucket  = "RenderTask"
	eventsBucket = "ActivityEvent"
	metaBucket   = "LibrarySnapshot"
	sequenceKey  = "sequence"
)

var bucketNames = []string{assetsBucket, specsBucket, tasksBucket, eventsBucket, metaBucket}

func bucketForEntity(entity string) string {
	switch entity {
	case "VideoAsset":
		return assetsBucket
	case "PreviewSpec":
		return specsBucket
	case "RenderTask":
		return tasksBucket
	case "ActivityEvent":
		return eventsBucket
	case "LibrarySnapshot":
		return metaBucket
	default:
		return ""
	}
}

func isEntityBucket(name string) bool {
	return bucketForEntity(name) != ""
}
