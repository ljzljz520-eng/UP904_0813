package service

import (
	"fmt"

	"slowpreview/analysis"
	"slowpreview/catalog"
	"slowpreview/domain"
	"slowpreview/ingest"
	"slowpreview/playback"
	"slowpreview/store"
)

type LibraryView struct {
	Assets       []domain.VideoAsset
	Previews     []domain.PreviewRecord
	ReviewQueue  catalog.ReviewQueue
	AssetMatches []domain.VideoAsset
}

func (s *Service) ImportManifest(lines []string) (int, []string, error) {
	manifest, parseIssues := ingest.ParseManifest(lines)
	normalized, normalizeIssues := ingest.NormalizeManifest(manifest)
	issues := append(parseIssues, normalizeIssues...)
	imported := 0
	for _, item := range normalized {
		if err := s.RegisterAsset(item.Asset); err != nil {
			issues = append(issues, fmt.Sprintf("%s: %v", item.Asset.ID, err))
			continue
		}
		imported++
	}
	return imported, issues, nil
}

func (s *Service) BuildLibraryView(query string) (LibraryView, error) {
	snapshot, err := s.repo.Snapshot()
	if err != nil {
		return LibraryView{}, err
	}
	index := catalog.NewAssetIndex(snapshot.Assets)
	return LibraryView{Assets: snapshot.Assets, Previews: snapshot.Previews, ReviewQueue: catalog.BuildReviewQueue(snapshot.Previews), AssetMatches: index.Find(query)}, nil
}

func (s *Service) SearchPreviews(query catalog.PreviewQuery) ([]catalog.PreviewResult, error) {
	snapshot, err := s.repo.Snapshot()
	if err != nil {
		return nil, err
	}
	ranked := catalog.RankPreviews(snapshot.Previews, func(record domain.PreviewRecord) int {
		assessment, assessErr := s.AssessPreview(record.Spec)
		if assessErr != nil {
			return 0
		}
		return assessment.Score
	})
	filtered := make([]catalog.PreviewResult, 0, len(ranked))
	for _, result := range ranked {
		if catalog.MatchPreview(result.Record, query, result.Score) {
			filtered = append(filtered, result)
		}
	}
	return filtered, nil
}

func (s *Service) AssessPreview(spec domain.PreviewSpec) (analysis.Assessment, error) {
	asset, err := s.repo.GetAsset(spec.AssetID)
	if err != nil {
		return analysis.Assessment{}, err
	}
	return analysis.Assess(asset, spec), nil
}

func (s *Service) Recommendations(assetID string, crop domain.CropWindow) ([]analysis.Recommendation, error) {
	asset, err := s.repo.GetAsset(assetID)
	if err != nil {
		return nil, err
	}
	return analysis.Recommend(asset, crop), nil
}

func (s *Service) PlaybackSession(previewID string) (playback.Session, error) {
	spec, err := s.repo.GetSpec(previewID)
	if err != nil {
		return playback.Session{}, err
	}
	asset, err := s.repo.GetAsset(spec.AssetID)
	if err != nil {
		return playback.Session{}, err
	}
	return playback.NewSession(asset, spec), nil
}

func RepositoryHealth(repo *store.Repository) (map[string]int, error) {
	health := make(map[string]int)
	for _, entity := range []string{"VideoAsset", "PreviewSpec", "RenderTask", "ActivityEvent", "LibrarySnapshot"} {
		count, err := repo.Count(entity)
		if err != nil {
			return nil, err
		}
		health[entity] = count
	}
	return health, nil
}
