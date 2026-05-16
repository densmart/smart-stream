package com.smartstream.tvclient.ui.details

import android.os.Bundle
import android.util.Log
import androidx.fragment.app.FragmentActivity
import com.smartstream.tvclient.R

/**
 * Activity for displaying media details.
 * Shows poster, title, description, duration, size, and "Play" button.
 */
class MediaDetailsActivity : FragmentActivity() {

    companion object {
        private const val TAG = "MediaDetailsActivity"
        const val EXTRA_MEDIA_ID = "media_id"
        const val EXTRA_MEDIA_NAME = "media_name"
        const val EXTRA_MEDIA_POSTER = "media_poster"
        const val EXTRA_MEDIA_DURATION = "media_duration"
        const val EXTRA_MEDIA_SIZE = "media_size"
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        Log.d(TAG, "onCreate: =============== MediaDetailsActivity started ===============")
        super.onCreate(savedInstanceState)

        Log.d(TAG, "onCreate: Extras - mediaId=${intent.getStringExtra(EXTRA_MEDIA_ID)}")
        Log.d(TAG, "onCreate: Extras - mediaName=${intent.getStringExtra(EXTRA_MEDIA_NAME)}")
        Log.d(TAG, "onCreate: Extras - mediaPoster=${intent.getStringExtra(EXTRA_MEDIA_POSTER)}")
        Log.d(TAG, "onCreate: Extras - mediaDuration=${intent.getIntExtra(EXTRA_MEDIA_DURATION, 0)}")
        Log.d(TAG, "onCreate: Extras - mediaSize=${intent.getLongExtra(EXTRA_MEDIA_SIZE, 0L)}")

        setContentView(R.layout.activity_media_details)
        Log.d(TAG, "onCreate: Layout set")

        if (savedInstanceState == null) {
            Log.d(TAG, "onCreate: Creating MediaDetailsFragmentSimple")
            val fragment = MediaDetailsFragmentSimple()
            supportFragmentManager.beginTransaction()
                .replace(R.id.details_fragment, fragment)
                .commit()
            Log.d(TAG, "onCreate: Fragment transaction committed")
        } else {
            Log.d(TAG, "onCreate: Fragment already exists (savedInstanceState not null)")
        }
    }
}
