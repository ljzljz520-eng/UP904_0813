package store

import (
	"testing"

	"slowpreview/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/library.db"
	database, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	repository := NewRepository(database)
	asset := domain.VideoAsset{ID: "asset-persist", SourcePath: "persist.mp4", Title: "persistence drill", DurationMS: 9000, FrameRate: 60, Width: 1920, Height: 1080}
	if err := repository.SaveAsset(asset); err != nil {
		t.Fatal(err)
	}
	if _, err := repository.SaveSpecAndTaskForTest(); err != nil {
		t.Fatal(err)
	}
	if err := database.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	reloaded, err := NewRepository(reopened).GetAsset(asset.ID)
	if err != nil || reloaded.Title != asset.Title {
		t.Fatalf("reloaded=%#v err=%v", reloaded, err)
	}
}

func (r *Repository) SaveSpecAndTaskForTest() (domain.PreviewSpec, error) {
	spec := domain.PreviewSpec{ID: "persist-preview", AssetID: "asset-persist", Speed: domain.SpeedHalf, Resolution: domain.ResolutionMedium, Crop: domain.CropWindow{StartMS: 0, EndMS: 4000}, OutputPath: "previews/persist.mp4", RequestedLabel: "0.5x"}
	if err := r.SaveSpec(spec); err != nil {
		return spec, err
	}
	return spec, r.SaveTask(domain.RenderTask{ID: "task-persist-preview", PreviewID: spec.ID, Speed: spec.Speed, Status: domain.StatusDraft})
}
