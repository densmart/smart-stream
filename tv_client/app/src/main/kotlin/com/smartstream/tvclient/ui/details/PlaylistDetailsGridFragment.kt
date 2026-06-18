package com.smartstream.tvclient.ui.details

import android.content.Intent
import android.graphics.Rect
import android.os.Bundle
import android.util.Log
import android.view.KeyEvent
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ImageView
import android.widget.TextView
import android.widget.Toast
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import androidx.recyclerview.widget.GridLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.bumptech.glide.Glide
import com.smartstream.tvclient.R
import com.smartstream.tvclient.data.model.Media
import com.smartstream.tvclient.data.model.Playlist
import com.smartstream.tvclient.data.repository.MediaRepository
import com.smartstream.tvclient.data.repository.PlaylistRepository
import com.smartstream.tvclient.ui.main.MediaCardAdapter
import com.smartstream.tvclient.ui.player.PlayerActivity
import com.smartstream.tvclient.utils.PaginationHelper
import com.smartstream.tvclient.utils.SharedPrefsManager
import kotlinx.coroutines.launch

/**
 * Fragment for displaying playlist children in a grid layout (like main screen)
 * with breadcrumb navigation
 */
class PlaylistDetailsGridFragment : Fragment() {

    private val playlistRepository = PlaylistRepository()
    private val mediaRepository = MediaRepository()

    private lateinit var recyclerView: RecyclerView
    private lateinit var adapter: MediaCardAdapter
    private lateinit var breadcrumbText: TextView
    private lateinit var btnSearch: View

    private lateinit var detailPoster: ImageView
    private lateinit var detailTitle: TextView
    private lateinit var detailInfo: TextView
    private lateinit var detailCardContainer: View

    private lateinit var paginationHelper: PaginationHelper

    // Grid decoration
    private var currentItemDecoration: RecyclerView.ItemDecoration? = null

    private var playlistId: String? = null
    private var playlistName: String? = null
    private var playlistPoster: String? = null
    private var playlistHasChildren: Boolean = false
    private var currentOffset = 0

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        return inflater.inflate(R.layout.fragment_playlist_details_grid, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        // Get extras from activity
        activity?.intent?.let { intent ->
            playlistId = intent.getStringExtra(PlaylistDetailsActivity.EXTRA_PLAYLIST_ID)
            playlistName = intent.getStringExtra(PlaylistDetailsActivity.EXTRA_PLAYLIST_NAME)
            playlistPoster = intent.getStringExtra(PlaylistDetailsActivity.EXTRA_PLAYLIST_POSTER)
            playlistHasChildren = intent.getBooleanExtra(PlaylistDetailsActivity.EXTRA_PLAYLIST_HAS_CHILDREN, false)
        }

        initViews(view)
        setupPagination(view)
        setupRecyclerView()
        setupBreadcrumb()

        // Load playlist children
        loadPlaylistChildren()
    }

