package parser

import (
	"bufio"

	"berty.tech/berty/v2/go/pkg/tyber"
	"github.com/pkg/errors"

	"berty.tech/berty/tool/tyber/go/session"
)

func (p *Parser) parseLogs(s *session.Session, srcScanner *bufio.Scanner) error {
	latestTime := s.Started
	for srcScanner.Scan() {
		log := srcScanner.Text()
		tl := parseTypedLog(log)

		if tl == nil {
			// Not a Tyber log
			continue
		}

		latestTime = tl.Time

		switch tl.LogType {
		case tyber.TraceType:
			trl, err := p.parseTraceLog(log)
			if err != nil {
				p.logger.Error(err.Error())
				continue
			}

			trace := trl.ToAppTrace()

			s.TracesLock.Lock()
			s.Traces = append(s.Traces, trace)
			s.RunningTraces[trace.ID] = trace
			s.TracesLock.Unlock()

			if s.Openned {
				p.EventChan <- trace.ToCreateTraceEvent()
			}

		case tyber.StepType:
			sl, err := parseStepLog(log)
			if err != nil {
				p.logger.Error(err.Error())
				continue
			}

			step := sl.ToAppStep()

			s.TracesLock.Lock()
			parentTrace, ok := s.RunningTraces[sl.Step.ParentTraceID]
			if !ok && step.ForceReopen {
				for _, t := range s.Traces {
					if t.ID == sl.Step.ParentTraceID {
						parentTrace = t
						s.RunningTraces[t.ID] = t
						ok = true
						break
					}
				}
			}

			if !ok {
				p.logger.Errorf("parent trace not found in running traces: %s", log)
				s.TracesLock.Unlock()
				continue
			}

			shouldUpdateTraceName := len(step.UpdateTraceName) > 0
			if shouldUpdateTraceName {
				parentTrace.Name = step.UpdateTraceName
			}
			parentTrace.Steps = append(parentTrace.Steps, step)
			// TODO
			// if step.Status == tyber.Running {
			// 	s.runningSteps[]
			// } elseif sl.Step.EndTrace {
			if sl.Step.EndTrace {
				parentTrace.Finished = step.Started
				parentTrace.StatusType = step.StatusType
				delete(s.RunningTraces, parentTrace.ID)
			}
			s.TracesLock.Unlock()

			if s.Openned {
				cse := step.ToCreateStepEvent(parentTrace.ID)
				p.EventChan <- cse
				if sl.Step.EndTrace || shouldUpdateTraceName {
					p.EventChan <- parentTrace.ToUpdateTraceEvent()
				}
			}
		case tyber.SubscribeType:
			subl, err := parseSubscribeLog(log)
			if err != nil {
				p.logger.Error(err.Error())
				continue
			}

			subscribe := subl.ToAppSubscribe()

			s.TracesLock.Lock()
			parentTrace, ok := s.RunningTraces[subl.Subscribe.StepToAdd.ParentTraceID]

			if !ok {
				p.logger.Errorf("parent trace not found in running traces: %s", log)
				continue
			}

			parentTrace.Subs = append(parentTrace.Subs, session.SubTarget{TargetName: subscribe.TargetName, TargetDetails: subscribe.TargetDetails, StepToAdd: subl.Subscribe.StepToAdd})
			s.TracesLock.Unlock()

		case tyber.EventType:
			el, err := parseEventLog(log)
			if err != nil {
				p.logger.Error(err.Error())
				continue
			}

			s.TracesLock.Lock()
			for _, parentTrace := range s.RunningTraces {
				changed := false
				for _, sub := range parentTrace.Subs {
					if sub.TargetName == el.Message {
						match := true
						for _, tdet := range sub.TargetDetails {
							found := false
							for _, det := range el.Event.Details {
								if det.Name == tdet.Name && det.Description == tdet.Description {
									found = true
									break
								}
							}
							if !found {
								match = false
								break
							}
						}
						if match {
							step := el.ToAppStep(sub.StepToAdd)
							parentTrace.Steps = append(parentTrace.Steps, step)
							if sub.StepToAdd.EndTrace {
								parentTrace.Finished = step.Started
								parentTrace.StatusType = step.StatusType
								delete(s.RunningTraces, parentTrace.ID)
								changed = true
							}
						}
					}
				}
				if changed {
					p.EventChan <- parentTrace.ToUpdateTraceEvent()
				}
			}
			s.TracesLock.Unlock()
		}
	}

	s.Finished = latestTime

	if err := srcScanner.Err(); err != nil {
		s.StatusType = tyber.Failed
		p.EventChan <- sessionToUpdateEvent(s)
		// TODO: ADD ERROR TO ALL RUNNING TRACES
		return errors.Wrap(err, "parsing traces failed")
	}

	s.StatusType = tyber.Succeeded
	p.EventChan <- sessionToUpdateEvent(s)

	return nil
}
