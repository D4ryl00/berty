package tech.berty.gobridge.bledriver;

import android.bluetooth.BluetoothDevice;
import android.bluetooth.BluetoothGatt;
import android.bluetooth.BluetoothGattCallback;
import android.bluetooth.BluetoothGattCharacteristic;
import android.bluetooth.BluetoothGattService;
import android.bluetooth.BluetoothProfile;
import android.content.Context;
import android.util.Log;

import androidx.annotation.NonNull;

import java.util.List;

import static android.bluetooth.BluetoothDevice.BOND_BONDED;
import static android.bluetooth.BluetoothDevice.BOND_NONE;

public class PeerDevice {
    private static final String TAG = "bty.ble.PeerDevice";

    // Max MTU requested
    // See https://chromium.googlesource.com/aosp/platform/system/bt/+/29e794418452c8b35c2d42fe0cda81acd86bbf43/stack/include/gatt_api.h#123
    private static final int REQUEST_MTU = 517;

    public enum CONNECTION_STATE {
        DISCONNECTED,
        CONNECTED,
        CONNECTING,
        DISCONNECTING
    }
    private CONNECTION_STATE mState = CONNECTION_STATE.DISCONNECTED;

    private Context mContext;
    private BluetoothDevice mBluetoothDevice;
    private BluetoothGatt mBluetoothGatt;

    private final Object mLockState = new Object();
    private final Object mLockRemotePID = new Object();
    private final Object mLockMtu = new Object();
    private final Object mLockClient = new Object();
    private final Object mLockServer = new Object();
    private final Object mLockPeer = new Object();

    private BluetoothGattService mBertyService;
    private BluetoothGattCharacteristic mPeerIDCharacteristic;
    private BluetoothGattCharacteristic mWriterCharacteristic;

    private Peer mPeer;
    private String mRemotePID;
    private String mLocalPID;
    private boolean mClientReady = false;
    private boolean mServerReady = false;

    //private int mMtu = 0;
    // default MTU is 23
    private int mMtu = 23;

    public PeerDevice(@NonNull Context context, @NonNull BluetoothDevice bluetoothDevice, String localPID) {
        mContext = context;
        mBluetoothDevice = bluetoothDevice;
        mLocalPID = localPID;
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

        if (!isConnected()) {
            setBluetoothGatt(mBluetoothDevice.connectGatt(mContext, false,
                mGattCallback, BluetoothDevice.TRANSPORT_LE));
        }
    }

    public boolean isConnected() {
        return getState() == CONNECTION_STATE.CONNECTED;
    }

    public boolean isDisconnected() {
        return getState() == CONNECTION_STATE.DISCONNECTED;
    }

    // setters and getters are accessed by the DeviceManager thread et this thread so we need to
    // synchronize them.
    public void setState(CONNECTION_STATE state) {
        synchronized (mLockState) {
            mState = state;
        }
    }

    public CONNECTION_STATE getState() {
        synchronized (mLockState) {
            return mState;
        }
    }

    public void setBluetoothGatt(BluetoothGatt gatt) {
        synchronized (mLockState) {
            mBluetoothGatt = gatt;
        }
    }

    public BluetoothGatt getBluetoothGatt() {
        synchronized (mLockState) {
            return mBluetoothGatt;
        }
    }

    public void setPeerIDCharacteristic(BluetoothGattCharacteristic peerID) {
        mPeerIDCharacteristic = peerID;
    }

    public BluetoothGattCharacteristic getPeerIDCharacteristic() {
        return mPeerIDCharacteristic;
    }

    public void setWriterCharacteristic(BluetoothGattCharacteristic write) {
        mWriterCharacteristic = write;
    }

    public BluetoothGattCharacteristic getWriterCharacteristic() {
        return mWriterCharacteristic;
    }

