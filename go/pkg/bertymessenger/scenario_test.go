package bertymessenger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"testing"
	"time"

	"berty.tech/berty/v2/go/internal/ipfsutil"
	"berty.tech/berty/v2/go/internal/tinder"
	"berty.tech/berty/v2/go/internal/tracer"
	"berty.tech/berty/v2/go/pkg/bertyprotocol"
	"berty.tech/berty/v2/go/pkg/bertytypes"
	"berty.tech/berty/v2/go/pkg/errcode"
	datastore "github.com/ipfs/go-datastore"
	sync_ds "github.com/ipfs/go-datastore/sync"
	config "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	discovery "github.com/libp2p/go-libp2p-discovery"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	libp2p_mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/api/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	grpc "google.golang.org/grpc"
)

type BertyClient struct {
	Protocol  *bertyprotocol.TestingProtocol
	Messenger MessengerServiceServer
	cancel    func()

	config *bertytypes.InstanceGetConfiguration_Reply
	group  *bertytypes.GroupInfo_Reply
}

func addAsContact(ctx context.Context, t *testing.T, senders, receivers []*bertyprotocol.TestingProtocol) {
	t.Log(logTree("Add Senders/Receivers as Contact", 0, true))
	start := time.Now()
	var sendDuration, receiveDuration, acceptDuration, activateDuration time.Duration

	for i, sender := range senders {
		for _, receiver := range receivers {
			substart := time.Now()

			// Get sender/receiver configs
			senderCfg, err := sender.Client.InstanceGetConfiguration(ctx, &bertytypes.InstanceGetConfiguration_Request{})
			require.NoError(t, err)
			require.NotNil(t, senderCfg)
			receiverCfg, err := receiver.Client.InstanceGetConfiguration(ctx, &bertytypes.InstanceGetConfiguration_Request{})
			require.NoError(t, err)
			require.NotNil(t, receiverCfg)

			// Setup receiver's sharable contact
			_, err = receiver.Client.ContactRequestEnable(ctx, &bertytypes.ContactRequestEnable_Request{})
			require.NoError(t, err)
			receiverRDV, err := receiver.Client.ContactRequestResetReference(ctx, &bertytypes.ContactRequestResetReference_Request{})
			require.NoError(t, err)
			require.NotNil(t, receiverRDV)
			if i == 1 {
				break
			}

			receiverSharableContact := &bertytypes.ShareableContact{
				PK:                   receiverCfg.AccountPK,
				PublicRendezvousSeed: receiverRDV.PublicRendezvousSeed,
			}

			// Sender sends contact request
			_, err = sender.Client.ContactRequestSend(ctx, &bertytypes.ContactRequestSend_Request{
				Contact: receiverSharableContact,
			})

			// Check if sender and receiver are the same account, should return the right error and skip
			if bytes.Compare(senderCfg.AccountPK, receiverCfg.AccountPK) == 0 {
				require.Equal(t, errcode.LastCode(err), errcode.ErrContactRequestSameAccount)
				continue
			}

			// Check if contact request was already sent, should return right error and skip
			receiverWasSender := false
			for j := 0; j < i; j++ {
				if senders[j] == receiver {
					receiverWasSender = true
				}
			}

			senderWasReceiver := false
			if receiverWasSender {
				for _, r := range receivers {
					if r == sender {
						senderWasReceiver = true
					}
				}
			}

			if receiverWasSender && senderWasReceiver {
				require.Equal(t, errcode.LastCode(err), errcode.ErrContactRequestContactAlreadyAdded)
				continue
			}

			// No other error should occur
			require.NoError(t, err)

			sendDuration += time.Since(substart)
			substart = time.Now()

			// Receiver subcribes to handle incoming contact request
			subCtx, subCancel := context.WithCancel(ctx)
			subReceiver, err := receiver.Client.GroupMetadataSubscribe(subCtx, &bertytypes.GroupMetadataSubscribe_Request{
				GroupPK: receiverCfg.AccountGroupPK,
				Since:   []byte("give me everything"),
			})
			require.NoError(t, err)
			found := false

			// Receiver waits for valid contact request coming from sender
			for {
				evt, err := subReceiver.Recv()
				if err == io.EOF || subReceiver.Context().Err() != nil {
					break
				}

				require.NoError(t, err)

				if evt == nil || evt.Metadata.EventType != bertytypes.EventTypeAccountContactRequestIncomingReceived {
					continue
				}

				req := &bertytypes.AccountContactRequestReceived{}
				err = req.Unmarshal(evt.Event)

				require.NoError(t, err)

				if bytes.Compare(senderCfg.AccountPK, req.ContactPK) == 0 {
					found = true
					break
				}
			}

			subCancel()
			require.True(t, found)

			receiveDuration += time.Since(substart)
			substart = time.Now()

			// Receiver accepts contact request
			_, err = receiver.Client.ContactRequestAccept(ctx, &bertytypes.ContactRequestAccept_Request{
				ContactPK: senderCfg.AccountPK,
			})

			require.NoError(t, err)

			acceptDuration += time.Since(substart)
			substart = time.Now()

			// Both receiver and sender activate the contact group
			grpInfo, err := sender.Client.GroupInfo(ctx, &bertytypes.GroupInfo_Request{
				ContactPK: receiverCfg.AccountPK,
			})
			require.NoError(t, err)

			_, err = sender.Client.ActivateGroup(ctx, &bertytypes.ActivateGroup_Request{
				GroupPK: grpInfo.Group.PublicKey,
			})

			require.NoError(t, err)

			_, err = receiver.Client.ActivateGroup(ctx, &bertytypes.ActivateGroup_Request{
				GroupPK: grpInfo.Group.PublicKey,
			})

			require.NoError(t, err)

			activateDuration += time.Since(substart)
			substart = time.Now()
		}
	}

	t.Log(logTree("Send Contact Requests", 1, true))
	t.Logf(logTree("duration: %s", 1, false), sendDuration)
	t.Log(logTree("Receive Contact Requests", 1, true))
	t.Logf(logTree("duration: %s", 1, false), receiveDuration)
	t.Log(logTree("Accept Contact Requests", 1, true))
	t.Logf(logTree("duration: %s", 1, false), acceptDuration)
	t.Log(logTree("Activate Contact Groups", 1, true))
	t.Logf(logTree("duration: %s", 1, false), activateDuration)

	t.Logf(logTree("duration: %s", 0, false), time.Since(start))
}

