package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	ds "github.com/ipfs/go-datastore"
	ipfs_ds "github.com/ipfs/go-datastore"
	ds_sync "github.com/ipfs/go-datastore/sync"
	ipfs_cfg "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	ipfs_p2p "github.com/ipfs/go-ipfs/core/node/libp2p"
	ipfs_repo "github.com/ipfs/go-ipfs/repo"

	ci "github.com/libp2p/go-libp2p-core/crypto"

	// "github.com/ipfs/go-ipfs/core"
	"github.com/libp2p/go-libp2p"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"

	"github.com/libp2p/go-libp2p-core/event"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

var tcpBertyRelays = []string{
	"/ip4/51.159.21.214/tcp/4040/p2p/QmdT7AmhhnbuwvCpa5PH1ySK9HJVB82jr3fo1bxMxBPW6p",
	// "/ip4/51.15.25.224/tcp/4040/p2p/12D3KooWHhDBv6DJJ4XDWjzEXq6sVNEs6VuxsV1WyBBEhPENHzcZ",
}

var quicBertyRelays = []string{
	"/ip4/51.159.21.214/udp/4040/quic/p2p/QmdT7AmhhnbuwvCpa5PH1ySK9HJVB82jr3fo1bxMxBPW6p",
	"/ip4/51.15.25.224/udp/4040/quic/p2p/12D3KooWHhDBv6DJJ4XDWjzEXq6sVNEs6VuxsV1WyBBEhPENHzcZ",
}

var tcpIPFSRelays = []string{
	"/ip4/147.75.80.110/tcp/4001/p2p/QmbFgm5zan8P6eWWmeyfncR5feYEMPbht5b1FW1C37aQ7y",
	"/ip4/147.75.195.153/tcp/4001/p2p/QmW9m57aiBDHAkKj9nmFSEn7ZqrcF1fZS4bipsTCHburei",
	"/ip4/147.75.70.221/tcp/4001/p2p/Qme8g49gm3q4Acp7xWBKg3nAa9fxZ1YmyDJdyGgoG6LsXh",
}

var quicIPFSRelays = []string{
	"/ip4/147.75.80.110/udp/4001/quic/p2p/QmbFgm5zan8P6eWWmeyfncR5feYEMPbht5b1FW1C37aQ7y",
	"/ip4/147.75.195.153/udp/4001/quic/p2p/QmW9m57aiBDHAkKj9nmFSEn7ZqrcF1fZS4bipsTCHburei",
	"/ip4/147.75.70.221/udp/4001/quic/p2p/Qme8g49gm3q4Acp7xWBKg3nAa9fxZ1YmyDJdyGgoG6LsXh",
}

const (
	staticBertyRelayMode = "static-berty"
	staticIPFSRelayMode  = "static-ipfs"
	discoveryRelayMode   = "discovery"
	disabledRelayMode    = "none"
)

type serverOpts struct {
	port  int
	ip6   bool
	relay string
}

