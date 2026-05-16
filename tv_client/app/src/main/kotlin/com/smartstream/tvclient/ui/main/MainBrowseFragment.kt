package com.smartstream.tvclient.ui.main

import android.content.Intent
import android.os.Bundle
import android.widget.Toast
import androidx.core.content.ContextCompat
import androidx.leanback.app.BrowseSupportFragment
import androidx.leanback.widget.ArrayObjectAdapter
import androidx.leanback.widget.HeaderItem
import androidx.leanback.widget.ListRow
import androidx.leanback.widget.ListRowPresenter
import androidx.leanback.widget.OnItemViewClickedListener
import androidx.lifecycle.lifecycleScope
import com.smartstream.tvclient.R
import com.smartstream.tvclient.data.model.Media
import com.smartstream.tvclient.data.model.Playlist
import com.smartstream.tvclient.data.repository.MediaRepository
import com.smartstream.tvclient.data.repository.PlaylistRepository
import com.smartstream.tvclient.ui.details.MediaDetailsActivity
import com.smartstream.tvclient.ui.details.PlaylistDetailsActivity
import kotlinx.coroutines.launch

/**
 * Main browse fragment showing playlists and media categories.
 * Uses Leanback BrowseSupportFragment for TV-optimized navigation.
 */
class MainBrowseFragment : BrowseSupportFragment() {

    private val playlistRepository = PlaylistRepository()
    private val mediaRepository = MediaRepository()

    private lateinit var rowsAdapter: ArrayObjectAdapter

    override fun onActivityCreated(savedInstanceState: Bundle?) {
        super.onActivityCreated(savedInstanceState)

        setupUI()
        setupEventListeners()
        loadContent()
    }

    private fun setupUI() {
        // Set badge and title
        badgeDrawable = ContextCompat.getDrawable(requireContext(), R.drawable.ic_app_logo)
        title = getString(R.string.main_title)

        // Set headers state
        headersState = HEADERS_ENABLED
        isHeadersTransitionOnBackEnabled = true

        // Set brand color
        brandColor = ContextCompat.getColor(requireContext(), R.color.primary)

        // Set search icon color
        searchAffordanceColor = ContextCompat.getColor(requireContext(), R.color.primary_light)
    }

    private fun setupEventListeners() {
        // Item click listener
        setOnItemViewClickedListener(OnItemViewClickedListener { _, item, _, _ ->
            android.util.Log.d("MainBrowseFragment", "Item clicked: $item")
            when (item) {
                is Playlist -> {
                    android.util.Log.d("MainBrowseFragment", "Opening playlist details for: ${item.name}")
                    openPlaylistDetails(item)
                }
                is Media -> {
                    android.util.Log.d("MainBrowseFragment", "Opening media details for: ${item.name}")
                    openMediaDetails(item)
                }
                else -> {
                    android.util.Log.w("MainBrowseFragment", "Unknown item type clicked: ${item?.javaClass?.simpleName}")
                }
            }
        })
    }

    private fun loadContent() {
        rowsAdapter = ArrayObjectAdapter(ListRowPresenter())

        viewLifecycleOwner.lifecycleScope.launch {
            try {
                // Load playlists
                loadPlaylists()

                // Load unassigned media
                loadUnassignedMedia()

                adapter = rowsAdapter
            } catch (e: Exception) {
                showError(getString(R.string.error) + ": ${e.message}")
            }
        }
    }

    private suspend fun loadPlaylists() {
        try {
            val result = playlistRepository.getPlaylists()
            result.onSuccess { playlists ->
                if (playlists.isNotEmpty()) {
                    val cardPresenter = CardPresenter()
                    val listRowAdapter = ArrayObjectAdapter(cardPresenter)

                    playlists.forEach { playlist ->
                        listRowAdapter.add(playlist)
                    }

                    val header = HeaderItem(0, getString(R.string.main_playlists))
                    rowsAdapter.add(ListRow(header, listRowAdapter))
                }
            }.onFailure { error ->
                showError(getString(R.string.playlists_error) + ": ${error.message}")
            }
        } catch (e: Exception) {
            showError(getString(R.string.playlists_error) + ": ${e.message}")
        }
    }

    private suspend fun loadUnassignedMedia() {
        try {
            val result = mediaRepository.getMediaList(onlyUnassigned = true)
            result.onSuccess { mediaList ->
                if (mediaList.isNotEmpty()) {
                    val cardPresenter = CardPresenter()
                    val listRowAdapter = ArrayObjectAdapter(cardPresenter)

                    mediaList.forEach { media ->
                        listRowAdapter.add(media)
                    }

                    val header = HeaderItem(1, getString(R.string.media_unassigned))
                    rowsAdapter.add(ListRow(header, listRowAdapter))
                }
            }.onFailure { error ->
                showError(getString(R.string.media_error) + ": ${error.message}")
            }
        } catch (e: Exception) {
            showError(getString(R.string.media_error) + ": ${e.message}")
        }
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

    private fun openMediaDetails(media: Media) {
        android.util.Log.d("MainBrowseFragment", "openMediaDetails: Creating intent for media: ${media.id}")
        android.util.Log.d("MainBrowseFragment", "openMediaDetails: Media details - name=${media.name}, duration=${media.duration}, size=${media.size}")

        val intent = Intent(requireContext(), MediaDetailsActivity::class.java).apply {
            putExtra(MediaDetailsActivity.EXTRA_MEDIA_ID, media.id)
            putExtra(MediaDetailsActivity.EXTRA_MEDIA_NAME, media.name)
            putExtra(MediaDetailsActivity.EXTRA_MEDIA_POSTER, media.poster)
            putExtra(MediaDetailsActivity.EXTRA_MEDIA_DURATION, media.duration ?: 0)
            putExtra(MediaDetailsActivity.EXTRA_MEDIA_SIZE, media.size ?: 0L)
        }

        android.util.Log.d("MainBrowseFragment", "openMediaDetails: Starting MediaDetailsActivity")
        try {
            startActivity(intent)
            android.util.Log.d("MainBrowseFragment", "openMediaDetails: MediaDetailsActivity started successfully")
        } catch (e: Exception) {
            android.util.Log.e("MainBrowseFragment", "openMediaDetails: Failed to start activity", e)
            Toast.makeText(requireContext(), "Error: ${e.message}", Toast.LENGTH_LONG).show()
        }
    }

    private fun showError(message: String) {
        activity?.runOnUiThread {
            Toast.makeText(requireContext(), message, Toast.LENGTH_LONG).show()
        }
    }
}