func startMockedService(ctx context.Context, t *testing.T, logger *zap.Logger, amount int) ([]*BertyClient, func()) {
	opts := &bertyprotocol.TestingOpts{
		Mocknet:        libp2p_mocknet.New(ctx),
		TracerProvider: global.TraceProvider(),
	}
	rdvpeer, err := opts.Mocknet.GenPeer()
	require.NoError(t, err)
	require.NotNil(t, rdvpeer)

	_, cleanupRDVP := ipfsutil.TestingRDVP(ctx, t, rdvpeer)
	rdvpnet := opts.Mocknet.Net(rdvpeer.ID())
	require.NotNil(t, rdvpnet)

	opts.RDVPeer = rdvpeer.Peerstore().PeerInfo(rdvpeer.ID())

	tps := make([]*BertyClient, amount)
	for i := range tps {
		tps[i] = &BertyClient{}
		svcName := fmt.Sprintf("pt[%d]", i)
		opts.Logger = logger.Named(svcName)

		tps[i].Protocol, tps[i].cancel = bertyprotocol.NewTestingProtocol(ctx, t, opts)
		require.NotNil(t, tps[i])

		tps[i].Messenger, err = New(tps[i].Protocol.Client, &Opts{Logger: logger.Named("messenger")})
		if err != nil {
			cleanupRDVP()
		}
	}

	err = opts.Mocknet.LinkAll()
	require.NoError(t, err)

	for _, net := range opts.Mocknet.Nets() {
		if net != rdvpnet {
			_, err = opts.Mocknet.ConnectNets(net, rdvpnet)
			assert.NoError(t, err)
		}
	}
	return tps, cleanupRDVP
}

