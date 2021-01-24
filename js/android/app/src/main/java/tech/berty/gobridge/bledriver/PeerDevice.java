package tech.berty.gobridge.bledriver;

import android.bluetooth.BluetoothDevice;
import android.bluetooth.BluetoothGatt;
import android.bluetooth.BluetoothGattCallback;
import android.bluetooth.BluetoothGattCharacteristic;
import android.bluetooth.BluetoothGattService;
import android.bluetooth.BluetoothProfile;
import android.content.Context;
import android.os.Build;
import android.os.Handler;
import android.util.Log;

import androidx.annotation.NonNull;

import java.util.Base64;
import java.util.List;
import java.util.Queue;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

import static android.bluetooth.BluetoothDevice.BOND_BONDED;
import static android.bluetooth.BluetoothDevice.BOND_BONDING;
import static android.bluetooth.BluetoothDevice.BOND_NONE;
import static android.bluetooth.BluetoothGatt.GATT_SUCCESS;

public class PeerDevice {
    private static final String TAG = "bty.ble.PeerDevice";

    // Minimal and default MTU
    private static final int DEFAULT_MTU = 23;

    // Max MTU that Android can handle
    public static final int MAX_MTU = 517;

    public enum CONNECTION_STATE {
        DISCONNECTED,
        CONNECTED,
        CONNECTING,
        DISCONNECTING
    }
    private CONNECTION_STATE mClientState = CONNECTION_STATE.DISCONNECTED;
    private CONNECTION_STATE mServerState = CONNECTION_STATE.DISCONNECTED;


    private Context mContext;
    private BluetoothDevice mBluetoothDevice;
    private BluetoothGatt mBluetoothGatt;

    public final Object mLockClientState = new Object();
    public final Object mLockServerState = new Object();
    private final Object mLockRemotePID = new Object();
    private final Object mLockMtu = new Object();
    private final Object mLockClient = new Object();
    private final Object mLockServer = new Object();
    private final Object mLockPeer = new Object();

    private BluetoothGattService mBertyService;
    private BluetoothGattCharacteristic mReaderCharacteristic;
    private BluetoothGattCharacteristic mWriterCharacteristic;

    private Peer mPeer;
    private String mRemotePID;
    private String mLocalPID;
    private boolean mClientReady = false;
    private boolean mServerReady = false;
    private boolean mDiscoveryStarted = false;

    //private int mMtu = 0;
    // default MTU is 23
    private int mMtu = 23;

    public PeerDevice(@NonNull Context context, @NonNull BluetoothDevice bluetoothDevice, String localPID) {
        mContext = context;
        mBluetoothDevice = bluetoothDevice;
        mLocalPID = localPID;
    }

    public BluetoothDevice getBluetoothDevice() {
        return mBluetoothDevice;
    }

    public String getMACAddress() {
        return mBluetoothDevice.getAddress();
    }

    @NonNull
    @Override
    public java.lang.String toString() {
        return getMACAddress();
    }

    // Use TRANSPORT_LE for connections to remote dual-mode devices. This is a solution to prevent the error
    // status 133 in GATT connections:
    // https://android.jlelse.eu/lessons-for-first-time-android-bluetooth-le-developers-i-learned-the-hard-way-fee07646624
    // API level 23
    public void connectToDevice() {
        Log.d(TAG, "connectToDevice: " + getMACAddress());

        if (getServerState() != CONNECTION_STATE.DISCONNECTED) {
            Log.d(TAG, String.format("connectToDevice canceled, device %s is handled by GATT server", getMACAddress()));
        } else if (checkAndSetClientState(CONNECTION_STATE.DISCONNECTED, CONNECTION_STATE.CONNECTING)) {
            BleDriver.mainHandler.postDelayed(new Runnable() {
                @Override
                public void run() {
                    setBluetoothGatt(mBluetoothDevice.connectGatt(mContext, false,
                        mGattCallback, BluetoothDevice.TRANSPORT_LE));
                }
            }, 100);
        } else {
            Log.d(TAG, String.format("connectToDevice canceled, device %s is not disconnected", getMACAddress()));
        }
    }

    public boolean isConnected() {
        return getClientState() == CONNECTION_STATE.CONNECTED || getServerState() == CONNECTION_STATE.CONNECTED;
    }