func createServerHost(ctx context.Context, gOpts *globalOpts, sOpts *serverOpts) (host.Host, error) {
	if sOpts.relay == discoveryRelayMode || sOpts.relay == disabledRelayMode {
		gOpts.dht = true
	}

	opts, err := globalOptsToLibp2pOpts(ctx, gOpts) // Get identity and transport
	if err != nil {
		return nil, err
	}

	if sOpts.relay == disabledRelayMode { // If no relay, add relevant listener
		if sOpts.ip6 {
			if gOpts.tcp {
				opts = append(opts, libp2p.ListenAddrStrings(fmt.Sprintf("/ip6/::/tcp/%d", sOpts.port)))
			} else {
				opts = append(opts, libp2p.ListenAddrStrings(fmt.Sprintf("/ip6/::/udp/%d/quic", sOpts.port)))
			}
		} else {
			if gOpts.tcp {
				opts = append(opts, libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", sOpts.port)))
			} else {
				opts = append(opts, libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic", sOpts.port)))
			}
		}
		opts = append(opts, libp2p.NATPortMap()) // Open port on NAT for access through public IP
		opts = append(opts, libp2p.DisableRelay())
	} else {
		fmt.Println("forcing manual listening on 4001")
		opts = append(opts, libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/4001"))
		// opts = append(opts, libp2p.ListenAddrs()) // If using relay, set no listener
		opts = append(opts, libp2p.EnableAutoRelay())
		opts = append(opts, libp2p.ForceReachabilityPrivate())
	}

	if sOpts.relay != discoveryRelayMode && sOpts.relay != disabledRelayMode {
		var (
			maddrRelays  []string
			staticRelays []peer.AddrInfo
		)

		if sOpts.relay == staticBertyRelayMode {
			if gOpts.tcp {
				maddrRelays = tcpBertyRelays
			} else {
				maddrRelays = quicBertyRelays
			}
		} else {
			if gOpts.tcp {
				maddrRelays = tcpIPFSRelays
			} else {
				maddrRelays = quicIPFSRelays
			}
		}

		for _, addr := range maddrRelays {
			maddr, err := ma.NewMultiaddr(addr)
			if err != nil {
				log.Printf("error: can't parse Multiaddr: %v\n", err)
				continue
			}

			pi, err := peer.AddrInfoFromP2pAddr(maddr)
			if err != nil {
				log.Printf("error: can't parse AddrInfo: %v\n", err)
				continue
			}

			staticRelays = append(staticRelays, *pi)
		}

		opts = append(opts, libp2p.StaticRelays(staticRelays))
	}

	h, err := libp2p.New(opts...) // Create host
	if err != nil {
		return nil, err
	}

	monitorConnsCount(h, gOpts.limit)

	return h, nil
}

func printHint(h host.Host, gOpts *globalOpts, sOpts *serverOpts) {
	var serverAddr ma.Multiaddr

	if sOpts.relay == disabledRelayMode {
		log.Print("Waiting for public address...")
	} else {
		log.Print("Waiting for relay address... pid=", h.ID().Pretty())
	}

	eventReceiver, err := h.EventBus().Subscribe(new(event.EvtLocalAddressesUpdated))
	if err != nil {
		log.Fatalf("can't subscribe to local addresses updated events: %v", err)
	}
	defer eventReceiver.Close()

	for ev := range eventReceiver.Out() {
		serverAddr = nil
		update := ev.(event.EvtLocalAddressesUpdated)

		for _, addr := range update.Current {
			if addr.Action != event.Added {
				continue
			}

			if sOpts.relay != disabledRelayMode {
				if _, err := addr.Address.ValueForProtocol(ma.P_CIRCUIT); err != nil {
					continue
				}

				if gOpts.tcp {
					if _, err := addr.Address.ValueForProtocol(ma.P_TCP); err != nil {
						continue
					}
					serverAddr = addr.Address
				} else {
					if _, err := addr.Address.ValueForProtocol(ma.P_QUIC); err != nil {
						continue
					}
					serverAddr = addr.Address
				}
			} else if sOpts.ip6 {
				if _, err := addr.Address.ValueForProtocol(ma.P_IP6); err != nil {
					continue
				}
				if manet.IsPublicAddr(addr.Address) {
					serverAddr = addr.Address
				}
			} else {
				if _, err := addr.Address.ValueForProtocol(ma.P_IP4); err != nil {
					continue
				}
				if manet.IsPublicAddr(addr.Address) {
					serverAddr = addr.Address
				}
			}
		}

		if serverAddr != nil {
			hostAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", h.ID().Pretty()))
			if err != nil {
				panic(err)
			}
			fullAddr := serverAddr.Encapsulate(hostAddr)

			hint := "Now run: './bench"
			if gOpts.insecure {
				hint += " -insecure"
			}
			if gOpts.tcp {
				hint += " -tcp"
			}
			hint += fmt.Sprintf(" client -dest %s [-request ...] [-size X] [-reco]'", fullAddr)
			log.Println(hint)
		}
	}
}

func genIpfsHostWithDHT(ctx context.Context) (host.Host, error) {
	dsync := ds_sync.MutexWrap(ds.NewMapDatastore())
	repo, err := createDefaultMockedRepo(dsync)
	if err != nil {
		return nil, err
	}

	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: ipfs_p2p.DHTServerOption,
		// Routing: ipfs_p2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		Repo: repo,
	}

	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, err
	}

	return node.PeerHost, nil
}

func createDefaultMockedRepo(dstore ipfs_ds.Batching) (ipfs_repo.Repo, error) {
	c := ipfs_cfg.Config{}
	priv, pub, err := ci.GenerateKeyPairWithReader(ci.RSA, 2048, rand.Reader)
	if err != nil {
		return nil, err
	}

	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		return nil, err
	}

	privkeyb, err := ci.MarshalPrivateKey(priv)
	if err != nil {
		return nil, err
	}

	c.Bootstrap = ipfs_cfg.DefaultBootstrapAddresses
	// c.Bootstrap = []string{} // @NOTE(gfanton): if we remove bootstrap peer the test works

	c.AutoNAT.ServiceMode = ipfs_cfg.AutoNATServiceDisabled
	c.Addresses.Swarm = []string{"/ip4/0.0.0.0/tcp/4001", "/ip4/0.0.0.0/udp/4001/quic"}
	c.Identity.PeerID = pid.Pretty()
	c.Identity.PrivKey = base64.StdEncoding.EncodeToString(privkeyb)

	return &ipfs_repo.Mock{
		D: dstore,
		C: c,
	}, nil
}

