package com.smartstream.tvclient.ui.details

import android.content.Intent
import android.graphics.Bitmap
import android.graphics.drawable.Drawable
import android.os.Bundle
import androidx.core.content.ContextCompat
import androidx.leanback.app.DetailsSupportFragment
import androidx.leanback.app.DetailsSupportFragmentBackgroundController
import androidx.leanback.widget.Action
import androidx.leanback.widget.ArrayObjectAdapter
import androidx.leanback.widget.ClassPresenterSelector
import androidx.leanback.widget.DetailsOverviewRow
import androidx.leanback.widget.FullWidthDetailsOverviewRowPresenter
import androidx.leanback.widget.FullWidthDetailsOverviewSharedElementHelper
import androidx.leanback.widget.ImageCardView
import com.bumptech.glide.Glide
import com.bumptech.glide.request.target.CustomTarget
import com.bumptech.glide.request.transition.Transition
import com.smartstream.tvclient.R
import com.smartstream.tvclient.ui.player.PlayerActivity
import com.smartstream.tvclient.utils.SharedPrefsManager

/**
 * Fragment displaying media details.
 * Uses Leanback DetailsSupportFragment.
 */
class MediaDetailsFragment : DetailsSupportFragment() {

    private lateinit var backgroundController: DetailsSupportFragmentBackgroundController
    private lateinit var rowsAdapter: ArrayObjectAdapter
    private lateinit var presenterSelector: ClassPresenterSelector

    private var mediaId: String? = null
    private var mediaName: String? = null
    private var mediaPoster: String? = null
    private var mediaDuration: Int = 0
    private var mediaSize: Long = 0L

    companion object {
        private const val ACTION_PLAY = 1L
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        android.util.Log.d("MediaDetailsFragment", "onCreate: =============== MediaDetailsFragment created ===============")
        super.onCreate(savedInstanceState)

        // Get extras from activity
        activity?.intent?.let { intent ->
            mediaId = intent.getStringExtra(MediaDetailsActivity.EXTRA_MEDIA_ID)
            mediaName = intent.getStringExtra(MediaDetailsActivity.EXTRA_MEDIA_NAME)
            mediaPoster = intent.getStringExtra(MediaDetailsActivity.EXTRA_MEDIA_POSTER)
            mediaDuration = intent.getIntExtra(MediaDetailsActivity.EXTRA_MEDIA_DURATION, 0)
            mediaSize = intent.getLongExtra(MediaDetailsActivity.EXTRA_MEDIA_SIZE, 0L)

            android.util.Log.d("MediaDetailsFragment", "onCreate: Got extras - mediaId=$mediaId")
            android.util.Log.d("MediaDetailsFragment", "onCreate: Got extras - mediaName=$mediaName")
            android.util.Log.d("MediaDetailsFragment", "onCreate: Got extras - mediaPoster=$mediaPoster")
            android.util.Log.d("MediaDetailsFragment", "onCreate: Got extras - mediaDuration=$mediaDuration")
            android.util.Log.d("MediaDetailsFragment", "onCreate: Got extras - mediaSize=$mediaSize")
        }

        android.util.Log.d("MediaDetailsFragment", "onCreate: Setting up background controller")
        backgroundController = DetailsSupportFragmentBackgroundController(this)
        backgroundController.enableParallax()
    }

    override fun onActivityCreated(savedInstanceState: Bundle?) {
        super.onActivityCreated(savedInstanceState)
        android.util.Log.d("MediaDetailsFragment", "onActivityCreated: Started")

        android.util.Log.d("MediaDetailsFragment", "onActivityCreated: Setting up adapter")
        setupAdapter()

        android.util.Log.d("MediaDetailsFragment", "onActivityCreated: Setting up details overview row")
        setupDetailsOverviewRow()

        android.util.Log.d("MediaDetailsFragment", "onActivityCreated: Loading poster image")
        loadPosterImage()

        android.util.Log.d("MediaDetailsFragment", "onActivityCreated: Fragment initialization complete")
    }

    private fun setupAdapter() {
        // Setup presenter selector
        val detailsPresenter = FullWidthDetailsOverviewRowPresenter(
            MediaDetailsDescriptionPresenter()
        ).apply {
            backgroundColor = ContextCompat.getColor(requireContext(), R.color.tv_card_background)
            initialState = FullWidthDetailsOverviewRowPresenter.STATE_HALF

            val helper = FullWidthDetailsOverviewSharedElementHelper()
            helper.setSharedElementEnterTransition(activity, "media_details")
            setListener(helper)

            // Setup action click listener on presenter
            setOnActionClickedListener { action ->
                android.util.Log.d("MediaDetailsFragment", "setOnActionClickedListener: Action clicked! id=${action.id}, label=${action.label1}")
                when (action.id) {
                    ACTION_PLAY -> {
                        android.util.Log.d("MediaDetailsFragment", "setOnActionClickedListener: ACTION_PLAY matched, calling playMedia()")
                        playMedia()
                    }
                    else -> {
                        android.util.Log.w("MediaDetailsFragment", "setOnActionClickedListener: Unknown action id=${action.id}")
                    }
                }
            }
        }

        presenterSelector = ClassPresenterSelector().apply {
            addClassPresenter(DetailsOverviewRow::class.java, detailsPresenter)
        }

        rowsAdapter = ArrayObjectAdapter(presenterSelector)
        adapter = rowsAdapter
    }

