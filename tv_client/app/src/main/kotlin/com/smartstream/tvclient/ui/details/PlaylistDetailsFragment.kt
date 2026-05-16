package com.smartstream.tvclient.ui.details

import android.content.Intent
import android.graphics.Bitmap
import android.graphics.drawable.Drawable
import android.os.Bundle
import android.widget.Toast
import androidx.core.content.ContextCompat
import androidx.leanback.app.DetailsSupportFragment
import androidx.leanback.app.DetailsSupportFragmentBackgroundController
import androidx.leanback.widget.ArrayObjectAdapter
import androidx.leanback.widget.ClassPresenterSelector
import androidx.leanback.widget.DetailsOverviewRow
import androidx.leanback.widget.FullWidthDetailsOverviewRowPresenter
import androidx.leanback.widget.FullWidthDetailsOverviewSharedElementHelper
import androidx.leanback.widget.HeaderItem
import androidx.leanback.widget.ListRow
import androidx.leanback.widget.ListRowPresenter
import androidx.leanback.widget.OnItemViewClickedListener
import androidx.lifecycle.lifecycleScope
import com.bumptech.glide.Glide
import com.bumptech.glide.request.target.CustomTarget
import com.bumptech.glide.request.transition.Transition
import com.smartstream.tvclient.R
import com.smartstream.tvclient.data.model.Media
import com.smartstream.tvclient.data.model.Playlist
import com.smartstream.tvclient.data.repository.MediaRepository
import com.smartstream.tvclient.data.repository.PlaylistRepository
import com.smartstream.tvclient.ui.main.CardPresenter
import com.smartstream.tvclient.ui.player.PlayerActivity
import com.smartstream.tvclient.utils.SharedPrefsManager
import kotlinx.coroutines.launch

/**
 * Fragment displaying playlist details and its media items.
 * Uses Leanback DetailsSupportFragment.
 */
class PlaylistDetailsFragment : DetailsSupportFragment() {

    private lateinit var backgroundController: DetailsSupportFragmentBackgroundController
    private lateinit var rowsAdapter: ArrayObjectAdapter
    private lateinit var presenterSelector: ClassPresenterSelector

    private val mediaRepository = MediaRepository()
    private val playlistRepository = PlaylistRepository()

    private var playlistId: String? = null
    private var playlistName: String? = null
    private var playlistPoster: String? = null
    private var playlistHasChildren: Boolean = false
    private var posterBitmap: Bitmap? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        android.util.Log.d(TAG, "onCreate: =============== PlaylistDetailsFragment created ===============")
        super.onCreate(savedInstanceState)

        // Get extras from activity
        activity?.intent?.let { intent ->
            playlistId = intent.getStringExtra(PlaylistDetailsActivity.EXTRA_PLAYLIST_ID)
            playlistName = intent.getStringExtra(PlaylistDetailsActivity.EXTRA_PLAYLIST_NAME)
            playlistPoster = intent.getStringExtra(PlaylistDetailsActivity.EXTRA_PLAYLIST_POSTER)
            playlistHasChildren = intent.getBooleanExtra(PlaylistDetailsActivity.EXTRA_PLAYLIST_HAS_CHILDREN, false)
        }

        android.util.Log.d(TAG, "onCreate: playlistId=$playlistId")
        android.util.Log.d(TAG, "onCreate: playlistName=$playlistName")
        android.util.Log.d(TAG, "onCreate: playlistPoster=$playlistPoster")
        android.util.Log.d(TAG, "onCreate: playlistHasChildren=$playlistHasChildren")