    public boolean isDisconnected() {
        return getClientState() == CONNECTION_STATE.DISCONNECTED && getServerState() == CONNECTION_STATE.DISCONNECTED;
    }

    private void disconnect() {
        synchronized (mLockClientState) {
            if (mClientState == CONNECTION_STATE.CONNECTED || mClientState == CONNECTION_STATE.CONNECTING) {
                mClientState = CONNECTION_STATE.DISCONNECTING;
                BleDriver.mainHandler.post(new Runnable() {
                    @Override
                    public void run() {
                        synchronized (mLockClientState) {
                            if (mClientState == CONNECTION_STATE.DISCONNECTING && getBluetoothGatt() != null) {
                                getBluetoothGatt().disconnect();
                                Log.i(TAG, String.format("force disconnect %s", mBluetoothDevice.getAddress()));
                            }
                        }
                    }
                });
            }
        }
    }

    public void close() {

        if (getBluetoothGatt() != null) {
            getBluetoothGatt().close();
            setBluetoothGatt(null);
        }
        BleQueue.clear();
        setClientReady(false);
        setServerReady(false);
        setPeer(null);
        PeerManager.unregister(mRemotePID);
    }

    // setters and getters are accessed by the DeviceManager thread et this thread so we need to
    // synchronize them.
    public void setClientState(CONNECTION_STATE state) {
        synchronized (mLockClientState) {
            mClientState = state;
        }
    }

    public CONNECTION_STATE getClientState() {
        synchronized (mLockClientState) {
            return mClientState;
        }
    }

    public boolean checkAndSetClientState(CONNECTION_STATE state, CONNECTION_STATE newState) {
        Log.v(TAG, String.format("checkAndSetClientState called for device %s, state=%s newState=%s", getMACAddress(), state, newState));
        synchronized (mLockClientState) {
            if (mClientState == state) {
                mClientState = newState;
                return true;
            }
            return false;
        }
    }

    public void setServerState(CONNECTION_STATE state) {
        synchronized (mLockServerState) {
            mServerState = state;
        }
    }

    public CONNECTION_STATE getServerState() {
        synchronized (mLockServerState) {
            return mServerState;
        }
    }

    public boolean checkAndSetServerState(CONNECTION_STATE state, CONNECTION_STATE newState) {
        Log.v(TAG, String.format("checkAndSetServerState called for device %s, state=%s newState=%s", getMACAddress(), state, newState));
        synchronized (mLockServerState) {
            if (mServerState == state) {
                mServerState = newState;
                return true;
            }
            return false;
        }
    }

    public synchronized void setBluetoothGatt(BluetoothGatt gatt) {
            mBluetoothGatt = gatt;
    }

    public synchronized BluetoothGatt getBluetoothGatt() {
            return mBluetoothGatt;
    }

    // isClient return if this PeerDevice is GATT client or server
    // A GATT client has a BluetoothGatt set, a server no.
    public boolean isClient() {
        synchronized (mLockClientState) {
            return mBluetoothGatt != null;
        }
    }

    public void setReaderCharacteristic(BluetoothGattCharacteristic peerID) {
        mReaderCharacteristic = peerID;
    }

    public BluetoothGattCharacteristic getReaderCharacteristic() {
        return mReaderCharacteristic;
    }

    public void setWriterCharacteristic(BluetoothGattCharacteristic write) {
        mWriterCharacteristic = write;
    }

    public BluetoothGattCharacteristic getWriterCharacteristic() {
        return mWriterCharacteristic;
    }

    public boolean setWriterValue(String value) {
        return getWriterCharacteristic().setValue(value);
    }

    public BluetoothGattService getBertyService() {
        return mBertyService;
    }

    public void setBertyService(BluetoothGattService service) {
        mBertyService = service;
    }

    public void setRemotePID(String peerID) {
        Log.d(TAG, "setPeerID: " + peerID + ", device: " + getMACAddress());
        synchronized (mLockRemotePID) {
            mRemotePID = peerID;
        }
    }

    public String getRemotePID() {
        synchronized (mLockRemotePID) {
            return mRemotePID;
        }
    }

    public void setPeer(Peer peer) {
        synchronized (mLockPeer) {
            mPeer = peer;
        }
    }

    public Peer getPeer() {
        synchronized (mLockPeer) {
            return mPeer;
        }
    }