    private fun setupDetailsOverviewRow() {
        android.util.Log.d("MediaDetailsFragment", "setupDetailsOverviewRow: Building description")
        val description = buildDescription()
        android.util.Log.d("MediaDetailsFragment", "setupDetailsOverviewRow: Description = $description")

        val row = DetailsOverviewRow(description)
        android.util.Log.d("MediaDetailsFragment", "setupDetailsOverviewRow: DetailsOverviewRow created")

        row.setImageDrawable(
            ContextCompat.getDrawable(
                requireContext(),
                R.drawable.ic_media_placeholder
            )
        )
        android.util.Log.d("MediaDetailsFragment", "setupDetailsOverviewRow: Placeholder image set")

        // Setup actions
        val actionsAdapter = ArrayObjectAdapter()
        val playAction = Action(
            ACTION_PLAY,
            getString(R.string.action_play),
            null
        )
        android.util.Log.d("MediaDetailsFragment", "setupDetailsOverviewRow: Created Play action with id=$ACTION_PLAY, label=${getString(R.string.action_play)}")
        actionsAdapter.add(playAction)
        row.actionsAdapter = actionsAdapter
        android.util.Log.d("MediaDetailsFragment", "setupDetailsOverviewRow: Actions adapter set with ${actionsAdapter.size()} actions")

        rowsAdapter.add(row)
        android.util.Log.d("MediaDetailsFragment", "setupDetailsOverviewRow: Row added to adapter, total rows=${rowsAdapter.size()}")
    }

    private fun buildDescription(): String {
        val parts = mutableListOf<String>()

        // Add duration
        if (mediaDuration > 0) {
            val hours = mediaDuration / 3600
            val minutes = (mediaDuration % 3600) / 60
            val seconds = mediaDuration % 60

            val durationStr = if (hours > 0) {
                String.format("%02d:%02d:%02d", hours, minutes, seconds)
            } else {
                String.format("%02d:%02d", minutes, seconds)
            }
            parts.add(getString(R.string.media_duration, durationStr))
        }

        // Add size
        if (mediaSize > 0) {
            val sizeStr = when {
                mediaSize < 1024 -> "$mediaSize B"
                mediaSize < 1024 * 1024 -> String.format("%.2f KB", mediaSize / 1024.0)
                mediaSize < 1024 * 1024 * 1024 -> String.format("%.2f MB", mediaSize / (1024.0 * 1024.0))
                else -> String.format("%.2f GB", mediaSize / (1024.0 * 1024.0 * 1024.0))
            }
            parts.add(getString(R.string.media_size, sizeStr))
        }

        return parts.joinToString("\n")
    }

    private fun loadPosterImage() {
        val row = rowsAdapter.get(0) as DetailsOverviewRow

        val posterUrl = mediaPoster?.let {
            if (it.startsWith("http")) it else "${SharedPrefsManager.getBaseUrl()}${it.removePrefix("/")}"
        }

        if (posterUrl != null) {
            Glide.with(requireContext())
                .asBitmap()
                .load(posterUrl)
                .error(R.drawable.ic_media_placeholder)
                .into(object : CustomTarget<Bitmap>() {
                    override fun onResourceReady(
                        resource: Bitmap,
                        transition: Transition<in Bitmap>?
                    ) {
                        row.setImageBitmap(requireContext(), resource)
                        rowsAdapter.notifyArrayItemRangeChanged(0, rowsAdapter.size())
                        backgroundController.coverBitmap = resource
                    }

                    override fun onLoadCleared(placeholder: Drawable?) {
                        // No-op
                    }
                })
        }
    }

    private fun playMedia() {
        android.util.Log.d("MediaDetailsFragment", "playMedia: called")
        android.util.Log.d("MediaDetailsFragment", "playMedia: mediaId=$mediaId")
        android.util.Log.d("MediaDetailsFragment", "playMedia: mediaName=$mediaName")

        mediaId?.let { id ->
            android.util.Log.d("MediaDetailsFragment", "playMedia: creating Intent for PlayerActivity")
            val intent = Intent(requireContext(), PlayerActivity::class.java).apply {
                putExtra(PlayerActivity.EXTRA_MEDIA_ID, id)
                putExtra(PlayerActivity.EXTRA_MEDIA_NAME, mediaName)
            }
            android.util.Log.d("MediaDetailsFragment", "playMedia: starting PlayerActivity")
            android.util.Log.d("MediaDetailsFragment", "playMedia: Intent extras - mediaId=$id, mediaName=$mediaName")

            try {
                startActivity(intent)
                android.util.Log.d("MediaDetailsFragment", "playMedia: PlayerActivity started successfully")
            } catch (e: Exception) {
                android.util.Log.e("MediaDetailsFragment", "playMedia: Failed to start PlayerActivity", e)
                android.widget.Toast.makeText(requireContext(), "Failed to start player: ${e.message}", android.widget.Toast.LENGTH_LONG).show()
            }
        } ?: run {
            android.util.Log.e("MediaDetailsFragment", "playMedia: mediaId is null!")
            android.widget.Toast.makeText(requireContext(), "Error: Media ID is missing", android.widget.Toast.LENGTH_LONG).show()
        }
    }
}
