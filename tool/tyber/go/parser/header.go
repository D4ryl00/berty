package parser

import (
	"bufio"
	"encoding/json"

	"github.com/pkg/errors"

	"berty.tech/berty/tool/tyber/go/session"
)

func (p *Parser) parseHeader(s *session.Session, srcScanner *bufio.Scanner) error {
	h := &session.Header{}
	for srcScanner.Scan() {
		log := srcScanner.Text()
		if err := json.Unmarshal([]byte(log), h); err != nil {
			return err
		}

		if h.Manager.Session.ID != "" {
			break
		}
	}

	if h.Manager.Session.ID == "" {
		return errors.New("invalid log: header not found")
	}

	// TODO: add other checks

	h.EpochToTime()
	s.Header = h
	s.ID = h.Manager.Session.ID
	s.DisplayName = h.Manager.Node.Messenger.DisplayName
	s.Started = h.Time

	return nil
}
