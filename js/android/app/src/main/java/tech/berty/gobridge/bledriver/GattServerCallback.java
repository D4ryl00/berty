package tech.berty.gobridge.bledriver;

import android.bluetooth.BluetoothDevice;
import android.bluetooth.BluetoothGatt;
import android.bluetooth.BluetoothGattCharacteristic;
import android.bluetooth.BluetoothGattServerCallback;
import android.bluetooth.BluetoothGattService;
import android.bluetooth.BluetoothProfile;
import android.content.Context;
import android.util.Log;

import java.util.Arrays;
import java.util.Base64;
import java.util.concurrent.CountDownLatch;

import static android.bluetooth.BluetoothGatt.GATT_SUCCESS;

public class GattServerCallback extends BluetoothGattServerCallback {
    private static final String TAG = "bty.ble.GattSrvCallback";

    // Size in bytes of the ATT MTU headers
    // see Bluetooth Core Specification 5.1: 4.8 Characteristic Value Read (p.2380)
    private static final int ATT_HEADER_READ_SIZE = 1;

    private Context mContext;
    private GattServer mGattServer;
    private CountDownLatch mDoneSignal;
    private String mLocalPID;

    public GattServerCallback(Context context, GattServer gattServer, CountDownLatch doneSignal) {
        mContext = context;
        mGattServer = gattServer;
        mDoneSignal = doneSignal;
    }

    public void setLocalPID(String peerID) {
        mLocalPID = peerID;
    }

    // When this callback is called, we assume that the GATT server is ready up.
    // We can enable scanner and advertiser.
    @Override
    public void onServiceAdded(int status, BluetoothGattService service) {
        Log.i(TAG, "onServiceAdded() called");
        super.onServiceAdded(status, service);
        if (status != BluetoothGatt.GATT_SUCCESS) {
            Log.e(TAG, "onServiceAdded error: failed to add service " + service);
        } else {
            mGattServer.setStarted(true);
        }
        mDoneSignal.countDown();
    }

    @Override
    public void onConnectionStateChange(BluetoothDevice device, int status, int newState) {
        super.onConnectionStateChange(device, status, newState);

        Log.v(TAG, String.format("onConnectionStateChange: device=%s status=%d newState=%d", device, status, newState));
        PeerDevice peerDevice = DeviceManager.get(device.getAddress());

        if (status == GATT_SUCCESS) {
            if (newState == BluetoothProfile.STATE_CONNECTED) {
                Log.d(TAG, "connected");

                if (peerDevice == null) {
                    Log.i(TAG, String.format("onConnectionStateChange(): connected with new device %s", device.getAddress()));
                    peerDevice = new PeerDevice(mContext, device, mLocalPID);
                    DeviceManager.addDevice(peerDevice);
                }

                if (peerDevice.getClientState() != PeerDevice.CONNECTION_STATE.DISCONNECTED) {
                    Log.d(TAG, String.format("onConnectionStateChange: connection already made as client, cancel this one for device %s", device.getAddress()));
                    //mGattServer.getGattServer().cancelConnection(device);
                } else if (!peerDevice.checkAndSetServerState(PeerDevice.CONNECTION_STATE.DISCONNECTED, PeerDevice.CONNECTION_STATE.CONNECTED)) {
                    Log.d(TAG, String.format("onConnectionStateChange: a server connection already exists, cancel this one for device %s", device.getAddress()));
                    //mGattServer.getGattServer().cancelConnection(device);
                }
            } else {
                Log.d(TAG, String.format("onConnectionStateChange: device %s disconnected", device.getAddress()));
                if (peerDevice != null) {
                    if (peerDevice.getServerState() == PeerDevice.CONNECTION_STATE.CONNECTED) {
                        BleInterface.BLEHandleLostPeer(peerDevice.getRemotePID());
                    }
                    peerDevice.setServerState(PeerDevice.CONNECTION_STATE.DISCONNECTED);
                }
            }
        } else {
            Log.e(TAG, String.format("onConnectionStateChange error with status %d", status));
        }
    }

