package service

import (
	"testing"

	"slowpreview/catalog"
	"slowpreview/domain"
	"slowpreview/store"
)

func TestWorkflowTwo(t *testing.T) {
	database, err := store.Open(t.TempDir() + "/workflow-two.db")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	svc := New(store.NewRepository(database))
	lines := []string{"asset-two|session-two.mp4|serve practice|14000|30|1280|720", "asset-three|session-three.mp4|drill three|18000|60|1920|1080"}
	imported, issues, err := svc.ImportManifest(lines)
	if err != nil || imported != 2 || len(issues) != 0 {
		t.Fatalf("imported=%d issues=%v err=%v", imported, issues, err)
	}
	if _, err := svc.DraftPreview(PreviewRequest{ID: "preview-two", AssetID: "asset-two", SpeedLabel: "half", Crop: domain.CropWindow{StartMS: 0, EndMS: 4000}, Interpolate: true}); err != nil {
		t.Fatal(err)
	}
	view, err := svc.BuildLibraryView("serve")
	if err != nil {
		t.Fatal(err)
	}
	if len(view.AssetMatches) != 1 || len(view.ReviewQueue.Items) != 1 {
		t.Fatalf("unexpected view: %#v", view)
	}
	results, err := svc.SearchPreviews(PreviewQueryForTest())
	if err != nil || len(results) != 1 {
		t.Fatalf("search: %#v %v", results, err)
	}
}

func PreviewQueryForTest() catalog.PreviewQuery {
	return catalog.PreviewQuery{Statuses: []string{domain.StatusDraft}}
}