    public void handleServerDataSent() {
        Log.v(TAG, String.format("handleServerDataSent for device %s", getMACAddress()));

        if (!isClientReady()) {
            if (getRemotePID() == null) {
                Log.e(TAG, String.format("handleServerDataSent: remotePID is null for device %s", getMACAddress()));
                return ;
            }

            setServerReady(true);
            setClientReady(true);
        } else {
            Log.e(TAG, String.format("handleServerDataSent: PID already sent for device %s", getMACAddress()));
            close();
        }
    }

    public boolean handleServerDataReceived(byte[] payload) {
        Log.v(TAG, String.format("handleServerDataReceived for device %s", getMACAddress()));

        boolean status = false;

        if (!isServerReady()) {
            Peer peer;
            String remotePID = new String(payload);

            // check if a connection already exists
            if ((peer = PeerManager.get(remotePID)) != null) {
                Log.i(TAG, String.format("handleServerDataReceived: canceling connection for device %s because a connection with the peer %s already exists for device %s", getMACAddress(), remotePID, peer.getPeerDevice().getMACAddress()));
                disconnect();
                return false;
            }

            setRemotePID(remotePID);
            peer = PeerManager.register(remotePID, this);
            setPeer(peer);
            status = true;
        } else {
            //BleInterface.BLEReceiveFromPeer(getRemotePID(), payload);
            status = BleQueue.add(new Runnable() {
                @Override
                public void run() {
                    Log.d(TAG, String.format("BleQueue: BLEReceiveFromPeer for device %s", getMACAddress()));
                    BleInterface.BLEReceiveFromPeer(getRemotePID(), payload);
                    BleQueue.completedCommand();
                }
            });

            if (status) {
                BleQueue.nextCommand();
            } else {
                Log.e(TAG, "could not enqueue requestMtu command");
            }
        }
        return status;
    }

    public void handleClientDataReceived(byte[] payload) {
        Log.v(TAG, String.format("handleClientDataReceived for device %s", getMACAddress()));

        if (!isClientReady()) {
            Peer peer;
            String remotePID = new String(payload);

            // check if a connection already exists
            if ((peer = PeerManager.get(remotePID)) != null) {
                Log.i(TAG, String.format("handleClientDataReceived: canceling connection for device %s because a connection with the peer %s already exists for device %s", getMACAddress(), remotePID, peer.getPeerDevice().getMACAddress()));
                disconnect();
                close();
                return ;
            }

            setRemotePID(remotePID);
            peer = PeerManager.register(remotePID, this);
            setPeer(peer);
            setClientReady(true);
            setServerReady(true);
            //peer.CallFoundPeer();
        } else {
            //BleInterface.BLEReceiveFromPeer(getRemotePID(), payload);
            boolean result = BleQueue.add(new Runnable() {
                @Override
                public void run() {
                    Log.d(TAG, String.format("BleQueue: BLEReceiveFromPeer for device %s", getMACAddress()));
                    BleInterface.BLEReceiveFromPeer(getRemotePID(), payload);
                    BleQueue.completedCommand();
                }
            });

            if (result) {
                BleQueue.nextCommand();
            } else {
                Log.e(TAG, "could not enqueue requestMtu command");
            }
        }
    }

    public void handleClientDataSent() {
    }

    public void handleWriteLocalPID() {
        Log.d(TAG, "handleWriteLocalPID called");
        setClientReady(true);
        PeerManager.register(mRemotePID, this);
    }

    public void setClientReady(boolean state) {
        Log.v(TAG, String.format("setClientReady called for device %s", getMACAddress()));

        synchronized (mLockClient) {
            mClientReady = state;
            Peer peer;

            if ((peer = PeerManager.get(getRemotePID())) != null) {
                Log.i(TAG, String.format("setClientReady: calling CallFoundPeer for device %s", getMACAddress()));
                peer.CallFoundPeer();
            }
        }
    }

    public boolean isClientReady() {
        synchronized (mLockClient) {
            return mClientReady;
        }
    }

    public void setServerReady(boolean state) {
        Log.v(TAG, String.format("setServerReady called for device %s", getMACAddress()));

        synchronized (mLockServer) {
            mServerReady = state;
            Peer peer;

            if ((peer = PeerManager.get(getRemotePID())) != null) {
                Log.i(TAG, String.format("setServerReady: calling CallFoundPeer for device %s", getMACAddress()));
                peer.CallFoundPeer();
            }
        }
    }