    private fun initViews(view: View) {
        recyclerView = view.findViewById(R.id.cards_recycler_view)
        breadcrumbText = view.findViewById(R.id.breadcrumb_text)
        btnSearch = view.findViewById(R.id.btn_search)
        detailCardContainer = view.findViewById(R.id.detail_card_container)
        detailPoster = view.findViewById(R.id.detail_poster)
        detailTitle = view.findViewById(R.id.detail_title)
        detailInfo = view.findViewById(R.id.detail_info)

        // Setup search button click
        btnSearch.setOnClickListener {
            Toast.makeText(requireContext(), "Search not implemented yet", Toast.LENGTH_SHORT).show()
        }
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
                    loadPlaylistChildren()
                }
            },
            onNextClick = {
                paginationHelper.goNext()?.let { newOffset ->
                    currentOffset = newOffset
                    loadPlaylistChildren()
                }
            }
        )
    }

    private fun setupRecyclerView() {
        adapter = MediaCardAdapter(
            onItemClick = { item ->
                when (item) {
                    is Media -> openMediaPlayer(item)
                    is Playlist -> openPlaylistDetails(item)
                }
            },
            onItemFocus = { item ->
                updateDetailCard(item)
            },
            onFocusLost = {
                // Will be handled by recycler view focus listener
            }
        )

        recyclerView.adapter = adapter

        // Use standard GridLayoutManager - custom navigation was causing focus and layout issues
        recyclerView.layoutManager = GridLayoutManager(requireContext(), 6)

        // Add spacing between cards
        val spacing = resources.getDimensionPixelSize(R.dimen.card_margin)
        currentItemDecoration = GridSpacingItemDecoration(6, spacing, true)
        recyclerView.addItemDecoration(currentItemDecoration!!)
    }

    /**
     * ItemDecoration for adding spacing between grid items
     */
    private class GridSpacingItemDecoration(
        private val spanCount: Int,
        private val spacing: Int,
        private val includeEdge: Boolean
    ) : RecyclerView.ItemDecoration() {

        override fun getItemOffsets(
            outRect: Rect,
            view: View,
            parent: RecyclerView,
            state: RecyclerView.State
        ) {
            val position = parent.getChildAdapterPosition(view)
            val column = position % spanCount

            if (includeEdge) {
                outRect.left = spacing - column * spacing / spanCount
                outRect.right = (column + 1) * spacing / spanCount

                if (position < spanCount) {
                    outRect.top = spacing
                }
                outRect.bottom = spacing
            } else {
                outRect.left = column * spacing / spanCount
                outRect.right = spacing - (column + 1) * spacing / spanCount
                if (position >= spanCount) {
                    outRect.top = spacing
                }
            }
        }
    }

    private fun setupBreadcrumb() {
        // Build breadcrumb: "Playlists / Parent Name"
        val breadcrumb = buildString {
            append(getString(R.string.main_playlists))
            if (playlistName != null) {
                append(" / ")
                append(playlistName)
            }
        }
        breadcrumbText.text = breadcrumb
    }

    private fun showDetailCard() {
        if (detailCardContainer.visibility == View.GONE) {
            detailCardContainer.visibility = View.VISIBLE
            updateGridLayout(4)
        }
    }

    private fun hideDetailCard() {
        if (detailCardContainer.visibility == View.VISIBLE) {
            detailCardContainer.visibility = View.GONE
            updateGridLayout(6)
        }
    }

    private fun updateGridLayout(spanCount: Int) {
        // Post the update to avoid modifying RecyclerView during layout pass
        recyclerView.post {
            (recyclerView.layoutManager as? GridLayoutManager)?.let { layoutManager ->
                // Remove old decoration
                currentItemDecoration?.let { decoration ->
                    recyclerView.removeItemDecoration(decoration)
                }

                // Update span count
                layoutManager.spanCount = spanCount

                // Add new decoration with updated span count
                val spacing = resources.getDimensionPixelSize(R.dimen.card_margin)
                currentItemDecoration = GridSpacingItemDecoration(spanCount, spacing, true)
                recyclerView.addItemDecoration(currentItemDecoration!!)

                // Request layout recalculation
                recyclerView.requestLayout()
            }
        }
    }

    private fun loadPlaylistChildren() {
        lifecycleScope.launch {
            try {
                playlistId?.let { id ->
                    Log.d(TAG, "Loading children for playlist: $id")
                    val result = playlistRepository.getPlaylistChildrenWithPagination(
                        playlistId = id,
                        offset = currentOffset
                    )
                    result.onSuccess { apiResponse ->
                        if (apiResponse.isSuccess() && apiResponse.result != null) {
                            Log.d(TAG, "Loaded ${apiResponse.result.size} child playlists")
                            adapter.submitList(apiResponse.result)
                            paginationHelper.update(apiResponse.pagination, currentOffset)

                            // Request focus on first item after data is loaded
                            recyclerView.postDelayed({
                                if (adapter.itemCount > 0 && recyclerView.childCount > 0) {
                                    recyclerView.getChildAt(0)?.requestFocus()
                                }
                            }, 100)
                        } else {
                            Log.e(TAG, "API error: ${apiResponse.error}")
                            showError(apiResponse.error ?: getString(R.string.playlists_error))
                        }
                    }.onFailure { error ->
                        Log.e(TAG, "Failed to load children", error)
                        showError(error.message ?: getString(R.string.playlists_error))
                    }
                }
            } catch (e: Exception) {
                Log.e(TAG, "Exception loading children", e)
                showError(e.message ?: getString(R.string.error))
            }
        }
    }

    private fun updateDetailCard(item: Any) {
        // Show detail card and adjust grid columns
        showDetailCard()

        when (item) {
            is Media -> {
                detailTitle.text = item.name
                detailTitle.visibility = View.VISIBLE
                detailInfo.text = buildMediaInfo(item)
                detailInfo.visibility = View.VISIBLE

                val posterUrl = item.getPosterUrl(SharedPrefsManager.getBaseUrl())
                if (posterUrl != null) {
                    Glide.with(this)
                        .load(posterUrl)
                        .error(R.drawable.ic_media_placeholder)
                        .into(detailPoster)
                } else {
                    detailPoster.setImageResource(R.drawable.ic_media_placeholder)
                }
            }
            is Playlist -> {
                detailTitle.text = item.name
                detailTitle.visibility = View.VISIBLE
                detailInfo.text = item.getTypeDisplayName(requireContext())
                detailInfo.visibility = View.VISIBLE

                val posterUrl = item.getPosterUrl(SharedPrefsManager.getBaseUrl())
                if (posterUrl != null) {
                    Glide.with(this)
                        .load(posterUrl)
                        .error(R.drawable.ic_playlist_placeholder)
                        .into(detailPoster)
                } else {
                    detailPoster.setImageResource(R.drawable.ic_playlist_placeholder)
                }
            }
        }
    }

    private fun buildMediaInfo(media: Media): String {
        val parts = mutableListOf<String>()
        media.getFormattedDuration()?.let { parts.add(it) }
        media.getFormattedSize()?.let { parts.add(it) }
        return parts.joinToString(" • ")
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

    companion object {
        private const val TAG = "PlaylistDetailsGrid"
    }
}
