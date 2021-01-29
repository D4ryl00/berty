package proximitytransport

import (
	"context"
	"io"
	"sync"

	"go.uber.org/zap"
)

type mplex struct {
	inputCaches []*RingBufferMap
	inputLock   sync.Mutex
	input       chan []byte

	output *io.PipeWriter

	ctx    context.Context
	logger *zap.Logger
}

func newMplex(ctx context.Context, logger *zap.Logger) *mplex {
	logger = logger.Named("mplex")
	return &mplex{
		input:  make(chan []byte),
		ctx:    ctx,
		logger: logger,
	}
}

func (m *mplex) setOutput(o *io.PipeWriter) {
	m.output = o
}

func (m *mplex) addInputCache(c *RingBufferMap) {
	m.inputLock.Lock()
	m.inputCaches = append(m.inputCaches, c)
	m.inputLock.Unlock()
}

func (m *mplex) write(s []byte) {
	m.logger.Debug("write", zap.Binary("payload", s))
	_, err := m.output.Write(s)
	if err != nil {
		m.logger.Error("write: write pipe error", zap.Error(err))
	} else {
		m.logger.Debug("write: successful write pipe")
	}
}

// run flushes caches and read input channel
func (m *mplex) run(peerID string) {
	m.logger.Debug("run: started")
	// flush caches
	m.inputLock.Lock()
	for _, cache := range m.inputCaches {
		m.logger.Debug("run: flushing one cache")

		payloads := cache.Flush(peerID)
		for payload := range payloads {
			m.write(payload)
		}
	}
	m.inputLock.Unlock()

	// read input
	m.logger.Debug("run: reading input channel")
	for {
		select {
		case payload := <-m.input:
			m.write(payload)
		case <-m.ctx.Done():
			return
		}
	}
}
