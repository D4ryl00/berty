package session

import (
	"math"
	"time"

	"berty.tech/berty/v2/go/pkg/tyber"
)

type baseLog struct {
	Level   string  `json:"level"`
	Epoch   float64 `json:"ts"`
	Logger  string  `json:"logger"`
	Caller  string  `json:"caller"`
	Message string  `json:"msg"`
	Time    time.Time
}

func (bl *baseLog) EpochToTime() {
	sec, dec := math.Modf(bl.Epoch)
	bl.Time = time.Unix(int64(sec), int64(dec*(1e9)))
}

type TypedLog struct {
	baseLog
	LogType tyber.LogType `json:"tyberLogType"`
}

type TraceLog struct {
	TypedLog
	Trace tyber.Trace `json:"trace"`
}

func (tl *TraceLog) ToAppTrace() *AppTrace {
	return &AppTrace{
		ID:          tl.Trace.TraceID,
		Name:        tl.Message,
		InitialName: tl.Message,
		Steps:       []*AppStep{},
		Status: Status{
			StatusType: tyber.Running,
			Started:    tl.Time,
		},
	}
}

type StepLog struct {
	TypedLog
	Step tyber.Step `json:"step"`
}

func (sl *StepLog) ToAppStep() *AppStep {
	return &AppStep{
		Name:    sl.Message,
		Details: sl.Step.Details,
		Status: Status{
			StatusType: sl.Step.Status,
			Started:    sl.Time,
		},
		ForceReopen:     sl.Step.ForceReopen,
		UpdateTraceName: sl.Step.UpdateTraceName,
	}
}

type SubscribeLog struct {
	TypedLog
	Subscribe tyber.Subscribe `json:"subscribe"`
}

func (subl *SubscribeLog) ToAppSubscribe() *AppSubscribe {
	sub := &AppSubscribe{
		TargetName:    subl.Subscribe.TargetStepName,
		TargetDetails: subl.Subscribe.TargetDetails,
		SubscribeStep: &AppStep{
			Name:    subl.Message,
			Details: append(subl.Subscribe.TargetDetails, tyber.Detail{Name: "TargetName", Description: subl.Subscribe.TargetStepName}),
			Status: Status{
				StatusType: tyber.Running,
				Started:    subl.Time,
			},
		},
	}
	return sub
}

type EventLog struct {
	TypedLog
	Event tyber.Event `json:"event"`
}

func (el *EventLog) ToAppStep(s tyber.Step) *AppStep {
	return &AppStep{
		Name:    el.Message,
		Details: append(el.Event.Details, s.Details...),
		Status: Status{
			StatusType: s.Status,
			Started:    el.Time,
		},
		ForceReopen:     s.ForceReopen,
		UpdateTraceName: s.UpdateTraceName,
	}
}