    public boolean isServerReady() {
        synchronized (mLockServer) {
            return mServerReady;
        }
    }

    private boolean takeBertyService(List<BluetoothGattService> services) {
        Log.v(TAG, String.format("takeBertyService: called for device %s", getMACAddress()));
        if (getBertyService() != null) {
            Log.d(TAG, String.format("Berty service already found for device %s", getMACAddress()));
            return true;
        }

        for (BluetoothGattService service : services) {
            if (service.getUuid().equals(GattServer.SERVICE_UUID)) {
                Log.d(TAG, String.format("Berty service found for device %s", getMACAddress()));
                setBertyService(service);
                return true;
            }
        }
        Log.i(TAG, String.format("Berty service not found for device", getMACAddress()));
        return false;
    }

    private boolean checkCharacteristicProperties(BluetoothGattCharacteristic characteristic,
                                                  int properties) {
        Log.d(TAG, "checkCharacteristicProperties: device: " + getMACAddress());

        if (characteristic.getProperties() == properties) {
            Log.d(TAG, "checkCharacteristicProperties() match, device: " + getMACAddress());
            return true;
        }
        Log.e(TAG, "checkCharacteristicProperties() doesn't match: " + characteristic.getProperties() + " / " + properties + ", device: " + getMACAddress());
        return false;
    }

    private boolean takeBertyCharacteristics() {
        Log.v(TAG, String.format("takeBertyCharacteristic called for device %s", getMACAddress()));

        if (getReaderCharacteristic() != null && getWriterCharacteristic() != null) {
            Log.d(TAG, "Berty characteristics already found");
            return true;
        }

        List<BluetoothGattCharacteristic> characteristics = mBertyService.getCharacteristics();
        for (BluetoothGattCharacteristic characteristic : characteristics) {
            if (characteristic.getUuid().equals(GattServer.READER_UUID)) {
                Log.d(TAG, String.format("reader characteristic found for device %s", getMACAddress()));
                if (checkCharacteristicProperties(characteristic,
                        BluetoothGattCharacteristic.PROPERTY_READ)) {
                    setReaderCharacteristic(characteristic);

                    if (!getBluetoothGatt().setCharacteristicNotification(characteristic, true)) {
                        Log.e(TAG, String.format("setCharacteristicNotification failed for device %s", getMACAddress()));
                        setReaderCharacteristic(null);
                        return false;
                    }
                }
            } else if (characteristic.getUuid().equals(GattServer.WRITER_UUID)) {
                Log.d(TAG, String.format("writer characteristic found for device: ", getMACAddress()));
                if (checkCharacteristicProperties(characteristic,
                        BluetoothGattCharacteristic.PROPERTY_WRITE)) {
                    setWriterCharacteristic(characteristic);
                }
            }
        }

        if (getReaderCharacteristic() != null && getWriterCharacteristic() != null) {
            return true;
        }

        Log.e(TAG, String.format("reader/writer characteristics not found for device %s", getMACAddress()));
        return false;
    }

    public boolean read() {
        Log.v(TAG, String.format("read() called for device %s", getMACAddress()));

        if (!isClient() || !isConnected()) {
            Log.e(TAG, String.format("read failed: device %s not connected", getMACAddress()));
            return false;
        }

        boolean result = BleQueue.add(new Runnable() {
            @Override
            public void run() {
                Log.v(TAG, String.format("BleQueue: read for device %s", getMACAddress()));
                if (isClient() && isConnected()) {
                    if (!getBluetoothGatt().readCharacteristic(getReaderCharacteristic())) {
                        Log.e(TAG, String.format("readCharacteristic failed for characteristic: %s", getReaderCharacteristic().getUuid()));
                        BleQueue.completedCommand();
                    } else {
                        Log.d(TAG, String.format("reading characteristic <%s>", getReaderCharacteristic().getUuid()));
                        //mNrTries++;
                    }
                } else {
                    Log.e(TAG, String.format("read failed: device %s not connected", getMACAddress()));
                    BleQueue.completedCommand();
                }
            }
        });

        if (result) {
            BleQueue.nextCommand();
        } else {
            Log.e(TAG, "could not enqueue read characteristic command");
        }
        return result;
    }