    public boolean updateWriterValue(String value) {
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

    public void handleServerPIDSent() {
        Log.d(TAG, "handleServerPIDSent: device: " + getMACAddress());
        Peer peer;

        if (!isServerReady()) {
            setServerReady(true);
            if ((peer = getPeer()) != null) {
                peer.CallFoundPeer();
            }
        } else {
            Log.e(TAG, "handleServerPIDSent: PID already sent, device: " + getMACAddress());
            close();
        }
    }

    public void handleReceivedRemotePID(String peerID) {
        Log.d(TAG, "handleReceivedRemotePID: device: " + getMACAddress());
        Peer peer;


        // test if a PeerDevice already exist for this remote PID
        if (PeerManager.get(peerID) != null) {
            Log.i(TAG, "handleReceivedRemotePID(): a connection already exists, close connection, device: " + getMACAddress());
            close();
            return ;
        }

        setRemotePID(peerID);
        setClientReady(true);
        peer = PeerManager.register(peerID, this);
        setPeer(peer);
        peer.CallFoundPeer();
    }

    public boolean handleServerDataReceived(byte[] data) {
        boolean status = false;

        if (updateWriterValue(new String(data))) {
            status = true;
        }
        BleInterface.BLEReceiveFromPeer(getRemotePID(), data);
        return status;
    }

    public void handleWriteLocalPID() {
        Log.d(TAG, "handleWriteLocalPID called");
        setClientReady(true);
        PeerManager.register(mRemotePID, this);
    }

    public void setClientReady(boolean state) {
        synchronized (mLockClient) {
            mClientReady = state;
        }
    }

    public boolean isClientReady() {
        synchronized (mLockClient) {
            return mClientReady;
        }
    }

    public void setServerReady(boolean state) {
        synchronized (mLockServer) {
            mServerReady = state;
        }
    }

    public boolean isServerReady() {
        synchronized (mLockServer) {
            return mServerReady;
        }
    }

    public void close() {
        if (isConnected()) {
            getBluetoothGatt().close();
            setState(CONNECTION_STATE.DISCONNECTING);
            setClientReady(false);
            setServerReady(false);
            setPeer(null);
            PeerManager.unregister(mRemotePID);
        }
    }

    private boolean takeBertyService(List<BluetoothGattService> services) {
        Log.d(TAG, "takeBertyService: device: " + getMACAddress());
        for (BluetoothGattService service : services) {
            if (service.getUuid().equals(GattServer.SERVICE_UUID)) {
                Log.d(TAG, "takeBertyService(): found, device: " + getMACAddress());
                setBertyService(service);
                return true;
            }
        }
        Log.i(TAG, "takeBertyService(): not found, device: " + getMACAddress());
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
        Log.d(TAG, "takeBertyCharacteristic(): device: " + getMACAddress());
        int nbOfFoundCharacteristics = 0;
        List<BluetoothGattCharacteristic> characteristics = mBertyService.getCharacteristics();
        for (BluetoothGattCharacteristic characteristic : characteristics) {
            if (characteristic.getUuid().equals(GattServer.PEER_ID_UUID)) {
                Log.d(TAG, "takeBertyCharacteristic(): peerID characteristic found, device: " + getMACAddress());
                if (checkCharacteristicProperties(characteristic,
                        BluetoothGattCharacteristic.PROPERTY_READ)) {
                    setPeerIDCharacteristic(characteristic);
                    nbOfFoundCharacteristics++;
                }
            } else if (characteristic.getUuid().equals(GattServer.WRITER_UUID)) {
                Log.d(TAG, "takeBertyCharacteristic(): writer characteristic found, device: " + getMACAddress());
                if (checkCharacteristicProperties(characteristic,
                        BluetoothGattCharacteristic.PROPERTY_WRITE)) {
                    setWriterCharacteristic(characteristic);
                    nbOfFoundCharacteristics++;
                }
            }
            if (nbOfFoundCharacteristics == 2) {
                return true;
            }
        }
        return false;
    }

    // Asynchronous operation that will be reported to the BluetoothGattCallback#onCharacteristicRead
    // callback.
    public boolean requestRemotePID() {
        Log.v(TAG, "requestRemotePID(): device: " + getMACAddress());
        if (!mBluetoothGatt.readCharacteristic(getPeerIDCharacteristic())) {
            Log.e(TAG, "requestRemotePID() error, device: " + getMACAddress());
            return false;
        }
        return true;
    }

    public boolean writeLocalPID() {
        Log.v(TAG, "writeLocalPID() called");

        BluetoothGattCharacteristic writer = getWriterCharacteristic();
        writer.setValue(mLocalPID);
        if (!mBluetoothGatt.writeCharacteristic(writer)) {
            Log.e(TAG, "writeLocalPID() error");
            return false;
        }
        return true;
    }

    public void setMtu(int mtu) {
        synchronized (mLockMtu) {
            mMtu = mtu;
        }
    }

    public int getMtu() {
        synchronized (mLockMtu) {
            return mMtu;
        }
    }

    private BluetoothGattCallback mGattCallback =
            new BluetoothGattCallback() {
                @Override
                public void onConnectionStateChange(BluetoothGatt gatt, int status, int newState) {
                    super.onConnectionStateChange(gatt, status, newState);
                    Log.d(TAG, "onConnectionStateChange() called by device " + gatt.getDevice().getAddress());
                    if (newState == BluetoothProfile.STATE_CONNECTED) {
                        Log.d(TAG, "onConnectionStateChange(): connected");
                        setState(CONNECTION_STATE.CONNECTED);
                        gatt.requestMtu(REQUEST_MTU);
                    } else if (newState == BluetoothProfile.STATE_DISCONNECTED) {
                        Log.d(TAG, "onConnectionStateChange(): disconnected");
                        BleInterface.BLEHandleLostPeer(getRemotePID());
                        setState(CONNECTION_STATE.DISCONNECTED);
                        setBluetoothGatt(null);
                    } else {
                        Log.e(TAG, "onConnectionStateChange(): unknown state");
                        close();
                    }
                }

                @Override
                public void onServicesDiscovered(BluetoothGatt gatt, int status) {
                    super.onServicesDiscovered(gatt, status);
                    Log.d(TAG, "onServicesDiscovered: device: " + getMACAddress());
                    if (takeBertyService(gatt.getServices())) {
                        if (takeBertyCharacteristics()) {
                            requestRemotePID();
                        }
                    }
                }

                @Override
                public void onCharacteristicRead(BluetoothGatt gatt,
                                                 BluetoothGattCharacteristic characteristic,
                                                 int status) {
                    super.onCharacteristicRead(gatt, characteristic, status);
                    Log.d(TAG, "onCharacteristicRead: device: " + getMACAddress());
                    if (status != BluetoothGatt.GATT_SUCCESS) {
                        Log.e(TAG, "onCharacteristicRead(): read error");
                        close();
                    }
                    if (characteristic.getUuid().equals(GattServer.PEER_ID_UUID)) {
                        String peerID;
                        if ((peerID = characteristic.getStringValue(0)) == null
                                || peerID.length() == 0) {
                            Log.e(TAG, "onCharacteristicRead() error: peerID is null");
                            return ;
                        }
                        handleReceivedRemotePID(peerID);
                    } else {
                        Log.e(TAG, "onCharacteristicRead(): wrong read characteristic");
                        close();
                    }
                }

                @Override
                public void onCharacteristicWrite (BluetoothGatt gatt,
                                                   BluetoothGattCharacteristic characteristic,
                                                   int status) {
                    super.onCharacteristicWrite(gatt, characteristic, status);
                }

                @Override
                public void onMtuChanged(BluetoothGatt gatt, int mtu, int status) {
                    Log.d(TAG, "onMtuChanged(): mtu:" + mtu + ", device: " + getMACAddress());
                    PeerDevice peerDevice;
                    if (status != BluetoothGatt.GATT_SUCCESS) {
                        Log.e(TAG, "onMtuChanged() error: transmission error");
                        close();
                        return ;
                    }
                    if ((peerDevice = DeviceManager.get(gatt.getDevice().getAddress())) == null) {
                        Log.e(TAG, "onMtuChanged() error: device not found");
                        gatt.close();
                        return ;
                    }
                    peerDevice.setMtu(mtu);
                    gatt.discoverServices();
                }
            };
}
