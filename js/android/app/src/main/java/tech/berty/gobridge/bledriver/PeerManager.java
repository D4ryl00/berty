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

    public static synchronized void register(String key, PeerDevice peerDevice) {
        Log.d(TAG, "register() called");
        Peer peer;

        if ((peer = mPeers.get(key)) == null) {
            Log.d(TAG, "register(): peer unknown");
            peer = new Peer(key);
            peer.setPeerDevice(peerDevice);
            mPeers.put(key, peer);
        } else {
            Log.d(TAG, "register(): peer known");
            if (peer.isReady()) {
                Log.d(TAG, "register(): peer ready");
                BleInterface.BLEHandleFoundPeer(key);
            }
        }
    }

    public static synchronized Peer get(String key) {
        return mPeers.get(key);
    }
}
