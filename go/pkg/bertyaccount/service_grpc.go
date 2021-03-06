package bertyaccount

import (
	"context"
	"strings"

	ma "github.com/multiformats/go-multiaddr"

	"berty.tech/berty/v2/go/pkg/errcode"
	"berty.tech/berty/v2/go/pkg/messengertypes"
)

// Get GRPC listener addresses
func (s *service) GetGRPCListenerAddrs(ctx context.Context, req *GetGRPCListenerAddrs_Request) (*GetGRPCListenerAddrs_Reply, error) {
	m, err := s.getInitManager()
	if err != nil {
		return nil, err
	}

	grpcListeners := m.GetGRPCListeners()
	entries := make([]*GetGRPCListenerAddrs_Reply_Entry, len(grpcListeners))
	for i, l := range grpcListeners {
		ps := make([]string, 0)
		ma.ForEach(l.GRPCMultiaddr(), func(c ma.Component) bool {
			ps = append(ps, c.Protocol().Name)
			return true
		})

		proto := strings.Join(ps, "/")
		entries[i] = &GetGRPCListenerAddrs_Reply_Entry{
			Maddr: l.Addr().String(),
			Proto: proto,
		}
	}

	return &GetGRPCListenerAddrs_Reply{
		Entries: entries,
	}, nil
}

// GetMessengerClient returns the Messenger Client of the actual Berty account if there is one selected.
func (s *service) GetMessengerClient() (messengertypes.MessengerServiceClient, error) {
	m, err := s.getInitManager()
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	messenger, err := m.GetMessengerClient()
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	return messenger, err
}
