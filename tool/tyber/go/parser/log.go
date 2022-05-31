package parser

import (
	"encoding/json"
	"fmt"

	"berty.tech/berty/tool/tyber/go/session"
	"berty.tech/berty/v2/go/pkg/tyber"
)

func parseTypedLog(log string) *session.TypedLog {
	tl := &session.TypedLog{}
	if err := json.Unmarshal([]byte(log), tl); err != nil {
		return nil
	}
	tl.EpochToTime()

	if !tl.LogType.IsKnown() {
		return nil
	}

	return tl
}

func (p *Parser) parseTraceLog(log string) (*session.TraceLog, error) {
	tl := &session.TraceLog{}
	if err := json.Unmarshal([]byte(log), tl); err != nil {
		return nil, err
	}
	tl.EpochToTime()

	if tl.Message == "" {
		return nil, fmt.Errorf("invalid trace log: message empty: %s", log)
	}

	return tl, nil
}

func parseStepLog(log string) (*session.StepLog, error) {
	sl := &session.StepLog{}
	if err := json.Unmarshal([]byte(log), sl); err != nil {
		return nil, err
	}
	sl.EpochToTime()

	if sl.Message == "" {
		return nil, fmt.Errorf("invalid step log: message empty: %s", log)
	}
	if sl.Step.Status != tyber.Running && sl.Step.Status != tyber.Succeeded && sl.Step.Status != tyber.Failed {
		return nil, fmt.Errorf("invalid step log: invalid status: %s", log)
	}

	return sl, nil
}

func parseSubscribeLog(log string) (*session.SubscribeLog, error) {
	sl := &session.SubscribeLog{}
	if err := json.Unmarshal([]byte(log), sl); err != nil {
		return nil, err
	}
	sl.EpochToTime()

	if sl.Message == "" {
		return nil, fmt.Errorf("invalid step log: message empty: %s", log)
	}

	return sl, nil
}

func parseEventLog(log string) (*session.EventLog, error) {
	el := &session.EventLog{}
	if err := json.Unmarshal([]byte(log), el); err != nil {
		return nil, err
	}
	el.EpochToTime()

	if el.Message == "" {
		return nil, fmt.Errorf("invalid event log: message empty: %s", log)
	}

	return el, nil
}

// func toAppStep(s tyber.Step, el *session.EventLog) *AppStep {
// 	return &AppStep{
// 		Name:    el.Message,
// 		Details: append(el.Event.Details, s.Details...),
// 		Status: Status{
// 			StatusType: s.Status,
// 			Started:    el.Time,
// 		},
// 		ForceReopen:     s.ForceReopen,
// 		UpdateTraceName: s.UpdateTraceName,
// 	}
// }
