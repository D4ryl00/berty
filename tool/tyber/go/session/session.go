package session

import (
	"io"
	"sync"

	"berty.tech/berty/tool/tyber/go/logger"
)

type Session struct {
	ID            string      `json:"id"`
	DisplayName   string      `json:"displayName"`
	SrcName       string      `json:"srcName"`
	SrcType       SrcType     `json:"srcType"`
	Header        *Header     `json:"header"`
	Traces        []*AppTrace `json:"traces"`
	SrcCloser     io.Closer
	TracesLock    sync.Mutex
	RunningTraces map[string]*AppTrace
	Openned       bool
	Status
	// Internals
	logger *logger.Logger
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
		SrcCloser:     srcCloser,
		Traces:        []*AppTrace{},
		RunningTraces: map[string]*AppTrace{},
		// runningSteps: map[string]*Step{}, TODO
	}
}
