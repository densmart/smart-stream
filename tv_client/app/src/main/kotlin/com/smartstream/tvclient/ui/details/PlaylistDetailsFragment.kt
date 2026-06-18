package com.smartstream.tvclient.ui.details

import android.content.Intent
import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ImageView
import android.widget.TextView
import android.widget.Toast
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.bumptech.glide.Glide
import com.smartstream.tvclient.R
import com.smartstream.tvclient.data.model.Media
import com.smartstream.tvclient.data.repository.MediaRepository
import com.smartstream.tvclient.ui.player.PlayerActivity
import com.smartstream.tvclient.utils.PaginationHelper
import com.smartstream.tvclient.utils.SharedPrefsManager
import kotlinx.coroutines.launch

/**
 * Fragment displaying playlist details with vertical list of media items
 */
class PlaylistDetailsFragment : Fragment() {

    private val mediaRepository = MediaRepository()

    private lateinit var backgroundPoster: ImageView
    private lateinit var playlistPoster: ImageView
    private lateinit var playlistName: TextView
    private lateinit var episodesCount: TextView
    private lateinit var mediaList: RecyclerView
    private lateinit var adapter: MediaListAdapter

    private lateinit var paginationHelper: PaginationHelper

    private var playlistId: String? = null
    private var playlistNameText: String? = null
    private var playlistPosterPath: String? = null
    private var currentOffset = 0

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        return inflater.inflate(R.layout.fragment_playlist_details_media, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        // Get extras from activity
        activity?.intent?.let { intent ->
            playlistId = intent.getStringExtra(PlaylistDetailsActivity.EXTRA_PLAYLIST_ID)
            playlistNameText = intent.getStringExtra(PlaylistDetailsActivity.EXTRA_PLAYLIST_NAME)
            playlistPosterPath = intent.getStringExtra(PlaylistDetailsActivity.EXTRA_PLAYLIST_POSTER)
        }

        initViews(view)
        setupPagination(view)
        setupRecyclerView()
        loadPoster()
        loadMedia()
    }

    private fun initViews(view: View) {
        backgroundPoster = view.findViewById(R.id.background_poster)
        playlistPoster = view.findViewById(R.id.playlist_poster)
        playlistName = view.findViewById(R.id.playlist_name)
        episodesCount = view.findViewById(R.id.episodes_count)
        mediaList = view.findViewById(R.id.media_list)

        // Set playlist name
        playlistName.text = playlistNameText ?: getString(R.string.playlists_loading)
    }

    private fun setupPagination(view: View) {
        val paginationContainer = view.findViewById<View>(R.id.pagination_container)
        val btnPrevious = view.findViewById<TextView>(R.id.btn_previous)
        val btnNext = view.findViewById<TextView>(R.id.btn_next)
        val paginationInfo = view.findViewById<TextView>(R.id.pagination_info)

        paginationHelper = PaginationHelper(
            paginationContainer,
            btnPrevious,
            btnNext,
            paginationInfo
        )

        paginationHelper.setupListeners(
            onPreviousClick = {
                paginationHelper.goPrevious()?.let { newOffset ->
                    currentOffset = newOffset
                    loadMedia()
                }
            },
            onNextClick = {
                paginationHelper.goNext()?.let { newOffset ->
                    currentOffset = newOffset
                    loadMedia()
                }
            }
        )
    }

    private fun setupRecyclerView() {
        adapter = MediaListAdapter(playlistId) { media ->
            openMediaPlayer(media)
        }

        mediaList.adapter = adapter
        mediaList.layoutManager = LinearLayoutManager(requireContext())

        // Request focus on first item after layout
        mediaList.postDelayed({
            if (adapter.itemCount > 0) {
                mediaList.getChildAt(0)?.requestFocus()
            }
        }, 100)
    }

    override fun onResume() {
        super.onResume()
        // Refresh adapter to update "continue watching" display
        adapter.refreshList()
    }