    public boolean write(byte[] payload) {
        Log.v(TAG, String.format("write() called for device %s", getMACAddress()));

        if (!isClient() || !isConnected()) {
            Log.e(TAG, "write failed: device not connected");
            return false;
        }

        Log.d(TAG, String.format("write: value is %s", Base64.getEncoder().encodeToString(payload)));

        boolean result = BleQueue.add(new Runnable() {
            @Override
            public void run() {
                Log.v(TAG, String.format("BleQueue: write for device %s", getMACAddress()));
                if (isClient() && isConnected()) {
                    if (!getWriterCharacteristic().setValue(payload) || !getBluetoothGatt().writeCharacteristic(getWriterCharacteristic())) {
                        Log.e(TAG, String.format("writerCharacteristic failed for characteristic: %s", getWriterCharacteristic().getUuid()));
                        BleQueue.completedCommand();
                    } else {
                        Log.d(TAG, String.format("writing characteristic %s", getWriterCharacteristic().getUuid()));
                        //mNrTries++;
                    }
                } else {
                    BleQueue.completedCommand();
                }
            }
        });

        if (result) {
            BleQueue.nextCommand();
        } else {
            Log.e(TAG, "could not enqueue read characteristic command");
        }
        return result;
    }

    private boolean requestMtu(final int mtu) {
        Log.v(TAG, "requestMtu called");

        if (mtu < DEFAULT_MTU || mtu > MAX_MTU) {
            Log.e(TAG, "mtu must be between 23 and 517");
            return false;
        }

        if (!isClient() || !isConnected()) {
            Log.e(TAG, "request mtu failed: device not connected");
            return false;
        }

        boolean result = BleQueue.add(new Runnable() {
            @Override
            public void run() {
                Log.v(TAG, String.format("BleQueue: requestMtu for device %s", getMACAddress()));
                if (isClient() && isConnected()) {
                    if (!getBluetoothGatt().requestMtu(mtu)) {
                        Log.e(TAG, "requestMtu failed");
                        BleQueue.completedCommand();
                    }
                } else {
                    Log.e(TAG, "request MTU failed: device not connected");
                    BleQueue.completedCommand();
                }
            }
        });

        if (result) {
            BleQueue.nextCommand();
        } else {
            Log.e(TAG, "could not enqueue requestMtu command");
        }

        return result;
    }

    public void setMtu(int mtu) {
        mMtu = mtu;
    }

    public int getMtu() {
        return mMtu;
    }

    private void handshake() {
        Log.d(TAG, "handshake: called");
        if (takeBertyService(getBluetoothGatt().getServices())) {
            if (takeBertyCharacteristics()) {
                //requestMtu(MAX_MTU);

                // send local PID
                if (!write(mLocalPID.getBytes())) {
                    Log.e(TAG, String.format("handshake: fail to send local PID for device %s", getMACAddress()));
                    disconnect();
                }

                // get remote PID
                if (!read()) {
                    Log.e(TAG, String.format("handshake: fail to read remote PID for device %s", getMACAddress()));
                    disconnect();
                }
            }
        }
    }

