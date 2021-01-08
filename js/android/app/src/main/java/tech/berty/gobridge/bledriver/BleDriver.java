package tech.berty.gobridge.bledriver;

import android.annotation.SuppressLint;
import android.app.Application;
import android.bluetooth.BluetoothAdapter;
import android.bluetooth.BluetoothGatt;
import android.bluetooth.BluetoothGattCharacteristic;
import android.bluetooth.BluetoothManager;
import android.bluetooth.le.AdvertiseData;
import android.bluetooth.le.AdvertiseSettings;
import android.bluetooth.le.BluetoothLeAdvertiser;
import android.bluetooth.le.BluetoothLeScanner;
import android.bluetooth.le.ScanFilter;
import android.bluetooth.le.ScanSettings;
import android.content.Context;
import android.content.IntentFilter;
import android.content.pm.PackageManager;
import android.util.Log;

import java.util.Base64;
import java.util.Collections;
import java.util.concurrent.CountDownLatch;

// Make the BleDriver class a Singleton
// see https://medium.com/@kevalpatel2106/how-to-make-the-perfect-singleton-de6b951dfdb0
public class BleDriver {
    private static final String TAG = "bty.ble.BleDriver";

    private static volatile BleDriver mBleDriver;

    static final String ACTION_PEER_FOUND = "BleDriver.ACTION_PEER_FOUND";

    private static Context mAppContext;
    private static BluetoothManager mBluetoothManager;
    private static BluetoothAdapter mBluetoothAdapter;
    private static GattServer mGattServer;

    private static Advertiser mAdvertiser;
    private static Scanner mScanner;

    private static boolean mInit = false;
    private static boolean mStarted = false;

    private BleDriver(Context context) {
        if (mBleDriver != null) {
            throw new RuntimeException("Use getInstance() method to get the singleton instance of this class");
        }
        mAppContext = context;
    }

    // Singleton method
    public static BleDriver getInstance(Context appContext) {
        if (mBleDriver == null) {
            synchronized (BleDriver.class) {
                if (mBleDriver == null) {
                    mBleDriver = new BleDriver(appContext);
                    initDriver();
                }
            }
        }
        return mBleDriver;
    }

    // Init BluetoothAdapter object and test if bluetooth is enabled.
    private static synchronized boolean initSystemBle() {
        Log.i(TAG, "initBluetoothAdapter(): init system bluetooth requirements");
        if (!mAppContext.getPackageManager().hasSystemFeature(PackageManager.FEATURE_BLUETOOTH_LE)) {
            Log.e(TAG, "initSystemBle: BLE is not supported by this device");
            return false;
        }
        // Initializes Bluetooth adapter.
        mBluetoothManager = (BluetoothManager)mAppContext.getSystemService(Context.BLUETOOTH_SERVICE);
        if ((mBluetoothAdapter = mBluetoothManager.getAdapter()) == null) {
            Log.e(TAG, "initBluetoothAdapter(): bluetooth adapter not available");
            return false;
        }
        if (!mBluetoothAdapter.isEnabled()) {
            Log.e(TAG, "initBluetoothAdapter(): bluetooth not enabled");
            return false;
        }
        Log.i(TAG, "initBluetoothAdapter(): bluetooth is supported on this hardware platform and enabled");
        return true;
    }

    // main initialization method
    private static synchronized void initDriver() {
        if (!initSystemBle()) {
            Log.e(TAG, "initDriver: initSystemBle failed");
            return;
        }

        mAdvertiser = new Advertiser(mBluetoothAdapter);
        if (!mAdvertiser.isInit()) {
            Log.e(TAG, "initDriver: Advertiser init failed");
        }

        mScanner = new Scanner(mAppContext, mBluetoothAdapter);
        if (!mScanner.isInit()) {
            Log.e(TAG, "initDriver: Scanner init failed");
        }

        // Setup context dependant objects
        PeerManager.setContext(mAppContext); // TODO: remove PeerManager
        mGattServer = new GattServer(mAppContext, mBluetoothManager);
        setInit(true);
    }

    public static synchronized void setInit(boolean status) {
        mInit = status;
    }

    public synchronized boolean isInit() {
        return mInit;
    }

    public synchronized boolean isStarted() {
        return mStarted;
    }

    public synchronized void setStarted(boolean status) {
        mStarted = status;
    }

    public synchronized void StartBleDriver(String localPeerID) {
        Log.d(TAG, "StartBleDriver() called");

        if (!isInit()) {
            Log.e(TAG, "StartBleDriver: driver not init");
            return ;
        }

        if (isStarted()) {
            Log.i(TAG, "StartBleDriver(): BLE driver is already on, one instance is allow");
            return ;
        }

        if (!mGattServer.start(localPeerID)) {
           return ;
        }
        setStarted(true);

        if (!mAdvertiser.start()) {
            Log.e(TAG, "StartBleDriver: failed to start advertising");
            StopBleDriver();
            return ;
        }

        if (!mScanner.start(localPeerID)) {
            Log.e(TAG, "StartBleDriver: failed to start scanning");
            StopBleDriver();
            return ;
        }

        Log.i(TAG, "StartBleDriver: initDriver completed");
    }

    public synchronized void StopBleDriver() {
        if (!isInit()) {
            Log.e(TAG, "StopBleDriver: driver not init");
            return ;
        }

        if (!isStarted()) {
            Log.d(TAG, "driver is not started");
            return ;
        }

        mAdvertiser.stop();
        mScanner.stop();
        DeviceManager.closeAllDeviceConnections();
        mGattServer.stop();
        setStarted(false);
    }

    public boolean SendToPeer(String remotePID, byte[] payload) {
        Log.d(TAG, "SendToPeer(): remotePID=" + remotePID + " payload=" + Base64.getEncoder().encodeToString(payload));
        PeerDevice peerDevice;
        BluetoothGattCharacteristic writer;
        BluetoothGatt gatt;

        if ((peerDevice = PeerManager.get(remotePID).getPeerDevice()) == null) {
            Log.e(TAG, "SendToPeer: remote device not found");
            return false;
        }
        if (peerDevice.isDisconnected()) {
            Log.e(TAG, "SendToPeer: remote device is disconnected");
            return false;
        }
        writer = peerDevice.getWriterCharacteristic();
        writer.setValue(payload);
        gatt = peerDevice.getBluetoothGatt();
        gatt.writeCharacteristic(writer);
        return true;
    }
}