    private fun loadPoster() {
        val posterUrl = playlistPosterPath?.let {
            when {
                it.startsWith("http") -> it
                it.startsWith("/static/posters/") -> "${SharedPrefsManager.getBaseUrl()}${it.removePrefix("/")}"
                it.startsWith("/") -> "${SharedPrefsManager.getBaseUrl()}${it.removePrefix("/")}"
                else -> "${SharedPrefsManager.getBaseUrl()}static/posters/$it"
            }
        }

        android.util.Log.d(TAG, "loadPoster: posterUrl=$posterUrl")

        if (posterUrl != null) {
            // Load small poster
            Glide.with(requireContext())
                .load(posterUrl)
                .centerCrop()
                .error(R.drawable.ic_playlist_placeholder)
                .placeholder(R.drawable.ic_playlist_placeholder)
                .into(playlistPoster)

            // Load background poster
            Glide.with(requireContext())
                .load(posterUrl)
                .centerCrop()
                .error(R.drawable.ic_playlist_placeholder)
                .placeholder(R.drawable.ic_playlist_placeholder)
                .into(backgroundPoster)
        } else {
            android.util.Log.d(TAG, "loadPoster: No poster URL available")
            playlistPoster.setImageResource(R.drawable.ic_playlist_placeholder)
            backgroundPoster.setImageResource(R.drawable.ic_playlist_placeholder)
        }
    }

    private fun loadMedia() {
        lifecycleScope.launch {
            try {
                playlistId?.let { id ->
                    android.util.Log.d(TAG, "loadMedia: Loading media for playlist $id")
                    val result = mediaRepository.getPlaylistMedia(
                        playlistId = id,
                        offset = currentOffset
                    )
                    result.onSuccess { apiResponse ->
                        if (apiResponse.isSuccess() && apiResponse.result != null) {
                            val media = apiResponse.result
                            android.util.Log.d(TAG, "loadMedia: Got ${media.size} media items")

                            // Update episodes count with total from pagination
                            val totalCount = apiResponse.pagination?.total ?: media.size.toLong()
                            episodesCount.text = getString(R.string.episodes_count, totalCount.toInt())

                            // Submit sorted list to adapter
                            adapter.submitList(media)
                            paginationHelper.update(apiResponse.pagination, currentOffset)

                            // Request focus on first item
                            mediaList.postDelayed({
                                if (adapter.itemCount > 0) {
                                    mediaList.getChildAt(0)?.requestFocus()
                                }
                            }, 200)
                        } else {
                            android.util.Log.e(TAG, "loadMedia: API error: ${apiResponse.error}")
                            showError(apiResponse.error ?: getString(R.string.media_error))
                        }
                    }.onFailure { error ->
                        android.util.Log.e(TAG, "loadMedia: Failed", error)
                        showError(error.message ?: getString(R.string.media_error))
                    }
                }
            } catch (e: Exception) {
                android.util.Log.e(TAG, "loadMedia: Exception", e)
                showError(e.message ?: getString(R.string.error))
            }
        }
    }

    private fun openMediaPlayer(media: Media) {
        // Get saved position - check if this is the current episode in playlist
        val pId = playlistId
        val savedPosition = if (pId != null) {
            val currentMediaId = SharedPrefsManager.getPlaylistCurrentMedia(pId)
            if (currentMediaId == media.id) {
                // This is the current episode, get its position
                SharedPrefsManager.getPlaylistPosition(pId)
            } else {
                // Different episode, start from beginning
                0L
            }
        } else {
            // Single media, get its saved position
            SharedPrefsManager.getMediaPosition(media.id)
        }

        val intent = Intent(requireContext(), PlayerActivity::class.java).apply {
            putExtra(PlayerActivity.EXTRA_MEDIA_ID, media.id)
            putExtra(PlayerActivity.EXTRA_MEDIA_NAME, media.name)
            // Pass playlist ID for autoplay
            putExtra(PlayerActivity.EXTRA_PLAYLIST_ID, pId)
            // Pass saved position for continue watching
            putExtra(PlayerActivity.EXTRA_START_POSITION, savedPosition)
        }
        startActivity(intent)
    }

    private fun showError(message: String) {
        activity?.runOnUiThread {
            Toast.makeText(requireContext(), message, Toast.LENGTH_LONG).show()
        }
    }

    companion object {
        private const val TAG = "PlaylistDetailsFrag"
    }
}