    private BluetoothGattCallback mGattCallback =
            new BluetoothGattCallback() {
                @Override
                public void onConnectionStateChange(BluetoothGatt gatt, int status, int newState) {
                    super.onConnectionStateChange(gatt, status, newState);
                    Log.v(TAG, "onConnectionStateChange() called by device " + gatt.getDevice().getAddress());
                    BluetoothDevice device = gatt.getDevice();

                    if (status == GATT_SUCCESS) {
                        if (newState == BluetoothProfile.STATE_CONNECTED) {
                            Log.d(TAG, "onConnectionStateChange(): connected");
                            setClientState(CONNECTION_STATE.CONNECTED);

                            int bondState = device.getBondState();        // Take action depending on the bond state
                            if(bondState == BOND_NONE || bondState == BOND_BONDED) {
                                // Connected to device, now proceed to discover it's services but delay a bit if needed
                                int delayWhenBonded = 0;
                                if (Build.VERSION.SDK_INT <= Build.VERSION_CODES.N) {
                                    delayWhenBonded = 1000;
                                }
                                final int delay = bondState == BOND_BONDED ? delayWhenBonded : 0;
                                final Runnable discoverRunnable = () -> {
                                    Log.d(TAG, "Continuing connection of " + device.getAddress() + " with delay of " + delay + " ms");
                                    if (gatt.discoverServices()) {
                                        mDiscoveryStarted = true;
                                    } else {
                                        Log.d(TAG, "discoverServices failed to start");
                                    }
                                };
                                BleDriver.mainHandler.postDelayed(discoverRunnable, delay);
                            } else if (bondState == BOND_BONDING) {
                                // Bonding process in progress, let it complete
                                Log.i(TAG, "waiting for bonding to complete");
                            }
                        } else if (newState == BluetoothProfile.STATE_DISCONNECTED) {
                            Log.d(TAG, "onConnectionStateChange(): disconnected");
                            BleInterface.BLEHandleLostPeer(getRemotePID());
                            if (getClientState() == CONNECTION_STATE.CONNECTED) {
                                close();
                            }
                            setClientState(CONNECTION_STATE.DISCONNECTED);
                        } else {
                            Log.e(TAG, "onConnectionStateChange(): unknown state");
                            close();
                        }
                    } else {
                        Log.e(TAG, "onConnectionStateChange(): status error=" + status);
                        close();
                    }
                }

                @Override
                public void onServicesDiscovered(BluetoothGatt gatt, int status) {
                    super.onServicesDiscovered(gatt, status);
                    Log.v(TAG, String.format("onServicesDiscovered for device %s", getMACAddress()));
                    mDiscoveryStarted = false;

                    if (status != GATT_SUCCESS) {
                        Log.e(TAG, String.format("service discovery failed due to internal error '%s', disconnecting", status));
                        disconnect();
                        return;
                    }

                    Log.i(TAG, String.format("discovered %d services for '%s'", gatt.getServices().size(), mBluetoothDevice.getAddress()));
                    handshake();
                }

                @Override
                public void onCharacteristicRead(BluetoothGatt gatt,
                                                 BluetoothGattCharacteristic characteristic,
                                                 int status) {
                    super.onCharacteristicRead(gatt, characteristic, status);
                    Log.v(TAG, String.format("onCharacteristicRead for device %s", getMACAddress()));

                    if (status != GATT_SUCCESS) {
                        Log.e(TAG, "onCharacteristicRead(): read error");
                        disconnect();
                        close();
                        BleQueue.completedCommand();
                        return ;
                    }
                    if (characteristic.getUuid().equals(GattServer.READER_UUID)) {
                        byte[] value = characteristic.getValue();
                        if (value.length == 0) {
                            Log.d(TAG, "onCharacteristicRead(): received data length is null");
                            BleQueue.completedCommand();
                            return ;
                        } else {
                            handleClientDataReceived(value);
                        }
                    } else {
                        Log.e(TAG, "onCharacteristicRead(): wrong read characteristic");
                        disconnect();
                        close();
                    }
                    BleQueue.completedCommand();
                }

                @Override
                public void onCharacteristicWrite (BluetoothGatt gatt,
                                                   BluetoothGattCharacteristic characteristic,
                                                   int status) {
                    super.onCharacteristicWrite(gatt, characteristic, status);
                    Log.v(TAG, String.format("onCharacteristicWrite for device %s", getMACAddress()));

                    handleClientDataSent();
                    BleQueue.completedCommand();
                }

                @Override
                public void onMtuChanged(BluetoothGatt gatt, int mtu, int status) {
                    super.onMtuChanged(gatt, mtu, status);
                    Log.v(TAG, String.format("onMtuChanged(): mtu %s for device %s", mtu, getMACAddress()));
                    PeerDevice peerDevice;

                    if (status != GATT_SUCCESS) {
                        Log.e(TAG, "onMtuChanged() error: transmission error");
                        BleQueue.completedCommand();
                        close();
                        return ;
                    }

                    mMtu = mtu;
                    BleQueue.completedCommand();
                }

                @Override
                public void onCharacteristicChanged (BluetoothGatt gatt, BluetoothGattCharacteristic characteristic) {
                    super.onCharacteristicChanged(gatt, characteristic);
                    Log.v(TAG, String.format("onCharacteristicChanged called for device %s", getMACAddress()));

                    byte[] value = characteristic.getValue();
                    Log.d(TAG, String.format("onCharacteristicChanged: value is %s", Base64.getEncoder().encodeToString(value)));
                    handleClientDataReceived(value);
                }
    };
}