    // onCharacteristicReadRequest is called when client wants the server device peer id
    @Override
    public void onCharacteristicReadRequest(BluetoothDevice device,
                                            int requestId,
                                            int offset,
                                            BluetoothGattCharacteristic characteristic) {
        super.onCharacteristicReadRequest(device, requestId, offset, characteristic);
        Log.v(TAG, String.format("onCharacteristicReadRequest() called: device=%s requestId=%d", device.getAddress(), requestId));

        boolean full = false;
        PeerDevice peerDevice;
        byte[] value;

        if ((peerDevice = DeviceManager.get(device.getAddress())) == null) {
            Log.e(TAG, String.format("onCharacteristicReadRequest(): device %s not found", device.getAddress()));
            /*peerDevice = new PeerDevice(mContext, device, mLocalPID);
            DeviceManager.addDevice(peerDevice);*/
            mGattServer.getGattServer().sendResponse(device, requestId, BluetoothGatt.GATT_FAILURE,
                offset, null);
            return ;
        } else if (peerDevice.getServerState() != PeerDevice.CONNECTION_STATE.CONNECTED) {
            Log.e(TAG, String.format("onCharacteristicWriteRequest: device not connected", device.getAddress()));
            mGattServer.getGattServer().sendResponse(device, requestId, BluetoothGatt.GATT_FAILURE,
                offset, null);
            return ;
        }
        if (characteristic.getUuid().equals(GattServer.READER_UUID)) {
            String peerID = characteristic.getStringValue(0);
            if ((peerID.length() - offset) <= peerDevice.getMtu() - ATT_HEADER_READ_SIZE) {
                Log.d(TAG, "onCharacteristicReadRequest: mtu is big enough (" + (peerID.length() - offset) + " bytes to read)");
                full = true;
            } else {
                Log.d(TAG, "onCharacteristicReadRequest: mtu is too small (" + (peerID.length() - offset) + " bytes to read)");
            }
            value = Arrays.copyOfRange(peerID.getBytes(), offset, peerID.length());
            mGattServer.getGattServer().sendResponse(device, requestId, BluetoothGatt.GATT_SUCCESS, offset, value);
            if (full) {
                Log.d(TAG, "onCharacteristicReadRequest: finished");
                peerDevice.handleServerDataSent();
            }
        } else {
            Log.e(TAG, "onCharacteristicReadRequest: try to read to a wrong characteristic");
            mGattServer.getGattServer().sendResponse(device, requestId, BluetoothGatt.GATT_FAILURE,
                    offset, null);
        }
    }

