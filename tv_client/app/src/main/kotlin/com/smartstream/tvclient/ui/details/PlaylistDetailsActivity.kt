package com.smartstream.tvclient.ui.details

import android.os.Bundle
import androidx.fragment.app.FragmentActivity
import com.smartstream.tvclient.R

/**
 * Activity for displaying playlist details.
 * Shows playlist poster, name, type, and list of media items.
 */
class PlaylistDetailsActivity : FragmentActivity() {

    companion object {
        private const val TAG = "PlaylistDetailsAct"
        const val EXTRA_PLAYLIST_ID = "playlist_id"
        const val EXTRA_PLAYLIST_NAME = "playlist_name"
        const val EXTRA_PLAYLIST_POSTER = "playlist_poster"
        const val EXTRA_PLAYLIST_HAS_CHILDREN = "playlist_has_children"
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        android.util.Log.d(TAG, "onCreate: =============== PlaylistDetailsActivity started ===============")
        super.onCreate(savedInstanceState)

        android.util.Log.d(TAG, "onCreate: Extras - playlistId=${intent.getStringExtra(EXTRA_PLAYLIST_ID)}")
        android.util.Log.d(TAG, "onCreate: Extras - playlistName=${intent.getStringExtra(EXTRA_PLAYLIST_NAME)}")
        android.util.Log.d(TAG, "onCreate: Extras - playlistPoster=${intent.getStringExtra(EXTRA_PLAYLIST_POSTER)}")
        android.util.Log.d(TAG, "onCreate: Extras - playlistHasChildren=${intent.getBooleanExtra(EXTRA_PLAYLIST_HAS_CHILDREN, false)}")

        setContentView(R.layout.activity_playlist_details)
        android.util.Log.d(TAG, "onCreate: Layout set")

        if (savedInstanceState == null) {
            android.util.Log.d(TAG, "onCreate: Creating PlaylistDetailsFragment")
            val fragment = PlaylistDetailsFragment()
            supportFragmentManager.beginTransaction()
                .replace(R.id.details_fragment, fragment)
                .commit()
            android.util.Log.d(TAG, "onCreate: Fragment transaction committed")
        } else {
            android.util.Log.d(TAG, "onCreate: Fragment already exists (savedInstanceState not null)")
        }
    }
}
