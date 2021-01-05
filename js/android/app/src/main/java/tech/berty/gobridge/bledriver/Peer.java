package tech.berty.gobridge.bledriver;

import android.util.Log;

public class Peer {
    private static final String TAG = "bty.ble.Peer";

    private String mPeerID;

    private PeerDevice mPeerDevice = null;

    public Peer(String peerID) {
        mPeerID = peerID;
    }

    public synchronized String getPeerID() {
        return mPeerID;
    }

    public synchronized void setPeerDevice(PeerDevice peerDevice) {
        mPeerDevice = peerDevice;
    }

    public synchronized PeerDevice getPeerDevice() {
        return mPeerDevice;
    }

    public synchronized boolean isReady() {
        if (mPeerDevice.isServerReady() && mPeerDevice.isClientReady()) {
            return true;
        }
        return false;
    }
}
