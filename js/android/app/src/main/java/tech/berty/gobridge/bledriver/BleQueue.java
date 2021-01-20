package tech.berty.gobridge.bledriver;

import android.util.Log;

import java.util.Queue;
import java.util.concurrent.ConcurrentLinkedQueue;

public class BleQueue {
    private static final String TAG = "bty.ble.BleQueue";

    // Maximum number of retries of commands
    private static final int MAX_TRIES = 3;

    private static final Queue<Runnable> mCommandQueue = new ConcurrentLinkedQueue<>();

    private static boolean mCommandQueueBusy = false;
    private static boolean mIsRetrying = false;
    private static int mNrTries = 0;


    public synchronized static boolean add(Runnable task) {
        return mCommandQueue.add(task);
    }

    /**
     * The current command has been completed, move to the next command in the queue (if any)
     */
    public synchronized static void completedCommand() {
        Log.v(TAG, "completedCommand called");

        mIsRetrying = false;
        mCommandQueue.poll();
        mCommandQueueBusy = false;
        nextCommand();
    }

    /**
     * Retry the current command. Typically used when a read/write fails and triggers a bonding procedure
     */
    public synchronized static void retryCommand() {
        Log.v(TAG, "retryCommand called");

        mCommandQueueBusy = false;
        Runnable currentCommand = mCommandQueue.peek();
        if (currentCommand != null) {
            if (mNrTries >= MAX_TRIES) {
                // Max retries reached, give up on this one and proceed
                Log.d(TAG, "max number of tries reached, not retrying operation anymore");
                mCommandQueue.poll();
            } else {
                mIsRetrying = true;
            }
        }
        nextCommand();
    }

    /**
     * Execute the next command in the subscribe queue.
     * A queue is used because the calls have to be executed sequentially.
     * If the read or write fails, the next command in the queue is executed.
     */
    public synchronized static void nextCommand() {
        Log.v(TAG, "nextCommand called");

        // If there is still a command being executed, then bail out
        if (mCommandQueueBusy) {
            Log.d(TAG, "nextCommand: another command is running, cancel");
            return;
        }

        // Check if there is something to do at all
        final Runnable bluetoothCommand = mCommandQueue.peek();
        if (bluetoothCommand == null) return;

        // Execute the next command in the queue
        mCommandQueueBusy = true;
        if (!mIsRetrying) {
            mNrTries = 0;
        }
        BleDriver.mainHandler.post(new Runnable() {
            @Override
            public void run() {
                try {
                    bluetoothCommand.run();
                } catch (Exception e) {
                    Log.e(TAG, "nextCommand: command exception", e);
                    completedCommand();
                }
            }
        });
    }

    public synchronized static void clear() {
        mCommandQueue.clear();
        mCommandQueueBusy = false;
    }
}
