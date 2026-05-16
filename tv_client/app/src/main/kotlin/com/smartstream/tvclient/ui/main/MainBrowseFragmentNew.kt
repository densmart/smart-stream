package com.smartstream.tvclient.ui.main

import android.content.Intent
import android.graphics.Color
import android.graphics.Rect
import android.os.Bundle
import android.util.Log
import android.view.KeyEvent
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ImageButton
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
import com.smartstream.tvclient.ui.details.MediaDetailsActivity
import com.smartstream.tvclient.ui.details.PlaylistDetailsActivity
import com.smartstream.tvclient.ui.settings.ServerSettingsDialogFragment
import com.smartstream.tvclient.utils.SharedPrefsManager
import kotlinx.coroutines.launch

/**
 * Main browse fragment with custom design
 * Shows tabs for Media/Playlists with grid of cards and detail panel
 */
class MainBrowseFragmentNew : Fragment() {

    private val playlistRepository = PlaylistRepository()
    private val mediaRepository = MediaRepository()

    private lateinit var recyclerView: RecyclerView
    private lateinit var adapter: MediaCardAdapter

    private lateinit var tabMedia: TextView
    private lateinit var tabPlaylists: TextView
    private lateinit var btnSettings: ImageButton
    private lateinit var btnSearch: ImageButton

    private lateinit var detailPoster: ImageView
    private lateinit var detailTitle: TextView
    private lateinit var detailInfo: TextView
    private lateinit var detailCardContainer: View

    private var currentTab = Tab.MEDIA

    enum class Tab {
        MEDIA, PLAYLISTS
    }

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        return inflater.inflate(R.layout.fragment_main_browse_custom, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        initViews(view)
        setupRecyclerView()
        setupTabButtons()
        setupSettingsButton()

        // Load initial content (Media by default)
        loadMediaContent()
    }

    private fun initViews(view: View) {
        recyclerView = view.findViewById(R.id.cards_recycler_view)
        tabMedia = view.findViewById(R.id.tab_media)
        tabPlaylists = view.findViewById(R.id.tab_playlists)
        btnSettings = view.findViewById(R.id.btn_settings)
        btnSearch = view.findViewById(R.id.btn_search)

        detailCardContainer = view.findViewById(R.id.detail_card_container)

        // Force visibility and text properties
        tabMedia.visibility = View.VISIBLE
        tabPlaylists.visibility = View.VISIBLE
        tabMedia.alpha = 1f
        tabPlaylists.alpha = 1f

        // Set text programmatically to ensure it's there
        tabMedia.text = "Media"
        tabPlaylists.text = "Playlists"
        detailPoster = view.findViewById(R.id.detail_poster)
        detailTitle = view.findViewById(R.id.detail_title)
        detailInfo = view.findViewById(R.id.detail_info)
    }

