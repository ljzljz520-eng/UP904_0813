package store

import (
	"sort"

	"go.etcd.io/bbolt"
	"slowpreview/domain"
)

func (r *Repository) ListSpecs() ([]domain.PreviewSpec, error) {
	specs := make([]domain.PreviewSpec, 0)
	err := r.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(specsBucket)).ForEach(func(_, value []byte) error {
			var spec domain.PreviewSpec
			if err := decode(value, &spec); err != nil {
				return err
			}
			specs = append(specs, spec)
			return nil
		})
	})
	sort.SliceStable(specs, func(i, j int) bool { return specs[i].ID < specs[j].ID })
	return specs, err
}

func (r *Repository) ListTasks() ([]domain.RenderTask, error) {
	tasks := make([]domain.RenderTask, 0)
	err := r.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(tasksBucket)).ForEach(func(_, value []byte) error {
			var task domain.RenderTask
			if err := decode(value, &task); err != nil {
				return err
			}
			tasks = append(tasks, task)
			return nil
		})
	})
	sort.SliceStable(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })
	return tasks, err
}

func (r *Repository) Snapshot() (domain.LibrarySnapshot, error) {
	assets, err := r.ListAssets()
	if err != nil {
		return domain.LibrarySnapshot{}, err
	}
	specs, err := r.ListSpecs()
	if err != nil {
		return domain.LibrarySnapshot{}, err
	}
	tasks, err := r.ListTasks()
	if err != nil {
		return domain.LibrarySnapshot{}, err
	}
	events, err := r.ListEvents("")
	if err != nil {
		return domain.LibrarySnapshot{}, err
	}
	previews := make([]domain.PreviewRecord, 0, len(specs))
	for _, spec := range specs {
		task, taskErr := r.GetTask("task-" + spec.ID)
		if taskErr != nil {
			continue
		}
		previews = append(previews, domain.PreviewRecord{Spec: spec, Task: task, Status: task.Status})
	}
	if len(previews) == 0 && len(tasks) > 0 {
		for _, task := range tasks {
			previews = append(previews, domain.PreviewRecord{Task: task, Status: task.Status})
		}
	}
	return domain.LibrarySnapshot{Assets: assets, Previews: previews, Events: events}, nil
}

func (r *Repository) Count(entity string) (int, error) {
	bucket := bucketForEntity(entity)
	if bucket == "" {
		return 0, ErrNotFound
	}
	count := 0
	err := r.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(bucket)).ForEach(func(_, _ []byte) error {
			count++
			return nil
		})
	})
	return count, err
}
