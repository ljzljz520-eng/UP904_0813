package catalog

import (
	"sort"
	"strings"

	"slowpreview/domain"
)

type AssetIndex struct {
	assets map[string]domain.VideoAsset
	order  []string
}

func NewAssetIndex(assets []domain.VideoAsset) *AssetIndex {
	index := &AssetIndex{assets: make(map[string]domain.VideoAsset), order: make([]string, 0, len(assets))}
	for _, asset := range assets {
		index.Upsert(asset)
	}
	return index
}

func (i *AssetIndex) Upsert(asset domain.VideoAsset) {
	if i.assets == nil {
		i.assets = make(map[string]domain.VideoAsset)
	}
	if _, exists := i.assets[asset.ID]; !exists {
		i.order = append(i.order, asset.ID)
	}
	i.assets[asset.ID] = asset
}

func (i *AssetIndex) Remove(id string) bool {
	if _, exists := i.assets[id]; !exists {
		return false
	}
	delete(i.assets, id)
	filtered := i.order[:0]
	for _, candidate := range i.order {
		if candidate != id {
			filtered = append(filtered, candidate)
		}
	}
	i.order = filtered
	return true
}

func (i *AssetIndex) Get(id string) (domain.VideoAsset, bool) {
	asset, ok := i.assets[id]
	return asset, ok
}

func (i *AssetIndex) IDs() []string {
	ids := append([]string(nil), i.order...)
	return ids
}

func (i *AssetIndex) Snapshot() []domain.VideoAsset {
	assets := make([]domain.VideoAsset, 0, len(i.order))
	for _, id := range i.order {
		if asset, ok := i.assets[id]; ok {
			assets = append(assets, asset)
		}
	}
	return assets
}

func (i *AssetIndex) Find(query string) []domain.VideoAsset {
	query = strings.ToLower(strings.TrimSpace(query))
	result := make([]domain.VideoAsset, 0)
	for _, asset := range i.Snapshot() {
		if query == "" || strings.Contains(strings.ToLower(asset.Title), query) || strings.Contains(strings.ToLower(asset.ID), query) {
			result = append(result, asset)
		}
	}
	sort.SliceStable(result, func(a, b int) bool {
		if result[a].Title == result[b].Title {
			return result[a].ID < result[b].ID
		}
		return result[a].Title < result[b].Title
	})
	return result
}

func (i *AssetIndex) FilterByResolution(resolution domain.Resolution) []domain.VideoAsset {
	result := make([]domain.VideoAsset, 0)
	for _, asset := range i.Snapshot() {
		if resolution == domain.ResolutionHigh && asset.Width >= 2560 && asset.Height >= 1440 {
			result = append(result, asset)
		} else if resolution == domain.ResolutionMedium && asset.Width >= 1920 && asset.Height >= 1080 {
			result = append(result, asset)
		} else if resolution == domain.ResolutionLow && asset.Width < 1920 {
			result = append(result, asset)
		}
	}
	return result
}
