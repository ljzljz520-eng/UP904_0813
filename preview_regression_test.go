package slowpreview

import (
	"strings"
	"testing"

	"slowpreview/domain"
	"slowpreview/service"
	"slowpreview/store"
)

func TestPreviewUsesSelectedSpeed(t *testing.T) {
	database, err := store.Open(t.TempDir() + "/preview.db")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	svc := service.New(store.NewRepository(database))
	asset := domain.VideoAsset{ID: "asset-regression", SourcePath: "training.mp4", Title: "jump drill", DurationMS: 12000, FrameRate: 60, Width: 1920, Height: 1080}
	if err := svc.RegisterAsset(asset); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.DraftPreview(service.PreviewRequest{ID: "preview-regression", AssetID: asset.ID, SpeedLabel: "0.5x", Resolution: domain.ResolutionMedium, Crop: domain.CropWindow{StartMS: 1000, EndMS: 5000}}); err != nil {
		t.Fatal(err)
	}
	record, err := svc.GeneratePreview("preview-regression")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(record.Task.Command, "setpts=2.00x*PTS") {
		t.Fatalf("selected half speed missing from command: %s", record.Task.Command)
	}
}