func startBertyService(t *testing.T, logger *zap.Logger) *BertyClient {
	t.Log("Starting service")
	var (
		node *core.IpfsNode
		api  ipfsutil.ExtendedCoreAPI
		disc tinder.Driver
		ps   *pubsub.PubSub
	)
	ctx := context.Background()
	rdvpeer, err := parseRdvpMaddr(ctx, "/dnsaddr/rdvp.berty.io/ipfs/QmdT7AmhhnbuwvCpa5PH1ySK9HJVB82jr3fo1bxMxBPW6p", logger)
	if err != nil {
		require.NoError(t, err)
	}
	opts := &ipfsutil.CoreAPIConfig{
		DisableCorePubSub: true,
		BootstrapAddrs:    config.DefaultBootstrapAddresses,
		SwarmAddrs: []string{
			"/ip4/0.0.0.0/tcp/0",
			"/ip6/0.0.0.0/tcp/0",
			"/ip4/0.0.0.0/udp/0/quic",
			"/ip6/0.0.0.0/udp/0/quic",
		},
		HostConfig: func(h host.Host, _ routing.Routing) error {
			var err error

			h.Peerstore().AddAddrs(rdvpeer.ID, rdvpeer.Addrs, peerstore.PermanentAddrTTL)
			rdvClient := tinder.NewRendezvousDiscovery(logger, h, rdvpeer.ID,
				rand.New(rand.NewSource(rand.Int63())))

			minBackoff, maxBackoff := time.Second, time.Minute
			rng := rand.New(rand.NewSource(rand.Int63()))
			disc, err = tinder.NewService(
				logger,
				rdvClient,
				discovery.NewExponentialBackoff(minBackoff, maxBackoff, discovery.FullJitter, time.Second, 5.0, 0, rng),
			)
			if err != nil {
				return err
			}

			ps, err = pubsub.NewGossipSub(ctx, h,
				pubsub.WithMessageSigning(true),
				pubsub.WithFloodPublish(true),
				pubsub.WithDiscovery(disc),
				pubsub.WithPeerExchange(true),
			)
			if err != nil {
				return err
			}

			return nil
		},
	}
	require.NoError(t, err)

	psapi := ipfsutil.NewPubSubAPI(ctx, logger.Named("pubsub"), disc, ps)
	api, node, err = ipfsutil.NewCoreAPI(ctx, opts)
	require.NoError(t, err)

	api = ipfsutil.InjectPubSubCoreAPIExtendedAdaptater(api, psapi)
	require.NoError(t, err)

	ipfsutil.EnableConnLogger(ctx, logger, node.PeerHost)

	rootDS := sync_ds.MutexWrap(datastore.NewMapDatastore())
	mk := bertyprotocol.NewMessageKeystore(ipfsutil.NewNamespacedDatastore(rootDS, datastore.NewKey("messages")))
	ks := ipfsutil.NewDatastoreKeystore(ipfsutil.NewNamespacedDatastore(rootDS, datastore.NewKey("account")))
	orbitdbDS := ipfsutil.NewNamespacedDatastore(rootDS, datastore.NewKey("orbitdb"))
	protoOpts := &bertyprotocol.Opts{
		Logger:          logger.Named("protocol"),
		PubSub:          ps,
		TinderDriver:    disc,
		IpfsCoreAPI:     api,
		DeviceKeystore:  bertyprotocol.NewDeviceKeystore(ks),
		RootDatastore:   rootDS,
		MessageKeystore: mk,
		OrbitCache:      bertyprotocol.NewOrbitDatastoreCache(orbitdbDS),
		Host:            node.PeerHost,
	}
	service, err := bertyprotocol.New(ctx, *protoOpts)

	if err != nil {
		require.NoError(t, err)
	}

	grpcServer := grpc.NewServer()

	client, err := bertyprotocol.NewClientFromServer(ctx, grpcServer, service)
	if err != nil {
		require.NoError(t, err)
	}

	messenger, err := New(client, &Opts{Logger: logger.Named("messenger")})
	if err != nil {
		defer node.Close()
		defer service.Close()
		defer client.Close()
	}

	return &BertyClient{
		Protocol: &bertyprotocol.TestingProtocol{
			Opts:    protoOpts,
			Client:  client,
			Service: service,
		},
		Messenger: messenger,
		cancel: func() {
			defer node.Close()
			defer service.Close()
			defer client.Close()
		},
	}
}

