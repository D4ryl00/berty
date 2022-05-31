package parser

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"berty.tech/berty/tool/tyber/go/logger"
	"berty.tech/berty/tool/tyber/go/session"
	"berty.tech/berty/v2/go/pkg/tyber"
)

type Parser struct {
	ctx            context.Context
	logger         *logger.Logger
	sessionManager *session.Manager
	listener       *net.TCPListener
	cancelNetwork  context.CancelFunc
	EventChan      chan interface{} // TODO: base event definition
}

func New(ctx context.Context, l *logger.Logger, sessionManager *session.Manager) *Parser {
	return &Parser{
		ctx:            ctx,
		logger:         l.Named("parser"),
		sessionManager: sessionManager,
		EventChan:      make(chan interface{}),
	}
}

func (p *Parser) ParseFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		p.logger.Errorf("opening file failed: %v", err)
		return err
	}

	p.startParsingSession(path, session.FileType, file)

	return nil
}

func (p *Parser) NetworkListen(address, port string) error {
	if p.listener != nil {
		p.cancelNetwork()
		p.listener = nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	p.cancelNetwork = cancel

	localAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		return err
	}

	p.listener, err = net.ListenTCP("tcp", localAddr)
	if err != nil {
		return err
	}

	go func(l *net.TCPListener) {
		p.logger.Infof("started listening on %s", l.Addr().String())
		defer l.Close()

		for {
			select {
			case <-ctx.Done():
				p.logger.Infof("stopped listening on %s", l.Addr().String())
				return

			default:
				if err := l.SetDeadline(time.Now().Add(time.Second)); err != nil {
					p.logger.Errorf("can't set deadline on tcp listener: %v", err)
					return
				}

				conn, err := l.Accept()
				if err != nil {
					if os.IsTimeout(err) {
						continue
					}
					p.logger.Errorf("TCP accept error: %v", err)
					return
				}

				p.startParsingSession(conn.RemoteAddr().String(), session.NetworkType, conn)
			}
		}
	}(p.listener)

	return nil
}

func (p *Parser) startParsingSession(srcName string, srcType session.SrcType, srcIO io.ReadCloser) {
	s := session.NewSession(p.logger, srcName, srcType, srcIO)
	srcScanner := bufio.NewScanner(srcIO)

	if err := p.parseHeader(s, srcScanner); err != nil {
		p.logger.Errorf("starting session failed with logs from %s (%s): %v", srcType, srcName, err)
		return
	}

	p.logger.Infof("started session %s with logs from %s (%s)", s.ID, srcType, srcName)

	s.StatusType = tyber.Running

	p.sessionManager.AddSession(s)

	// wait := make(chan interface{})

	go func() {
		if err := p.parseLogs(s, srcScanner); err != nil {
			p.logger.Errorf("parsing session %s logs from %s (%s) failed: %v", s.ID, srcType, srcName, err)
			s.StatusType = tyber.Failed
		} else {
			p.logger.Infof("successfully parsed session %s logs from %s (%s)", s.ID, srcType, srcName)
			s.StatusType = tyber.Succeeded
		}

		// if srcType == session.FileType {
		// 	wait <- struct{}{}
		// }
	}()

	// if srcType == session.FileType {
	// 	<-wait
	// }

	p.EventChan <- session.SessionToCreateEvent(s)
}