    // When receiving data, there are two cases:
    // * MTU is big enough, thus the whole message is transmitted, prepareWrite is false.
    // * Otherwise, we need to wait that all packets are transmitted, prepareWrite is true for
    //   all this transmissions. Data packets are put in a buffer.
    //   When all packets are sent, onExecuteWrite is called.
    @Override
    public void onCharacteristicWriteRequest(BluetoothDevice device,
                                             int requestId,
                                             BluetoothGattCharacteristic characteristic,
                                             boolean prepareWrite,
                                             boolean responseNeeded,
                                             int offset,
                                             byte[] value) {
        super.onCharacteristicWriteRequest(device, requestId, characteristic, prepareWrite,
                responseNeeded, offset, value);
        Log.v(TAG, String.format("onCharacteristicWriteRequest called: device=%s requestId=%d", device.getAddress(), requestId));
        PeerDevice peerDevice;
        boolean status = false;

        if ((peerDevice = DeviceManager.get(device.getAddress())) == null) {
            Log.e(TAG, String.format("onCharacteristicWriteRequest: device %s not found", device.getAddress()));
        } else if (peerDevice.getServerState() != PeerDevice.CONNECTION_STATE.CONNECTED) {
            Log.e(TAG, String.format("onCharacteristicWriteRequest: device not connected", device.getAddress()));
        } else {
            if (characteristic.getUuid().equals(GattServer.WRITER_UUID)) {
                Log.d(TAG, String.format("onCharacteristicWriteRequest: device=%s value=%s base64=%s size=%d offset=%d preparedWrite=%b needResponse=%b", device.getAddress(), BleDriver.bytesToHex(value), Base64.getEncoder().encodeToString(value), value.length, offset, prepareWrite, responseNeeded));
                if (prepareWrite) {
                    peerDevice.addToBuffer(value);
                    status = true;
                } else {
                    status = peerDevice.handleServerDataReceived(value); // TODO: copy value
                }
            } else {
                Log.e(TAG, "onCharacteristicWriteRequest: try to write to a wrong characteristic");
            }
        }
        if (responseNeeded && status) {
            mGattServer.getGattServer().sendResponse(device, requestId, BluetoothGatt.GATT_SUCCESS,
                    offset, value);
        } else if (responseNeeded) {
            mGattServer.getGattServer().sendResponse(device, requestId, BluetoothGatt.GATT_FAILURE,
                    0, null);
        }
    }

    // This callback is called when this GATT server has received all incoming data packets of one
    // transmission.
    // Thus we know we can handle data put in the buffer.
    @Override
    public void onExecuteWrite(BluetoothDevice device, int requestId, boolean execute) {
        super.onExecuteWrite(device, requestId, execute);
        Log.v(TAG, String.format("onExecuteWrite called: device=%s requestId=%d", device.getAddress(), requestId));
        PeerDevice peerDevice;
        boolean status = true;

        if (execute) {
            if ((peerDevice = DeviceManager.get(device.getAddress())) == null) {
                Log.e(TAG, String.format("onExecuteWrite: device %s not found", device.getAddress()));
                status = false;
            } else if (peerDevice.getServerState() != PeerDevice.CONNECTION_STATE.CONNECTED) {
                Log.e(TAG, String.format("onExecuteWrite: device not connected", device.getAddress()));
            } else {
                Log.v(TAG, String.format("onExecuteWrite: device=%s base64=%s value=%s", device.getAddress(), Base64.getEncoder().encodeToString(peerDevice.getBuffer()), BleDriver.bytesToHex(peerDevice.getBuffer())));
                status = peerDevice.handleServerDataReceived(peerDevice.getBuffer());
            }
            peerDevice.resetBuffer();
        }
        if (!status) {
            Log.e(TAG, "onExecuteWrite error");
            mGattServer.getGattServer().sendResponse(device, requestId, BluetoothGatt.GATT_FAILURE,
                0, null);
        } else {
            Log.d(TAG, String.format("onExecuteWrite successful for device %s", device.getAddress()));
            mGattServer.getGattServer().sendResponse(device, requestId, BluetoothGatt.GATT_SUCCESS,
                0, null);
        }
    }

    @Override
    public void onMtuChanged(BluetoothDevice device, int mtu) {
        super.onMtuChanged(device, mtu);
        Log.v(TAG, String.format("onMtuChanged called for device %s and mtu=%d", device.getAddress(), mtu));
        PeerDevice peerDevice;

        if ((peerDevice = DeviceManager.get(device.getAddress())) == null) {
            Log.e(TAG, "onMtuChanged() error: device not found");
            return ;
        }
        peerDevice.setMtu(mtu);
    }

    @Override
    public void onNotificationSent(BluetoothDevice device, int status) {
        super.onNotificationSent(device, status);
        Log.v(TAG, String.format("onNotificationSent called for device %s", device.getAddress()));

        if (status != GATT_SUCCESS) {
            Log.e(TAG, String.format("onNotificationSent: status error: %d for device %s, status", status, device));
        }
        BleQueue.completedCommand();
    }
}
