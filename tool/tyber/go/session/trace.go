package session

import "berty.tech/berty/v2/go/pkg/tyber"

type AppTrace struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	InitialName string     `json:"initialName"`
	Steps       []*AppStep `json:"steps"`
	Subs        []SubTarget
	Status
}

type SubTarget struct {
	TargetName    string
	TargetDetails []tyber.Detail
	StepToAdd     tyber.Step
}

func (t *AppTrace) ToCreateTraceEvent() CreateTraceEvent {
	return CreateTraceEvent{*t}
}

func (t *AppTrace) ToUpdateTraceEvent() UpdateTraceEvent {
	return UpdateTraceEvent{
		ID:     t.ID,
		Status: t.Status,
		Name:   t.Name,
	}
}

type AppStep struct {
	Name    string         `json:"name"`
	Details []tyber.Detail `json:"details"`
	Status
	ForceReopen     bool   `json:"forceReopen"`
	UpdateTraceName string `json:"updateTraceName"`
}

func (s *AppStep) ToCreateStepEvent(id string) CreateStepEvent {
	return CreateStepEvent{
		ID:      id,
		AppStep: *s,
	}
}

type AppSubscribe struct {
	TargetDetails []tyber.Detail `json:"details"`
	TargetName    string         `json:"targetName"`
	ParentTraceID string
	SubscribeStep *AppStep
}

func (sub *AppSubscribe) ToInitialCreateStepEvent(id string) CreateStepEvent {
	return CreateStepEvent{
		ID:      id,
		AppStep: *sub.SubscribeStep,
	}
}
