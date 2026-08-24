package service

import (
	"errors"
	"fmt"
	"strings"

	"slowpreview/domain"
	"slowpreview/engine"
	"slowpreview/store"
	"slowpreview/validate"
)

var ErrInvalidWorkflow = errors.New("workflow operation is invalid")

type PreviewRequest struct {
	ID          string
	AssetID     string
	SpeedLabel  string
	Resolution  domain.Resolution
	Crop        domain.CropWindow
	Interpolate bool
	OutputPath  string
}

type Service struct {
	repo    *store.Repository
	builder engine.CommandBuilder
}

func New(repo *store.Repository) *Service {
	return &Service{repo: repo, builder: engine.NewCommandBuilder()}
}

func (s *Service) RegisterAsset(asset domain.VideoAsset) error {
	asset.Title = validate.NormalizeTitle(asset.Title)
	if result := validate.Asset(asset); !result.OK() {
		return fmt.Errorf("asset rejected: %s", result.Error())
	}
	if err := s.repo.SaveAsset(asset); err != nil {
		return err
	}
	_, err := s.repo.AppendEvent(domain.ActivityEvent{PreviewID: asset.ID, Kind: domain.EventAssetRegistered, Message: domain.EventDescription(domain.EventAssetRegistered)})
	return err
}

func (s *Service) DraftPreview(request PreviewRequest) (domain.PreviewRecord, error) {
	asset, err := s.repo.GetAsset(request.AssetID)
	if err != nil {
		return domain.PreviewRecord{}, err
	}
	speed, ok := validate.ParseSpeed(request.SpeedLabel)
	if !ok {
		return domain.PreviewRecord{}, fmt.Errorf("speed rejected: %s", request.SpeedLabel)
	}
	if request.Resolution == "" {
		request.Resolution = validate.SuggestedResolution(asset)
	}
	if request.OutputPath == "" {
		request.OutputPath = validate.DefaultOutputPath(request.ID)
	}
	request.Crop = validate.ClampCrop(asset, request.Crop)
	spec := domain.PreviewSpec{ID: request.ID, AssetID: request.AssetID, Speed: speed, Resolution: request.Resolution, Crop: request.Crop, Interpolate: request.Interpolate, OutputPath: request.OutputPath, RequestedLabel: request.SpeedLabel}
	if result := validate.Preview(asset, spec); !result.OK() {
		return domain.PreviewRecord{}, fmt.Errorf("preview rejected: %s", result.Error())
	}
	record := domain.PreviewRecord{Spec: spec, Status: domain.StatusDraft, Message: "configuration saved", CreatedSeq: 0}
	if err := s.repo.SaveSpec(spec); err != nil {
		return domain.PreviewRecord{}, err
	}
	if err := s.repo.SaveTask(domain.RenderTask{ID: "task-" + spec.ID, PreviewID: spec.ID, Speed: speed, RequestedLabel: request.SpeedLabel, Resolution: spec.Resolution, Crop: spec.Crop, Interpolate: spec.Interpolate, Status: domain.StatusDraft}); err != nil {
		return domain.PreviewRecord{}, err
	}
	event, err := s.repo.AppendEvent(domain.ActivityEvent{PreviewID: spec.ID, Kind: domain.EventPreviewDrafted, Message: spec.Summary()})
	if err != nil {
		return domain.PreviewRecord{}, err
	}
	record.CreatedSeq = event.Seq
	record.UpdatedSeq = event.Seq
	return record, nil
}

func (s *Service) GeneratePreview(id string) (domain.PreviewRecord, error) {
	spec, err := s.repo.GetSpec(id)
	if err != nil {
		return domain.PreviewRecord{}, err
	}
	asset, err := s.repo.GetAsset(spec.AssetID)
	if err != nil {
		return domain.PreviewRecord{}, err
	}
	plan := engine.BuildPlan(asset, spec)
	plan.Task.Command = spec.OutputPath
	command := s.builder.Build(asset, plan)
	plan.Task.Command = command
	receipt := engine.SimulateRender(plan.Task, command)
	status := receipt.Status
	if status == "" {
		status = domain.StatusFailed
	}
	plan.Task.Status = status
	if err := s.repo.SaveTask(plan.Task); err != nil {
		return domain.PreviewRecord{}, err
	}
	kind := domain.EventForStatus(status)
	event, eventErr := s.repo.AppendEvent(domain.ActivityEvent{PreviewID: id, Kind: kind, Message: engine.StatusMessage(receipt)})
	if eventErr != nil {
		return domain.PreviewRecord{}, eventErr
	}
	record := domain.PreviewRecord{Spec: spec, Task: plan.Task, Status: status, Message: engine.StatusMessage(receipt), UpdatedSeq: event.Seq}
	if status == domain.StatusReady {
		record.Message = fmt.Sprintf("%s; filters=%s", record.Message, plan.FilterSummary())
	}
	return record, nil
}

func (s *Service) QueuePreview(id string) error {
	task, err := s.repo.GetTask("task-" + id)
	if err != nil {
		return err
	}
	if !domain.CanTransition(task.Status, domain.StatusQueued) {
		return fmt.Errorf("%w: %s", ErrInvalidWorkflow, domain.TransitionMessage(task.Status, domain.StatusQueued))
	}
	task.Status = domain.StatusQueued
	if err := s.repo.SaveTask(task); err != nil {
		return err
	}
	_, err = s.repo.AppendEvent(domain.ActivityEvent{PreviewID: id, Kind: domain.EventPreviewQueued, Message: "preview queued"})
	return err
}

func (s *Service) ArchivePreview(id string) error {
	task, err := s.repo.GetTask("task-" + id)
	if err != nil {
		return err
	}
	if !domain.CanTransition(task.Status, domain.StatusArchived) {
		return fmt.Errorf("%w: cannot archive %s", ErrInvalidWorkflow, task.Status)
	}
	task.Status = domain.StatusArchived
	if err := s.repo.SaveTask(task); err != nil {
		return err
	}
	if err := s.repo.DeletePreview(id); err != nil {
		return err
	}
	_, err = s.repo.AppendEvent(domain.ActivityEvent{PreviewID: id, Kind: domain.EventPreviewArchived, Message: "preview archived"})
	return err
}

func (s *Service) InspectPreview(id string) (domain.PreviewRecord, error) {
	spec, err := s.repo.GetSpec(id)
	if err != nil {
		return domain.PreviewRecord{}, err
	}
	task, err := s.repo.GetTask("task-" + id)
	if err != nil {
		return domain.PreviewRecord{}, err
	}
	events, err := s.repo.ListEvents(id)
	if err != nil {
		return domain.PreviewRecord{}, err
	}
	status := task.Status
	message := "no activity"
	if len(events) > 0 {
		message = events[len(events)-1].Message
	}
	return domain.PreviewRecord{Spec: spec, Task: task, Status: status, Message: strings.TrimSpace(message), UpdatedSeq: lastSequence(events)}, nil
}

func lastSequence(events []domain.ActivityEvent) int64 {
	if len(events) == 0 {
		return 0
	}
	return events[len(events)-1].Seq
}