func parseRdvpMaddr(ctx context.Context, rdvpMaddr string, logger *zap.Logger) (*peer.AddrInfo, error) {
	if rdvpMaddr == "" {
		logger.Debug("no rendezvous peer set")
		return nil, nil
	}

	resoveCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	rdvpeer, err := ipfsutil.ParseAndResolveIpfsAddr(resoveCtx, rdvpMaddr)
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	fds := make([]zapcore.Field, len(rdvpeer.Addrs))
	for i, maddr := range rdvpeer.Addrs {
		key := fmt.Sprintf("#%d", i)
		fds[i] = zap.String(key, maddr.String())
	}
	logger.Debug("rdvp peer resolved addrs", fds...)
	return rdvpeer, nil
}

func TestScenario(t *testing.T) {
	ctx := context.Background()

	logger, err := zap.NewDevelopment()
	if err != nil {
		require.NoError(t, err)
	}

	flush := tracer.InitTracer("localhost:14268", "berty")
	defer flush()

	num := 2
	//clients := make([]*BertyClient, num)
	// Start Mocked protocol
	clients, cleanup := startMockedService(ctx, t, logger, num)
	defer cleanup()

	/*clientsT := make([]*bertyprotocol.TestingProtocol, num)
	clientsT[0] = clients[0].Protocol
	clientsT[1] = clients[1].Protocol
	addAsContact(ctx, t, clientsT, clientsT)*/

	for i := range clients {
		// Start real protocol
		//clients[i] = startBertyService(t, logger)

		clients[i].config, err = clients[i].Protocol.Client.InstanceGetConfiguration(ctx, &bertytypes.InstanceGetConfiguration_Request{})
		if err != nil {
			require.NoError(t, err)
		}
	}

	t.Log("ShareableBertyID")
	_, err = clients[0].Protocol.Client.ContactRequestEnable(ctx, &bertytypes.ContactRequestEnable_Request{})
	require.NoError(t, err)
	receiverRDV, err := clients[0].Protocol.Client.ContactRequestResetReference(ctx, &bertytypes.ContactRequestResetReference_Request{})
	require.NoError(t, err)
	require.NotNil(t, receiverRDV)

	contact := &bertytypes.ShareableContact{
		PK:                   clients[0].config.AccountPK,
		PublicRendezvousSeed: receiverRDV.PublicRendezvousSeed,
	}

	t.Log("Send ContactRequest")
	_, err = clients[1].Protocol.Client.ContactRequestSend(ctx, &bertytypes.ContactRequestSend_Request{
		Contact:     contact,
		OwnMetadata: []byte("client[1]"),
	})
	require.NoError(t, err)

	subscribeMetaDataEvents(t, ctx, clients[0])

	// Activate the contact group
	t.Log("GroupInfo query")
	clients[1].group, err = clients[1].Protocol.Client.GroupInfo(ctx, &bertytypes.GroupInfo_Request{
		ContactPK: clients[0].config.AccountPK,
	})
	require.NoError(t, err)

	clients[0].group, err = clients[0].Protocol.Client.GroupInfo(ctx, &bertytypes.GroupInfo_Request{
		ContactPK: clients[1].config.AccountPK,
	})
	require.NoError(t, err)

	t.Log("Activate group1")
	_, err = clients[1].Protocol.Client.ActivateGroup(ctx, &bertytypes.ActivateGroup_Request{
		GroupPK: clients[1].group.Group.PublicKey,
	})
	require.NoError(t, err)

	t.Log("Activate group0")
	_, err = clients[0].Protocol.Client.ActivateGroup(ctx, &bertytypes.ActivateGroup_Request{
		GroupPK: clients[0].group.Group.PublicKey,
	})
	require.NoError(t, err)

	/*_, err = clients[0].Protocol.Client.ContactRequestDisable(ctx, &bertytypes.ContactRequestDisable_Request{})
	require.NoError(t, err)*/

	t.Log("Send message")
	time.Sleep(1 * time.Second)

	/*_, err = clients[1].Messenger.SendMessage(ctx, &SendMessage_Request{
		GroupPK: clients[1].group.Group.PublicKey,
		Message: "test",
	})
	require.NoError(t, err)*/

	wg := sync.WaitGroup{}
	wg.Add(1)
	subscribeMessageEvents(t, ctx, clients[0], &wg)

	_, err = clients[1].Protocol.Client.AppMessageSend(ctx, &bertytypes.AppMessageSend_Request{
		GroupPK: clients[1].group.Group.PublicKey,
		Payload: []byte("test"),
	})
	require.NoError(t, err)

	wg.Wait()

	t.Log("Send message")
	_, err = clients[0].Messenger.SendMessage(ctx, &SendMessage_Request{
		GroupPK: clients[1].group.Group.PublicKey,
		Message: "test2",
	})
	require.NoError(t, err)

}

