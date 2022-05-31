package session

import (
	"io"
	"sync"

	"berty.tech/berty/tool/tyber/go/logger"
)

type Session struct {
	ID          string      `json:"id"`
	DisplayName string      `json:"displayName"`
	SrcName     string      `json:"srcName"`
	SrcType     SrcType     `json:"srcType"`
	Header      *Header     `json:"header"`
	Traces      []*AppTrace `json:"traces"`
	Status
	// Internals
	srcCloser     io.Closer
	runningTraces map[string]*AppTrace
	tracesLock    sync.Mutex
	openned       bool
	logger        *logger.Logger
	// runningSteps map[string]*Step TODO
}

type SrcType string

const (
	FileType    SrcType = "file"
	NetworkType SrcType = "network"
)

func NewSession(l *logger.Logger, srcName string, srcType SrcType, srcCloser io.Closer) *Session {
	return &Session{
		logger:        l,
		SrcName:       srcName,
		SrcType:       srcType,
		Traces:        []*AppTrace{},
		srcCloser:     srcCloser,
		runningTraces: map[string]*AppTrace{},
		// runningSteps: map[string]*Step{}, TODO
	}
}

func (s *Session) SrcCloser() io.Closer {
	return s.srcCloser
}

func (s *Session) SetOpenned(openned bool) {
	s.tracesLock.Lock()
	s.openned = openned
	s.tracesLock.Unlock()
}

func (s *Session) IsOpenned() bool {
	s.tracesLock.Lock()
	defer s.tracesLock.Unlock()

	return s.openned
}

func (s *Session) GetRunningTrace(id string) (*AppTrace, bool) {
	s.tracesLock.Lock()
	defer s.tracesLock.Unlock()

	v, ok := s.runningTraces[id]
	return v, ok
}

func (s *Session) SetRunningTrace(id string, v *AppTrace) {
	s.tracesLock.Lock()
	s.runningTraces[id] = v
	s.tracesLock.Unlock()
}

func (s *Session) DeleteRunningTrace(id string) {
	s.tracesLock.Lock()
	delete(s.runningTraces, id)
	s.tracesLock.Unlock()
}

func (s *Session) ForEachRunningTrace(f func(*AppTrace)) {
	s.tracesLock.Lock()
	for _, trace := range s.runningTraces {
		f(trace)
	}
	s.tracesLock.Unlock()
}

func (s *Session) AddTrace(trace *AppTrace) {
	s.tracesLock.Lock()
	s.Traces = append(s.Traces, trace)
	s.runningTraces[trace.ID] = trace
	s.tracesLock.Unlock()
}

func (s *Session) GenerateCreateTraceEvents() []CreateTraceEvent {
	var events []CreateTraceEvent

	s.tracesLock.Lock()
	for _, t := range s.Traces {
		events = append(events, t.ToCreateTraceEvent())
	}
	s.tracesLock.Unlock()

	return events
}
