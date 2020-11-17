package bertybridge

import (
	proximity "berty.tech/berty/v2/go/internal/proximity-transport"
)

type NativeBleDriver interface {
	proximity.NativeDriver
}

type ProximityTransport struct {
	proximity.ProximityTransport
}

func GetProximityTransport(protocolName string) *ProximityTransport {
	t, ok := proximity.TransportMap.Load(protocolName)
	if ok {
		return t.(*ProximityTransport)
	}
	return nil
}
