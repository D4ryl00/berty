// +build android

package ble

import (
	"go.uber.org/zap"

	proximity "berty.tech/berty/v2/go/internal/proximity-transport"
)

// Supported is used by main package as default value for enable the BLE  driver.
// While UI actually enable or not the Java BLE driver.
// TODO: remove this when UI will be able to handle this for the first App launching.
const Supported = true

// Noop implementation for Android
// Real driver is given from Java directly
func NewDriver(logger *zap.Logger) proximity.NativeDriver {
	logger = logger.Named("BLE")
	logger.Info("NewDriver(): Java driver not found")

	return proximity.NewNoopNativeDriver(ProtocolCode, ProtocolName, DefaultAddr)
}