func subscribeMetaDataEvents(t *testing.T, ctx context.Context, client *BertyClient) {
	t.Log("Subscribe MetaDataEvents")

	var evt *bertytypes.GroupMetadataEvent

	req := &bertytypes.GroupMetadataSubscribe_Request{GroupPK: client.config.AccountGroupPK}
	cl, err := client.Protocol.Client.GroupMetadataSubscribe(ctx, req)
	if err != nil {
		require.NoError(t, err)
	}

	for {
		evt, err = cl.Recv()
		if err != nil {
			require.NoError(t, err)
		}

		if evt.Metadata.EventType == bertytypes.EventTypeAccountContactRequestIncomingReceived {
			t.Log("ContactRequest received")
			casted := &bertytypes.AccountContactRequestReceived{}
			if err := casted.Unmarshal(evt.Event); err != nil {
				require.NoError(t, err)
			}
			_, err = client.Protocol.Client.ContactRequestAccept(ctx, &bertytypes.ContactRequestAccept_Request{
				ContactPK: casted.ContactPK,
			})
			if err != nil {
				require.NoError(t, err)
			}
			return
		}
	}
}

func subscribeMessageEvents(t *testing.T, ctx context.Context, receiver *BertyClient, wg *sync.WaitGroup) {
	t.Log("Subscribe MessageEvents")

	var evt *bertytypes.GroupMessageEvent

	req := &bertytypes.GroupMessageSubscribe_Request{
		GroupPK: receiver.group.Group.PublicKey,
		//Since:   []byte("give me everything"),
	}
	cl, err := receiver.Protocol.Client.GroupMessageSubscribe(ctx, req)
	require.NoError(t, err)

	go func() {
		defer wg.Done()
		for {
			evt, err = cl.Recv()
			t.Log("Message Event found")
			if err == io.EOF {
				return
			} else if err != nil {
				require.NoError(t, err)
			}
			require.Equal(t, "test", string(evt.Message))
			return
		}
	}()
}

func logTree(log string, indent int, title bool) string {
	if !title {
		log = "└── " + log
	}

	for i := 0; i < indent; i++ {
		log = "│  " + log
	}

	return log
}