    private fun setupRecyclerView() {
        adapter = MediaCardAdapter(
            onItemClick = { item ->
                when (item) {
                    is Media -> openMediaDetails(item)
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

        // Create custom GridLayoutManager with fixed navigation
        val customLayoutManager = object : GridLayoutManager(requireContext(), 6) {
            override fun onInterceptFocusSearch(focused: View, direction: Int): View? {
                val position = getPosition(focused)
                if (position == RecyclerView.NO_POSITION) {
                    return super.onInterceptFocusSearch(focused, direction)
                }

                val spanCount = this.spanCount
                val itemCount = adapter.itemCount
                val currentRow = position / spanCount
                val currentColumn = position % spanCount

                when (direction) {
                    View.FOCUS_RIGHT -> {
                        // If on last column, go to first item of next row
                        if (currentColumn == spanCount - 1) {
                            val nextRowFirstItem = (currentRow + 1) * spanCount
                            if (nextRowFirstItem < itemCount) {
                                // Scroll to make sure the view is visible
                                recyclerView.post {
                                    recyclerView.smoothScrollToPosition(nextRowFirstItem)
                                    recyclerView.postDelayed({
                                        findViewByPosition(nextRowFirstItem)?.requestFocus()
                                    }, 100)
                                }
                                // Return focused to prevent default behavior
                                return focused
                            } else {
                                // No next row, stay on current item
                                return focused
                            }
                        } else if (position == itemCount - 1) {
                            // Last item overall, prevent moving right
                            return focused
                        }
                    }
                    View.FOCUS_LEFT -> {
                        // If on first column, go to last item of previous row
                        if (currentColumn == 0) {
                            if (currentRow > 0) {
                                val prevRowLastItem = currentRow * spanCount - 1
                                if (prevRowLastItem >= 0) {
                                    recyclerView.post {
                                        recyclerView.smoothScrollToPosition(prevRowLastItem)
                                        recyclerView.postDelayed({
                                            findViewByPosition(prevRowLastItem)?.requestFocus()
                                        }, 100)
                                    }
                                    return focused
                                }
                            } else {
                                // First row, first column - prevent moving left
                                return focused
                            }
                        }
                    }
                }
                return super.onInterceptFocusSearch(focused, direction)
            }
        }

        recyclerView.layoutManager = customLayoutManager

        // Add spacing between cards
        val spacing = resources.getDimensionPixelSize(R.dimen.card_margin)
        recyclerView.addItemDecoration(GridSpacingItemDecoration(6, spacing, true))

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

    private fun showDetailCard() {
        if (detailCardContainer.visibility == View.GONE) {
            detailCardContainer.visibility = View.VISIBLE
            // Change to 4 columns when preview is shown
            (recyclerView.layoutManager as? GridLayoutManager)?.spanCount = 4
        }
    }

    private fun hideDetailCard() {
        if (detailCardContainer.visibility == View.VISIBLE) {
            detailCardContainer.visibility = View.GONE
            // Change back to 6 columns when preview is hidden
            (recyclerView.layoutManager as? GridLayoutManager)?.spanCount = 6
        }
    }

    private fun setupTabButtons() {
        tabMedia.setOnClickListener {
            switchTab(Tab.MEDIA)
        }

        tabPlaylists.setOnClickListener {
            switchTab(Tab.PLAYLISTS)
        }

        // Intercept DOWN key to focus first card
        val downKeyListener = View.OnKeyListener { _, keyCode, event ->
            if (keyCode == KeyEvent.KEYCODE_DPAD_DOWN && event.action == KeyEvent.ACTION_DOWN) {
                recyclerView.post {
                    recyclerView.scrollToPosition(0)
                    recyclerView.postDelayed({
                        val firstViewHolder = recyclerView.findViewHolderForAdapterPosition(0)
                        firstViewHolder?.itemView?.requestFocus()
                    }, 50)
                }
                true
            } else {
                false
            }
        }

        tabMedia.setOnKeyListener(downKeyListener)
        tabPlaylists.setOnKeyListener(downKeyListener)
        btnSettings.setOnKeyListener(downKeyListener)

        // Hide preview when tabs receive focus and update text color
        val grayColor = Color.parseColor("#999999")

        tabMedia.setOnFocusChangeListener { _, hasFocus ->
            if (hasFocus) {
                hideDetailCard()
                // Dark text on focused background (light gray)
                tabMedia.setTextColor(Color.BLACK)
            } else {
                // Restore color based on selection: selected = black, not selected = gray
                tabMedia.setTextColor(if (currentTab == Tab.MEDIA) Color.BLACK else grayColor)
            }
        }

        tabPlaylists.setOnFocusChangeListener { _, hasFocus ->
            if (hasFocus) {
                hideDetailCard()
                // Dark text on focused background (light gray)
                tabPlaylists.setTextColor(Color.BLACK)
            } else {
                // Restore color based on selection: selected = black, not selected = gray
                tabPlaylists.setTextColor(if (currentTab == Tab.PLAYLISTS) Color.BLACK else grayColor)
            }
        }

        // Set initial selection
        updateTabSelection()
    }

    private fun setupSettingsButton() {
        btnSettings.setOnClickListener {
            showSettingsDialog()
        }

        btnSearch.setOnClickListener {
            // TODO: Implement search functionality
            Toast.makeText(requireContext(), "Search not implemented yet", Toast.LENGTH_SHORT).show()
        }

        // Intercept DOWN key for search button to focus first card
        btnSearch.setOnKeyListener { _, keyCode, event ->
            if (keyCode == KeyEvent.KEYCODE_DPAD_DOWN && event.action == KeyEvent.ACTION_DOWN) {
                recyclerView.post {
                    recyclerView.scrollToPosition(0)
                    recyclerView.postDelayed({
                        val firstViewHolder = recyclerView.findViewHolderForAdapterPosition(0)
                        firstViewHolder?.itemView?.requestFocus()
                    }, 50)
                }
                true
            } else {
                false
            }
        }

        // Hide preview when settings or search receive focus
        btnSettings.setOnFocusChangeListener { _, hasFocus ->
            if (hasFocus) {
                hideDetailCard()
            }
        }

        btnSearch.setOnFocusChangeListener { _, hasFocus ->
            if (hasFocus) {
                hideDetailCard()
            }
        }
    }

    private fun switchTab(tab: Tab) {
        if (currentTab == tab) return

        currentTab = tab
        updateTabSelection()

        when (tab) {
            Tab.MEDIA -> loadMediaContent()
            Tab.PLAYLISTS -> loadPlaylistsContent()
        }
    }

    private fun updateTabSelection() {
        tabMedia.isSelected = (currentTab == Tab.MEDIA)
        tabPlaylists.isSelected = (currentTab == Tab.PLAYLISTS)

        // Update text color: selected = black on white, not selected = gray on transparent
        val grayColor = Color.parseColor("#999999")
        tabMedia.setTextColor(if (currentTab == Tab.MEDIA) Color.BLACK else grayColor)
        tabPlaylists.setTextColor(if (currentTab == Tab.PLAYLISTS) Color.BLACK else grayColor)

        // Ensure alpha is set to fully opaque
        tabMedia.alpha = 1f
        tabPlaylists.alpha = 1f
    }

    private fun loadMediaContent() {
        lifecycleScope.launch {
            try {
                val result = mediaRepository.getMediaList(onlyUnassigned = true)
                result.onSuccess { mediaList ->
                    adapter.submitList(mediaList)
                }.onFailure { error ->
                    showError("${getString(R.string.media_error)}: ${error.message}")
                }
            } catch (e: Exception) {
                showError("${getString(R.string.error)}: ${e.message}")
            }
        }
    }

    private fun loadPlaylistsContent() {
        lifecycleScope.launch {
            try {
                val result = playlistRepository.getPlaylists()
                result.onSuccess { playlists ->
                    adapter.submitList(playlists)
                }.onFailure { error ->
                    showError("${getString(R.string.playlists_error)}: ${error.message}")
                }
            } catch (e: Exception) {
                showError("${getString(R.string.error)}: ${e.message}")
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

                Log.d("MainBrowseFragment", "Updated detail title: ${item.name}")

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
                detailInfo.text = "Type: ${item.getTypeDisplayName()}"
                detailInfo.visibility = View.VISIBLE

                Log.d("MainBrowseFragment", "Updated detail title: ${item.name}")

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

    private fun openMediaDetails(media: Media) {
        val intent = Intent(requireContext(), MediaDetailsActivity::class.java).apply {
            putExtra(MediaDetailsActivity.EXTRA_MEDIA_ID, media.id)
            putExtra(MediaDetailsActivity.EXTRA_MEDIA_NAME, media.name)
            putExtra(MediaDetailsActivity.EXTRA_MEDIA_POSTER, media.poster)
            putExtra(MediaDetailsActivity.EXTRA_MEDIA_DURATION, media.duration ?: 0)
            putExtra(MediaDetailsActivity.EXTRA_MEDIA_SIZE, media.size ?: 0L)
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

    private fun showSettingsDialog() {
        val dialog = ServerSettingsDialogFragment()
        dialog.show(childFragmentManager, "ServerSettingsDialog")
    }

    private fun showError(message: String) {
        activity?.runOnUiThread {
            Toast.makeText(requireContext(), message, Toast.LENGTH_LONG).show()
        }
    }
}
