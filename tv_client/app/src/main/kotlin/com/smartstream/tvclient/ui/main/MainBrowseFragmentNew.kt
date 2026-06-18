package com.smartstream.tvclient.ui.main

import android.animation.ValueAnimator
import android.content.Context
import android.content.Intent
import android.graphics.Color
import android.graphics.Rect
import android.os.Bundle
import android.text.Editable
import android.text.TextWatcher
import android.util.Log
import android.view.KeyEvent
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.view.inputmethod.InputMethodManager
import android.widget.EditText
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
import com.smartstream.tvclient.utils.PaginationHelper
import com.smartstream.tvclient.utils.SharedPrefsManager
import kotlinx.coroutines.Job
import kotlinx.coroutines.delay
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

    private lateinit var paginationHelper: PaginationHelper

    // Search variables
    private lateinit var searchInput: EditText
    var isSearchActive = false // Public for MainActivity to check
        private set
    private var currentSearchQuery = ""
    private var searchJob: Job? = null
    private val searchWidthExpanded = 200 // dp

    // Grid decoration
    private var currentItemDecoration: RecyclerView.ItemDecoration? = null

    private var currentTab = Tab.MEDIA
    private var mediaOffset = 0
    private var playlistsOffset = 0

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
        setupPagination(view)
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
        searchInput = view.findViewById(R.id.search_input)

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
                    when (currentTab) {
                        Tab.MEDIA -> {
                            mediaOffset = newOffset
                            loadMediaContent()
                        }
                        Tab.PLAYLISTS -> {
                            playlistsOffset = newOffset
                            loadPlaylistsContent()
                        }
                    }
                }
            },
            onNextClick = {
                paginationHelper.goNext()?.let { newOffset ->
                    when (currentTab) {
                        Tab.MEDIA -> {
                            mediaOffset = newOffset
                            loadMediaContent()
                        }
                        Tab.PLAYLISTS -> {
                            playlistsOffset = newOffset
                            loadPlaylistsContent()
                        }
                    }
                }
            }
        )
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
            openSearch()
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

        // Load content for the selected tab (maintains current page/offset)
        // If search is active, apply same query to new tab
        if (isSearchActive && currentSearchQuery.length >= 3) {
            lifecycleScope.launch {
                executeSearch(currentSearchQuery)
            }
        } else {
            when (tab) {
                Tab.MEDIA -> {
                    loadMediaContent()
                }
                Tab.PLAYLISTS -> {
                    loadPlaylistsContent()
                }
            }
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
                val result = mediaRepository.getMediaListWithPagination(
                    onlyUnassigned = true,
                    offset = mediaOffset
                )
                result.onSuccess { apiResponse ->
                    if (apiResponse.isSuccess() && apiResponse.result != null) {
                        adapter.submitList(apiResponse.result)
                        paginationHelper.update(apiResponse.pagination, mediaOffset)

                        // Request focus on first item after data is loaded
                        recyclerView.postDelayed({
                            if (adapter.itemCount > 0 && recyclerView.childCount > 0) {
                                recyclerView.getChildAt(0)?.requestFocus()
                            }
                        }, 100)
                    } else {
                        showError(apiResponse.error ?: getString(R.string.media_error))
                    }
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
                val result = playlistRepository.getPlaylistsWithPagination(
                    offset = playlistsOffset
                )
                result.onSuccess { apiResponse ->
                    if (apiResponse.isSuccess() && apiResponse.result != null) {
                        adapter.submitList(apiResponse.result)
                        paginationHelper.update(apiResponse.pagination, playlistsOffset)

                        // Request focus on first item after data is loaded
                        recyclerView.postDelayed({
                            if (adapter.itemCount > 0 && recyclerView.childCount > 0) {
                                recyclerView.getChildAt(0)?.requestFocus()
                            }
                        }, 100)
                    } else {
                        showError(apiResponse.error ?: getString(R.string.playlists_error))
                    }
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
                detailInfo.text = item.getTypeDisplayName(requireContext())
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

    // Search functionality
    private fun openSearch() {
        if (isSearchActive) return

        isSearchActive = true

        // Animate search input expansion
        val widthInPx = (searchWidthExpanded * resources.displayMetrics.density).toInt()
        val animator = ValueAnimator.ofInt(0, widthInPx)
        animator.duration = 300
        animator.addUpdateListener { valueAnimator ->
            val layoutParams = searchInput.layoutParams
            layoutParams.width = valueAnimator.animatedValue as Int
            searchInput.layoutParams = layoutParams
        }

        searchInput.visibility = View.VISIBLE
        animator.start()

        // Focus search input and show keyboard
        searchInput.postDelayed({
            searchInput.requestFocus()
            val imm = requireContext().getSystemService(Context.INPUT_METHOD_SERVICE) as InputMethodManager
            imm.showSoftInput(searchInput, InputMethodManager.SHOW_IMPLICIT)
        }, 350)

        // Setup text watcher for search
        setupSearchWatcher()
    }

    fun closeSearch() {
        if (!isSearchActive) return

        isSearchActive = false
        currentSearchQuery = ""

        // Hide keyboard
        val imm = requireContext().getSystemService(Context.INPUT_METHOD_SERVICE) as InputMethodManager
        imm.hideSoftInputFromWindow(searchInput.windowToken, 0)

        // Clear search input
        searchInput.setText("")

        // Animate search input collapse
        val animator = ValueAnimator.ofInt(searchInput.width, 0)
        animator.duration = 300
        animator.addUpdateListener { valueAnimator ->
            val layoutParams = searchInput.layoutParams
            layoutParams.width = valueAnimator.animatedValue as Int
            searchInput.layoutParams = layoutParams
        }
        animator.start()

        searchInput.postDelayed({
            searchInput.visibility = View.GONE
            btnSearch.requestFocus()
        }, 350)

        // Cancel any pending search
        searchJob?.cancel()

        // Reload normal content
        when (currentTab) {
            Tab.MEDIA -> {
                mediaOffset = 0
                loadMediaContent()
            }
            Tab.PLAYLISTS -> {
                playlistsOffset = 0
                loadPlaylistsContent()
            }
        }
    }

    private fun setupSearchWatcher() {
        searchInput.addTextChangedListener(object : TextWatcher {
            override fun beforeTextChanged(s: CharSequence?, start: Int, count: Int, after: Int) {}

            override fun onTextChanged(s: CharSequence?, start: Int, before: Int, count: Int) {
                val query = s?.toString() ?: ""
                performSearch(query)
            }

            override fun afterTextChanged(s: Editable?) {}
        })
    }

    private fun performSearch(query: String) {
        // Cancel previous search job
        searchJob?.cancel()

        // If query is less than 3 characters, show normal list
        if (query.length < 3) {
            if (query.isEmpty() && currentSearchQuery.isNotEmpty()) {
                currentSearchQuery = ""
                // Reload normal content
                when (currentTab) {
                    Tab.MEDIA -> {
                        mediaOffset = 0
                        loadMediaContent()
                    }
                    Tab.PLAYLISTS -> {
                        playlistsOffset = 0
                        loadPlaylistsContent()
                    }
                }
            }
            return
        }

        currentSearchQuery = query

        // Debounce search with 500ms delay
        searchJob = lifecycleScope.launch {
            delay(500)
            executeSearch(query)
        }
    }

    private suspend fun executeSearch(query: String) {
        try {
            when (currentTab) {
                Tab.MEDIA -> searchMedia(query)
                Tab.PLAYLISTS -> searchPlaylists(query)
            }
        } catch (e: Exception) {
            showError("${getString(R.string.error)}: ${e.message}")
        }
    }

    private suspend fun searchMedia(query: String) {
        val result = mediaRepository.searchMedia(
            query = query,
            offset = mediaOffset
        )
        result.onSuccess { apiResponse ->
            if (apiResponse.isSuccess() && apiResponse.result != null) {
                adapter.submitList(apiResponse.result)
                paginationHelper.update(apiResponse.pagination, mediaOffset)

                // Show empty message if no results
                if (apiResponse.result.isEmpty()) {
                    showError(getString(R.string.search_results_empty))
                }
            } else {
                showError(apiResponse.error ?: getString(R.string.media_error))
            }
        }.onFailure { error ->
            showError("${getString(R.string.media_error)}: ${error.message}")
        }
    }

    private suspend fun searchPlaylists(query: String) {
        val result = playlistRepository.searchPlaylists(
            query = query,
            offset = playlistsOffset
        )
        result.onSuccess { apiResponse ->
            if (apiResponse.isSuccess() && apiResponse.result != null) {
                adapter.submitList(apiResponse.result)
                paginationHelper.update(apiResponse.pagination, playlistsOffset)

                // Show empty message if no results
                if (apiResponse.result.isEmpty()) {
                    showError(getString(R.string.search_results_empty))
                }
            } else {
                showError(apiResponse.error ?: getString(R.string.playlists_error))
            }
        }.onFailure { error ->
            showError("${getString(R.string.playlists_error)}: ${error.message}")
        }
    }

    private fun showError(message: String) {
        activity?.runOnUiThread {
            Toast.makeText(requireContext(), message, Toast.LENGTH_LONG).show()
        }
    }
}