func runServer(ctx context.Context, gOpts *globalOpts, sOpts *serverOpts) error {
	h, err := genIpfsHostWithDHT(ctx)
	// h, err := createServerHost(ctx, gOpts, sOpts)
	if err != nil {
		return fmt.Errorf("server host creation failed: %v", err)
	}

	pirelay, err := peer.AddrInfoFromP2pAddr(ma.StringCast(tcpBertyRelays[0]))
	if err != nil {
		return fmt.Errorf("relay address is incorrect: %v", err)
	}
	h.Connect(ctx, *pirelay)

	go printHint(h, gOpts, sOpts)

	h.SetStreamHandler(benchDownloadPID, func(s network.Stream) {
		defer s.Close()

		remotePeerID := s.Conn().RemotePeer().Pretty()
		log.Printf("New download stream from: %s\n", remotePeerID)

		_, err := s.Write([]byte("\n"))
		if err != nil {
			log.Printf("Write error during stream opened ack to %s: %v\n", remotePeerID, err)
		}

		buf := bufio.NewReader(s)
		str, err := buf.ReadString('\n')
		if err != nil {
			log.Printf("Read error during download from %s: %v\n", remotePeerID, err)
			return
		}

		size, err := strconv.Atoi(string(str[:len(str)-1]))
		if err != nil {
			log.Printf("Invalid size received: %s: %v", str, err)
			return
		}

		data := make([]byte, size)
		rand.Read(data)

		_, err = s.Write(data)
		if err != nil {
			log.Printf("Write error during download to %s: %v\n", remotePeerID, err)
		}
		log.Printf("Sent %d bytes to %s", len(data), remotePeerID)
	})

	h.SetStreamHandler(benchUploadPID, func(s network.Stream) {
		defer s.Close()

		remotePeerID := s.Conn().RemotePeer().Pretty()
		log.Printf("New upload stream from: %s\n", remotePeerID)

		_, err := s.Write([]byte("\n"))
		if err != nil {
			log.Printf("Write error during stream opened ack to %s: %v\n", remotePeerID, err)
		}

		reader := bufio.NewReader(s)
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Printf("Read error during upload from %s: %v\n", remotePeerID, err)
			return
		}
		log.Printf("Received %d bytes from %s", len(data), remotePeerID)

		_, err = s.Write([]byte("\n"))
		if err != nil {
			log.Printf("Write error during uploaded ack to %s: %v\n", remotePeerID, err)
		}
	})

	select {}
}
