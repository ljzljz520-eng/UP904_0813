package service

import (
	"testing"

	"slowpreview/domain"
	"slowpreview/store"
)

func TestWorkflowOne(t *testing.T) {
	database, err := store.Open(t.TempDir() + "/workflow-one.db")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	svc := New(store.NewRepository(database))
	asset := domain.VideoAsset{ID: "asset-one", SourcePath: "one.mp4", Title: "footwork", DurationMS: 16000, FrameRate: 60, Width: 1920, Height: 1080}
	if err := svc.RegisterAsset(asset); err != nil {
		t.Fatal(err)
	}
	draft, err := svc.DraftPreview(PreviewRequest{ID: "preview-one", AssetID: asset.ID, SpeedLabel: "0.75x", Crop: domain.CropWindow{StartMS: 1000, EndMS: 6000}, Resolution: domain.ResolutionMedium})
	if err != nil || draft.Status != domain.StatusDraft {
		t.Fatalf("draft: %#v %v", draft, err)
	}
	if err := svc.QueuePreview(draft.Spec.ID); err != nil {
		t.Fatal(err)
	}
	record, err := svc.GeneratePreview(draft.Spec.ID)
	if err != nil || record.Status != domain.StatusReady {
		t.Fatalf("generate: %#v %v", record, err)
	}
	inspected, err := svc.InspectPreview(draft.Spec.ID)
	if err != nil || inspected.Status != domain.StatusReady {
		t.Fatalf("inspect: %#v %v", inspected, err)
	}
}
