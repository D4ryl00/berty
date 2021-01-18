package tech.berty.gobridge.bledriver;

import android.bluetooth.BluetoothAdapter;
import android.bluetooth.BluetoothGattCharacteristic;
import android.bluetooth.BluetoothGattServer;
import android.bluetooth.BluetoothGattService;
import android.bluetooth.BluetoothManager;
import android.content.Context;
import android.os.ParcelUuid;
import android.util.Log;

import java.util.UUID;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

import static android.bluetooth.BluetoothGattCharacteristic.PERMISSION_READ;
import static android.bluetooth.BluetoothGattCharacteristic.PERMISSION_WRITE;
import static android.bluetooth.BluetoothGattCharacteristic.PROPERTY_READ;
import static android.bluetooth.BluetoothGattCharacteristic.PROPERTY_WRITE;
import static android.bluetooth.BluetoothGattService.SERVICE_TYPE_PRIMARY;
import static android.content.Context.BLUETOOTH_SERVICE;

public class GattServer {
    private final String TAG = "bty.ble.GattServer";

    // GATT service UUID
    static final UUID SERVICE_UUID = UUID.fromString("A06C6AB8-886F-4D56-82FC-2CF8610D668D");
    static final UUID READER_UUID = UUID.fromString("0EF50D30-E208-4315-B323-D05E0A23E6B5");
    static final UUID WRITER_UUID = UUID.fromString("000CBD77-8D30-4EFF-9ADD-AC5F10C2CC1B");
    static final ParcelUuid P_SERVICE_UUID = new ParcelUuid(SERVICE_UUID);

    // GATT service objects
    private BluetoothGattService mService;
    private BluetoothGattCharacteristic mReaderCharacteristic;
    private BluetoothGattCharacteristic mWriterCharacteristic;

    private Context mContext;
    private BluetoothManager mBluetoothManager;
    private CountDownLatch mDoneSignal;
    private GattServerCallback mGattServerCallback;
    private BluetoothGattServer mBluetoothGattServer;
    private volatile boolean mInit = false;
    private volatile boolean mStarted = false;

    private Lock mLock = new ReentrantLock();

    public GattServer(Context context, BluetoothManager bluetoothManager) {
        mContext = context;
        mBluetoothManager = bluetoothManager;
        initGattService();
    }

    private void initGattService() {
        Log.i(TAG, "initGattService called");

        mService = new BluetoothGattService(SERVICE_UUID, SERVICE_TYPE_PRIMARY);
        mReaderCharacteristic = new BluetoothGattCharacteristic(READER_UUID, PROPERTY_READ, PERMISSION_READ);
        mWriterCharacteristic = new BluetoothGattCharacteristic(WRITER_UUID, PROPERTY_WRITE, PERMISSION_WRITE);

        if (!mService.addCharacteristic(mReaderCharacteristic) || !mService.addCharacteristic(mWriterCharacteristic)) {
            Log.e(TAG, "setupService failed: can't add characteristics to service");
            return ;
        }

        mDoneSignal = new CountDownLatch(1);
        mGattServerCallback = new GattServerCallback(mContext, this, mDoneSignal);

        mInit = true;
    }

    // After adding a new service, the success of this operation will be given to the callback
    // BluetoothGattServerCallback#onServiceAdded. It's only after this callback that the server
    // will be ready.
    public boolean start(final String peerID) {
        Log.i(TAG, "start called");

        if (!mInit) {
            Log.e(TAG, "start: GATT service not init");
            return false;
        }
        if (isStarted()) {
            Log.i(TAG, "start: GATT service already started");
            return true;
        }

        mGattServerCallback.setLocalPID(peerID);

        mBluetoothGattServer = mBluetoothManager.openGattServer(mContext, mGattServerCallback);

        mReaderCharacteristic.setValue(peerID);
        mWriterCharacteristic.setValue("");
        if (!mBluetoothGattServer.addService(mService)) {
            Log.e(TAG, "setupGattServer error: cannot add a new service");
            mBluetoothGattServer.close();
            mBluetoothGattServer = null;
            return false;
        }

        // wait that service starts
        try {
           mDoneSignal.await();
        } catch (InterruptedException e) {
            Log.e(TAG, "start: interrupted exception:", e);
        }

        // mStarted is updated by GattServerCallback
        return isStarted();
    }

    public BluetoothGattServer getGattServer() {
        BluetoothGattServer gattServer;
        mLock.lock();
        try {
            gattServer = mBluetoothGattServer;
        } finally {
            mLock.unlock();
        }
        return gattServer;
    }

    public void setStarted(boolean started) {
        mLock.lock();
        try {
            mStarted = started;
        } finally {
            mLock.unlock();
        }
    }

    public boolean isStarted() {
        boolean started;
        mLock.lock();
        try {
            started = mStarted;
        } finally {
            mLock.unlock();
        }
        return started;
    }

    public void stop() {
        Log.i(TAG, "stop() called");
        if (isStarted()) {
            setStarted(false);
            mBluetoothGattServer.close();
            mLock.lock();
            try {
                mBluetoothGattServer = null;
            } finally {
                mLock.unlock();
            }
        }
    }
}