        backgroundController = DetailsSupportFragmentBackgroundController(this)
        backgroundController.enableParallax()
    }

    override fun onActivityCreated(savedInstanceState: Bundle?) {
        super.onActivityCreated(savedInstanceState)
        android.util.Log.d(TAG, "onActivityCreated: Started")

        android.util.Log.d(TAG, "onActivityCreated: Setting up adapter")
        setupAdapter()

        android.util.Log.d(TAG, "onActivityCreated: Setting up details overview row")
        setupDetailsOverviewRow()

        android.util.Log.d(TAG, "onActivityCreated: Loading poster image")
        loadPosterImage()

        android.util.Log.d(TAG, "onActivityCreated: Loading playlist content (hasChildren=$playlistHasChildren)")
        loadPlaylistContent()
    }

    companion object {
        private const val TAG = "PlaylistDetailsFrag"
    }

    private fun setupAdapter() {
        // Setup presenter selector
        val detailsPresenter = FullWidthDetailsOverviewRowPresenter(
            MediaDetailsDescriptionPresenter()
        ).apply {
            backgroundColor = ContextCompat.getColor(requireContext(), R.color.tv_card_background)
            initialState = FullWidthDetailsOverviewRowPresenter.STATE_HALF

            val helper = FullWidthDetailsOverviewSharedElementHelper()
            helper.setSharedElementEnterTransition(activity, "playlist_details")
            setListener(helper)
        }

        presenterSelector = ClassPresenterSelector().apply {
            addClassPresenter(DetailsOverviewRow::class.java, detailsPresenter)
            addClassPresenter(ListRow::class.java, ListRowPresenter())
        }

        rowsAdapter = ArrayObjectAdapter(presenterSelector)
        adapter = rowsAdapter

        // Setup item click listener
        setOnItemViewClickedListener(OnItemViewClickedListener { _, item, _, _ ->
            when (item) {
                is Media -> openMediaPlayer(item)
                is Playlist -> openPlaylistDetails(item)
            }
        })
    }

    private fun setupDetailsOverviewRow() {
        val row = DetailsOverviewRow(getString(R.string.playlists_loading))

        row.setImageDrawable(
            ContextCompat.getDrawable(
                requireContext(),
                R.drawable.ic_playlist_placeholder
            )
        )

        rowsAdapter.add(row)
    }

    private fun loadPosterImage() {
        val row = rowsAdapter.get(0) as DetailsOverviewRow

        val posterUrl = playlistPoster?.let {
            when {
                it.startsWith("http") -> it
                it.startsWith("/static/posters/") -> "${SharedPrefsManager.getBaseUrl()}${it.removePrefix("/")}"
                it.startsWith("/") -> "${SharedPrefsManager.getBaseUrl()}${it.removePrefix("/")}"
                else -> "${SharedPrefsManager.getBaseUrl()}static/posters/$it"
            }
        }

        android.util.Log.d(TAG, "loadPosterImage: playlistPoster=$playlistPoster")
        android.util.Log.d(TAG, "loadPosterImage: posterUrl=$posterUrl")

        if (posterUrl != null) {
            Glide.with(requireContext())
                .asBitmap()
                .load(posterUrl)
                .error(R.drawable.ic_playlist_placeholder)
                .into(object : CustomTarget<Bitmap>() {
                    override fun onResourceReady(
                        resource: Bitmap,
                        transition: Transition<in Bitmap>?
                    ) {
                        android.util.Log.d(TAG, "loadPosterImage: Poster loaded successfully")
                        posterBitmap = resource
                        row.setImageBitmap(requireContext(), resource)
                        rowsAdapter.notifyArrayItemRangeChanged(0, 1)
                        backgroundController.coverBitmap = resource
                    }

                    override fun onLoadCleared(placeholder: Drawable?) {
                        // No-op
                    }
                })
        } else {
            android.util.Log.d(TAG, "loadPosterImage: No poster URL available")
        }
    }

    /**
     * Load playlist content based on hasChildren flag
     * If hasChildren == true: load child playlists
     * If hasChildren == false: load media items
     */
    private fun loadPlaylistContent() {
        android.util.Log.d(TAG, "loadPlaylistContent: called")
        lifecycleScope.launch {
            try {
                playlistId?.let { id ->
                    if (playlistHasChildren) {
                        android.util.Log.d(TAG, "loadPlaylistContent: Loading child playlists for $id")
                        loadPlaylistChildren(id)
                    } else {
                        android.util.Log.d(TAG, "loadPlaylistContent: Loading media for $id")
                        loadPlaylistMedia(id)
                    }
                }
            } catch (e: Exception) {
                android.util.Log.e(TAG, "loadPlaylistContent: Exception", e)
                showError(e.message ?: getString(R.string.error))
            }
        }
    }

    /**
     * Load child playlists
     */
    private suspend fun loadPlaylistChildren(playlistId: String) {
        android.util.Log.d(TAG, "loadPlaylistChildren: Starting request for $playlistId")
        val result = playlistRepository.getPlaylistChildren(playlistId, limit = 100, offset = 0)
        result.onSuccess { playlists ->
            android.util.Log.d(TAG, "loadPlaylistChildren: Success! Got ${playlists.size} playlists")
            displayChildPlaylists(playlists)
        }.onFailure { error ->
            android.util.Log.e(TAG, "loadPlaylistChildren: Failed", error)
            showError(error.message ?: getString(R.string.playlists_error))
        }
    }

    /**
     * Load media items
     */
    private suspend fun loadPlaylistMedia(playlistId: String) {
        val result = mediaRepository.getPlaylistMedia(playlistId, limit = 100, offset = 0)
        result.onSuccess { apiResponse ->
            if (apiResponse.isSuccess() && apiResponse.result != null) {
                val mediaList = apiResponse.result
                displayMediaList(mediaList)
            } else {
                showError(apiResponse.error ?: getString(R.string.media_error))
            }
        }.onFailure { error ->
            showError(error.message ?: getString(R.string.media_error))
        }
    }

    /**
     * Display child playlists
     */
    private fun displayChildPlaylists(playlists: List<Playlist>) {
        android.util.Log.d(TAG, "displayChildPlaylists: Displaying ${playlists.size} playlists")

        if (playlists.isEmpty()) {
            android.util.Log.d(TAG, "displayChildPlaylists: Empty playlists, showing empty message")
            // Create new row with empty message
            val emptyRow = DetailsOverviewRow(getString(R.string.playlists_empty))
            if (posterBitmap != null) {
                emptyRow.setImageBitmap(requireContext(), posterBitmap)
            } else {
                emptyRow.setImageDrawable(
                    ContextCompat.getDrawable(
                        requireContext(),
                        R.drawable.ic_playlist_placeholder
                    )
                )
            }
            rowsAdapter.replace(0, emptyRow)
            return
        }

        android.util.Log.d(TAG, "displayChildPlaylists: Updating details row with count")
        // Create new row with playlists count
        val updatedRow = DetailsOverviewRow("${playlists.size} ${getString(R.string.playlists_items)}")
        // Set poster image if available
        if (posterBitmap != null) {
            updatedRow.setImageBitmap(requireContext(), posterBitmap)
        } else {
            updatedRow.setImageDrawable(
                ContextCompat.getDrawable(
                    requireContext(),
                    R.drawable.ic_playlist_placeholder
                )
            )
        }
        rowsAdapter.replace(0, updatedRow)

        android.util.Log.d(TAG, "displayChildPlaylists: Creating list row with playlists")
        // Add playlists list
        val cardPresenter = CardPresenter()
        val listRowAdapter = ArrayObjectAdapter(cardPresenter)
        playlists.forEach { playlist ->
            android.util.Log.d(TAG, "displayChildPlaylists: Adding playlist: ${playlist.name}")
            listRowAdapter.add(playlist)
        }

        val header = HeaderItem(0, getString(R.string.main_playlists))
        rowsAdapter.add(ListRow(header, listRowAdapter))

        android.util.Log.d(TAG, "displayChildPlaylists: Adapter now has ${rowsAdapter.size()} rows")
        android.util.Log.d(TAG, "displayChildPlaylists: Done!")
    }

    private fun displayMediaList(mediaList: List<Media>) {
        android.util.Log.d(TAG, "displayMediaList: Displaying ${mediaList.size} media items")

        if (mediaList.isEmpty()) {
            android.util.Log.d(TAG, "displayMediaList: Empty media list, showing empty message")
            // Create new row with empty message
            val emptyRow = DetailsOverviewRow(getString(R.string.media_empty))
            if (posterBitmap != null) {
                emptyRow.setImageBitmap(requireContext(), posterBitmap)
            } else {
                emptyRow.setImageDrawable(
                    ContextCompat.getDrawable(
                        requireContext(),
                        R.drawable.ic_playlist_placeholder
                    )
                )
            }
            rowsAdapter.replace(0, emptyRow)
            return
        }

        android.util.Log.d(TAG, "displayMediaList: Updating details row with count")
        // Create new row with media count
        val updatedRow = DetailsOverviewRow("${mediaList.size} ${getString(R.string.media_items)}")
        // Set poster image if available
        if (posterBitmap != null) {
            updatedRow.setImageBitmap(requireContext(), posterBitmap)
        } else {
            updatedRow.setImageDrawable(
                ContextCompat.getDrawable(
                    requireContext(),
                    R.drawable.ic_playlist_placeholder
                )
            )
        }
        rowsAdapter.replace(0, updatedRow)

        android.util.Log.d(TAG, "displayMediaList: Creating list row with media")
        // Add media list
        val cardPresenter = CardPresenter()
        val listRowAdapter = ArrayObjectAdapter(cardPresenter)
        mediaList.forEach { media ->
            android.util.Log.d(TAG, "displayMediaList: Adding media: ${media.name}")
            listRowAdapter.add(media)
        }

        val header = HeaderItem(0, getString(R.string.main_media))
        rowsAdapter.add(ListRow(header, listRowAdapter))

        android.util.Log.d(TAG, "displayMediaList: Adapter now has ${rowsAdapter.size()} rows")
        android.util.Log.d(TAG, "displayMediaList: Done!")
    }

    private fun openMediaPlayer(media: Media) {
        val intent = Intent(requireContext(), PlayerActivity::class.java).apply {
            putExtra(PlayerActivity.EXTRA_MEDIA_ID, media.id)
            putExtra(PlayerActivity.EXTRA_MEDIA_NAME, media.name)
        }
        startActivity(intent)
    }

    private fun openPlaylistDetails(playlist: Playlist) {
        val intent = Intent(requireContext(), PlaylistDetailsActivity::class.java).apply {
            putExtra(PlaylistDetailsActivity.EXTRA_PLAYLIST_ID, playlist.id)
            putExtra(PlaylistDetailsActivity.EXTRA_PLAYLIST_NAME, playlist.name)
            putExtra(PlaylistDetailsActivity.EXTRA_PLAYLIST_POSTER, playlist.poster)
            putExtra(PlaylistDetailsActivity.EXTRA_PLAYLIST_HAS_CHILDREN, playlist.hasChildren)
        }
        startActivity(intent)
    }

    private fun showError(message: String) {
        activity?.runOnUiThread {
            Toast.makeText(requireContext(), message, Toast.LENGTH_LONG).show()
        }
    }
}
