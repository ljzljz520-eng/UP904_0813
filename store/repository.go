package store

import (
	"fmt"
	"sort"
	"strconv"

	"go.etcd.io/bbolt"
	"slowpreview/domain"
)

type Repository struct {
	db *Database
}

func NewRepository(db *Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveAsset(asset domain.VideoAsset) error {
	data, err := encode(asset)
	if err != nil {
		return err
	}
	return r.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte(assetsBucket)).Put([]byte(asset.ID), data) })
}

func (r *Repository) GetAsset(id string) (domain.VideoAsset, error) {
	var asset domain.VideoAsset
	err := r.db.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket([]byte(assetsBucket)).Get([]byte(id))
		if value == nil {
			return ErrNotFound
		}
		return decode(value, &asset)
	})
	return asset, err
}

func (r *Repository) ListAssets() ([]domain.VideoAsset, error) {
	assets := make([]domain.VideoAsset, 0)
	err := r.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(assetsBucket)).ForEach(func(_, value []byte) error {
			var asset domain.VideoAsset
			if err := decode(value, &asset); err != nil {
				return err
			}
			assets = append(assets, asset)
			return nil
		})
	})
	sort.Slice(assets, func(i, j int) bool { return assets[i].ID < assets[j].ID })
	return assets, err
}

func (r *Repository) SaveSpec(spec domain.PreviewSpec) error {
	data, err := encode(spec)
	if err != nil {
		return err
	}
	return r.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte(specsBucket)).Put([]byte(spec.ID), data) })
}

func (r *Repository) GetSpec(id string) (domain.PreviewSpec, error) {
	var spec domain.PreviewSpec
	err := r.db.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket([]byte(specsBucket)).Get([]byte(id))
		if value == nil {
			return ErrNotFound
		}
		return decode(value, &spec)
	})
	return spec, err
}

func (r *Repository) SaveTask(task domain.RenderTask) error {
	data, err := encode(task)
	if err != nil {
		return err
	}
	return r.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte(tasksBucket)).Put([]byte(task.ID), data) })
}

func (r *Repository) GetTask(id string) (domain.RenderTask, error) {
	var task domain.RenderTask
	err := r.db.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket([]byte(tasksBucket)).Get([]byte(id))
		if value == nil {
			return ErrNotFound
		}
		return decode(value, &task)
	})
	return task, err
}

func (r *Repository) AppendEvent(event domain.ActivityEvent) (domain.ActivityEvent, error) {
	var result domain.ActivityEvent
	err := r.db.Update(func(tx *bbolt.Tx) error {
		meta := tx.Bucket([]byte(metaBucket))
		sequence := uint64(0)
		if value := meta.Get([]byte(sequenceKey)); value != nil {
			parsed, parseErr := strconv.ParseUint(string(value), 10, 64)
			if parseErr != nil {
				return parseErr
			}
			sequence = parsed
		}
		sequence++
		result = event
		result.Seq = int64(sequence)
		data, encodeErr := encode(result)
		if encodeErr != nil {
			return encodeErr
		}
		if putErr := tx.Bucket([]byte(eventsBucket)).Put([]byte(fmt.Sprintf("%020d", sequence)), data); putErr != nil {
			return putErr
		}
		return meta.Put([]byte(sequenceKey), []byte(strconv.FormatUint(sequence, 10)))
	})
	return result, err
}

func (r *Repository) ListEvents(previewID string) ([]domain.ActivityEvent, error) {
	events := make([]domain.ActivityEvent, 0)
	err := r.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(eventsBucket)).ForEach(func(_, value []byte) error {
			var event domain.ActivityEvent
			if err := decode(value, &event); err != nil {
				return err
			}
			if previewID == "" || event.PreviewID == previewID {
				events = append(events, event)
			}
			return nil
		})
	})
	sort.Slice(events, func(i, j int) bool { return events[i].Seq < events[j].Seq })
	return events, err
}

func (r *Repository) DeletePreview(id string) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		if err := tx.Bucket([]byte(specsBucket)).Delete([]byte(id)); err != nil {
			return err
		}
		return tx.Bucket([]byte(tasksBucket)).Delete([]byte("task-" + id))
	})
}
