package domain

const (
	EventAssetRegistered = "asset_registered"
	EventPreviewDrafted  = "preview_drafted"
	EventPreviewQueued   = "preview_queued"
	EventRenderStarted   = "render_started"
	EventRenderCompleted = "render_completed"
	EventRenderFailed    = "render_failed"
	EventPreviewArchived = "preview_archived"
)

func EventForStatus(status string) string {
	switch status {
	case StatusDraft:
		return EventPreviewDrafted
	case StatusQueued:
		return EventPreviewQueued
	case StatusRendering:
		return EventRenderStarted
	case StatusReady:
		return EventRenderCompleted
	case StatusFailed:
		return EventRenderFailed
	case StatusArchived:
		return EventPreviewArchived
	default:
		return "unknown"
	}
}

func EventDescription(event string) string {
	switch event {
	case EventAssetRegistered:
		return "video asset entered the library"
	case EventPreviewDrafted:
		return "preview configuration was saved"
	case EventPreviewQueued:
		return "preview was queued for rendering"
	case EventRenderStarted:
		return "render simulation started"
	case EventRenderCompleted:
		return "preview is ready to inspect"
	case EventRenderFailed:
		return "preview render failed validation"
	case EventPreviewArchived:
		return "preview was archived"
	default:
		return "unrecognized activity"
	}
}
