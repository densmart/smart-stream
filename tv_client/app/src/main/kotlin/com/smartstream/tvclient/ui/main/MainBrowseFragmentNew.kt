package com.smartstream.tvclient.ui.main

import android.animation.ValueAnimator
import android.content.Context
import android.content.Intent
import android.graphics.Color
import android.graphics.Rect
import android.os.Bundle
import android.text.Editable
import android.text.TextWatcher
import android.view.KeyEvent
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.view.inputmethod.InputMethodManager
import android.widget.EditText
import android.widget.ImageButton
import android.widget.TextView
import android.widget.Toast
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import androidx.recyclerview.widget.GridLayoutManager
import androidx.recyclerview.widget.RecyclerView
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
    private lateinit var tabSeries: TextView
    private lateinit var tabFranchise: TextView
    private lateinit var tabCartoon: TextView
    private lateinit var btnSettings: ImageButton
    private lateinit var btnSearch: ImageButton

    private lateinit var paginationHelper: PaginationHelper

    // Search variables
    private lateinit var searchInput: EditText
    var isSearchActive = false // Public for MainActivity to check
        private set
    private var currentSearchQuery = ""
    private var searchJob: Job? = null
    private val searchWidthExpanded = 200 // dp

    private var currentTab = Tab.MEDIA
    private var mediaOffset = 0
    private var seriesOffset = 0
    private var franchiseOffset = 0
    private var cartoonOffset = 0

    enum class Tab {
        MEDIA, SERIES, FRANCHISE, CARTOON
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
        tabSeries = view.findViewById(R.id.tab_series)
        tabFranchise = view.findViewById(R.id.tab_franchise)
        tabCartoon = view.findViewById(R.id.tab_cartoon)
        btnSettings = view.findViewById(R.id.btn_settings)
        btnSearch = view.findViewById(R.id.btn_search)
        searchInput = view.findViewById(R.id.search_input)

        // Force visibility and text properties
        tabMedia.visibility = View.VISIBLE
        tabSeries.visibility = View.VISIBLE
        tabFranchise.visibility = View.VISIBLE
        tabCartoon.visibility = View.VISIBLE
        tabMedia.alpha = 1f
        tabSeries.alpha = 1f
        tabFranchise.alpha = 1f
        tabCartoon.alpha = 1f
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
                        Tab.SERIES -> {
                            seriesOffset = newOffset
                            loadSeriesContent()
                        }
                        Tab.FRANCHISE -> {
                            franchiseOffset = newOffset
                            loadFranchiseContent()
                        }
                        Tab.CARTOON -> {
                            cartoonOffset = newOffset
                            loadCartoonContent()
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
                        Tab.SERIES -> {
                            seriesOffset = newOffset
                            loadSeriesContent()
                        }
                        Tab.FRANCHISE -> {
                            franchiseOffset = newOffset
                            loadFranchiseContent()
                        }
                        Tab.CARTOON -> {
                            cartoonOffset = newOffset
                            loadCartoonContent()
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
                // Detail card preview disabled - no action needed
            },
            onFocusLost = {
                // Detail card preview disabled - no action needed
            }
        )

        recyclerView.adapter = adapter

        // Use standard GridLayoutManager with native TV focus handling
        recyclerView.layoutManager = GridLayoutManager(requireContext(), 6)

        // Enable native focus search for proper TV navigation
        recyclerView.isFocusable = false
        recyclerView.isFocusableInTouchMode = false

        // Add spacing between cards
        val spacing = resources.getDimensionPixelSize(R.dimen.card_margin)
        val itemDecoration = GridSpacingItemDecoration(6, spacing, true)
        recyclerView.addItemDecoration(itemDecoration)

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

    private fun setupTabButtons() {
        tabMedia.setOnClickListener {
            switchTab(Tab.MEDIA)
        }

        tabSeries.setOnClickListener {
            switchTab(Tab.SERIES)
        }

        tabFranchise.setOnClickListener {
            switchTab(Tab.FRANCHISE)
        }

        tabCartoon.setOnClickListener {
            switchTab(Tab.CARTOON)
        }

        // Remove custom DOWN key handling - let Android handle navigation naturally

        // Update text color on focus
        val grayColor = Color.parseColor("#999999")

        tabMedia.setOnFocusChangeListener { _, hasFocus ->
            if (hasFocus) {
                tabMedia.setTextColor(Color.BLACK)
            } else {
                tabMedia.setTextColor(if (currentTab == Tab.MEDIA) Color.BLACK else grayColor)
            }
        }

        tabSeries.setOnFocusChangeListener { _, hasFocus ->
            if (hasFocus) {
                tabSeries.setTextColor(Color.BLACK)
            } else {
                tabSeries.setTextColor(if (currentTab == Tab.SERIES) Color.BLACK else grayColor)
            }
        }

        tabFranchise.setOnFocusChangeListener { _, hasFocus ->
            if (hasFocus) {
                tabFranchise.setTextColor(Color.BLACK)
            } else {
                tabFranchise.setTextColor(if (currentTab == Tab.FRANCHISE) Color.BLACK else grayColor)
            }
        }

        tabCartoon.setOnFocusChangeListener { _, hasFocus ->
            if (hasFocus) {
                tabCartoon.setTextColor(Color.BLACK)
            } else {
                tabCartoon.setTextColor(if (currentTab == Tab.CARTOON) Color.BLACK else grayColor)
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
                Tab.MEDIA -> loadMediaContent()
                Tab.SERIES -> loadSeriesContent()
                Tab.FRANCHISE -> loadFranchiseContent()
                Tab.CARTOON -> loadCartoonContent()
            }
        }
    }

    private fun updateTabSelection() {
        tabMedia.isSelected = (currentTab == Tab.MEDIA)
        tabSeries.isSelected = (currentTab == Tab.SERIES)
        tabFranchise.isSelected = (currentTab == Tab.FRANCHISE)
        tabCartoon.isSelected = (currentTab == Tab.CARTOON)

        // Update text color: selected = black on white, not selected = gray on transparent
        val grayColor = Color.parseColor("#999999")
        tabMedia.setTextColor(if (currentTab == Tab.MEDIA) Color.BLACK else grayColor)
        tabSeries.setTextColor(if (currentTab == Tab.SERIES) Color.BLACK else grayColor)
        tabFranchise.setTextColor(if (currentTab == Tab.FRANCHISE) Color.BLACK else grayColor)
        tabCartoon.setTextColor(if (currentTab == Tab.CARTOON) Color.BLACK else grayColor)

        // Ensure alpha is set to fully opaque
        tabMedia.alpha = 1f
        tabSeries.alpha = 1f
        tabFranchise.alpha = 1f
        tabCartoon.alpha = 1f
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

    private fun loadSeriesContent() {
        lifecycleScope.launch {
            try {
                val result = playlistRepository.getPlaylistsWithPagination(
                    offset = seriesOffset,
                    type = "series"
                )
                result.onSuccess { apiResponse ->
                    if (apiResponse.isSuccess() && apiResponse.result != null) {
                        adapter.submitList(apiResponse.result)
                        paginationHelper.update(apiResponse.pagination, seriesOffset)
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

    private fun loadFranchiseContent() {
        lifecycleScope.launch {
            try {
                val result = playlistRepository.getPlaylistsWithPagination(
                    offset = franchiseOffset,
                    type = "franchise"
                )
                result.onSuccess { apiResponse ->
                    if (apiResponse.isSuccess() && apiResponse.result != null) {
                        adapter.submitList(apiResponse.result)
                        paginationHelper.update(apiResponse.pagination, franchiseOffset)
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

    private fun loadCartoonContent() {
        lifecycleScope.launch {
            try {
                val result = playlistRepository.getPlaylistsWithPagination(
                    offset = cartoonOffset,
                    type = "cartoon"
                )
                result.onSuccess { apiResponse ->
                    if (apiResponse.isSuccess() && apiResponse.result != null) {
                        adapter.submitList(apiResponse.result)
                        paginationHelper.update(apiResponse.pagination, cartoonOffset)
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
            Tab.SERIES -> {
                seriesOffset = 0
                loadSeriesContent()
            }
            Tab.FRANCHISE -> {
                franchiseOffset = 0
                loadFranchiseContent()
            }
            Tab.CARTOON -> {
                cartoonOffset = 0
                loadCartoonContent()
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
                    Tab.SERIES -> {
                        seriesOffset = 0
                        loadSeriesContent()
                    }
                    Tab.FRANCHISE -> {
                        franchiseOffset = 0
                        loadFranchiseContent()
                    }
                    Tab.CARTOON -> {
                        cartoonOffset = 0
                        loadCartoonContent()
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
                Tab.SERIES -> searchSeries(query)
                Tab.FRANCHISE -> searchFranchise(query)
                Tab.CARTOON -> searchCartoon(query)
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

    private suspend fun searchSeries(query: String) {
        val result = playlistRepository.searchPlaylists(
            query = query,
            offset = seriesOffset,
            type = "series"
        )
        result.onSuccess { apiResponse ->
            if (apiResponse.isSuccess() && apiResponse.result != null) {
                adapter.submitList(apiResponse.result)
                paginationHelper.update(apiResponse.pagination, seriesOffset)

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

    private suspend fun searchFranchise(query: String) {
        val result = playlistRepository.searchPlaylists(
            query = query,
            offset = franchiseOffset,
            type = "franchise"
        )
        result.onSuccess { apiResponse ->
            if (apiResponse.isSuccess() && apiResponse.result != null) {
                adapter.submitList(apiResponse.result)
                paginationHelper.update(apiResponse.pagination, franchiseOffset)

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

    private suspend fun searchCartoon(query: String) {
        val result = playlistRepository.searchPlaylists(
            query = query,
            offset = cartoonOffset,
            type = "cartoon"
        )
        result.onSuccess { apiResponse ->
            if (apiResponse.isSuccess() && apiResponse.result != null) {
                adapter.submitList(apiResponse.result)
                paginationHelper.update(apiResponse.pagination, cartoonOffset)

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
