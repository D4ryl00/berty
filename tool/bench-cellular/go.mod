module berty.tech/berty/v2/tool/bench-cellular

go 1.16

require (
	github.com/ipfs/go-datastore v0.5.1
	github.com/ipfs/go-ipfs v0.11.0
	github.com/ipfs/go-ipfs-config v0.18.0
	github.com/ipfs/go-log/v2 v2.3.0
	github.com/libp2p/go-libp2p v0.16.0
	github.com/libp2p/go-libp2p-circuit v0.4.0
	github.com/libp2p/go-libp2p-connmgr v0.2.4
	github.com/libp2p/go-libp2p-core v0.11.0
	github.com/libp2p/go-libp2p-kad-dht v0.15.0
	github.com/libp2p/go-libp2p-quic-transport v0.15.0
	github.com/libp2p/go-tcp-transport v0.4.0
	github.com/multiformats/go-multiaddr v0.4.1
	github.com/peterbourgon/ff/v3 v3.0.0
)

replace (
	github.com/ipfs/go-ipfs => ./go-ipfs
	github.com/libp2p/go-libp2p => ./go-libp2p
	github.com/libp2p/go-libp2p-circuit => ./go-libp2p-circuit
)
