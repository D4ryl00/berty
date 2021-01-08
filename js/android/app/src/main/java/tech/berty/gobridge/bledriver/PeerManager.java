package tech.berty.gobridge.bledriver;

import android.content.Context;
import android.util.Log;

import java.util.HashMap;

public class PeerManager {
    private static final String TAG = "bty.ble.PeerManager";

    private static HashMap<String, Peer> mPeers = new HashMap<>();
    private static Context mContext;

    public static void setContext(Context context) {
        mContext = context;
    }

    public static synchronized Peer register(String peerID, PeerDevice peerDevice) {
        Log.d(TAG, "register() called");
        Peer peer;

        if ((peer = mPeers.get(peerID)) == null) {
            Log.d(TAG, "register(): peer unknown");
            peer = new Peer(peerID);
            peer.setPeerDevice(peerDevice);
            mPeers.put(peerID, peer);
        } else {
            Log.d(TAG, "register(): peer already known");
        }
        return peer;
    }

    public static synchronized void unregister(String peerID) {
        Log.d(TAG, "unregister() called");
        Peer peer;

        if ((peer = mPeers.get(peerID)) != null) {
            Log.d(TAG, "unregister: peer removed");
            mPeers.remove(peerID);
        }
    }

    public static synchronized Peer get(String peerID) {
        return mPeers.get(peerID);
    }
}
