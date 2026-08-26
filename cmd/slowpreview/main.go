package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"slowpreview/domain"
	"slowpreview/report"
	"slowpreview/service"
	"slowpreview/store"
)

func main() {
	action := flag.String("action", "help", "register, draft, generate, inspect, history")
	dbPath := flag.String("db", "slowpreview.db", "bbolt database path")
	assetID := flag.String("asset", "training-01", "asset id")
	previewID := flag.String("preview", "preview-01", "preview id")
	speed := flag.String("speed", "0.5x", "preview speed")
	flag.Parse()
	db, err := store.Open(*dbPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()
	svc := service.New(store.NewRepository(db))
	if err := runAction(svc, *action, *assetID, *previewID, *speed); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runAction(svc *service.Service, action, assetID, previewID, speed string) error {
	switch strings.ToLower(action) {
	case "register":
		return svc.RegisterAsset(domain.VideoAsset{ID: assetID, SourcePath: "training.mp4", Title: "training session", DurationMS: 18000, FrameRate: 60, Width: 1920, Height: 1080})
	case "draft":
		record, err := svc.DraftPreview(service.PreviewRequest{ID: previewID, AssetID: assetID, SpeedLabel: speed, Crop: domain.CropWindow{StartMS: 2000, EndMS: 7000}, Resolution: domain.ResolutionMedium, Interpolate: true})
		if err != nil {
			return err
		}
		fmt.Println(report.ProgressText(record.Status), record.Spec.Summary())
		return nil
	case "generate":
		record, err := svc.GeneratePreview(previewID)
		if err != nil {
			return err
		}
		fmt.Println(report.StatusBadge(record.Status), record.Message)
		return nil
	case "inspect":
		record, err := svc.InspectPreview(previewID)
		if err != nil {
			return err
		}
		fmt.Printf("%s %s %s\n", record.Spec.ID, report.StatusBadge(record.Status), record.Message)
		return nil
	case "history":
		history, err := service.LoadHistoryFromService(svc, previewID)
		if err != nil {
			return err
		}
		fmt.Println(strings.Join(history.Kinds(), ","))
		return nil
	default:
		fmt.Println("slowpreview: use -action register|draft|generate|inspect|history")
		return nil
	}
}
