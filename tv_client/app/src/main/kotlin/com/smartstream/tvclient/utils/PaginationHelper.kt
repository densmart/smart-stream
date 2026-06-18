package com.smartstream.tvclient.utils

import android.view.View
import android.widget.TextView
import com.smartstream.tvclient.data.model.Pagination

/**
 * Helper class for managing pagination UI and logic
 */
class PaginationHelper(
    private val paginationContainer: View,
    private val btnPrevious: TextView,
    private val btnNext: TextView,
    private val paginationInfo: TextView,
    private val limit: Int = Constants.DEFAULT_PAGE_LIMIT
) {
    private var currentOffset = 0
    private var totalPages = 0
    private var totalItems = 0L

    /**
     * Update pagination UI with new data
     */
    fun update(pagination: Pagination?, offset: Int? = null) {
        // Sync internal offset with provided offset if given
        if (offset != null) {
            currentOffset = offset
        }

        if (pagination == null || pagination.total <= limit) {
            // Hide pagination if no data or all items fit on one page
            paginationContainer.visibility = View.GONE
            return
        }

        paginationContainer.visibility = View.VISIBLE
        totalPages = pagination.pages
        totalItems = pagination.total

        val currentPage = getCurrentPage()
        paginationInfo.text = "Страница $currentPage из $totalPages"

        // Update button states
        updateButtonStates()
    }

    /**
     * Get current page number (1-based)
     */
    fun getCurrentPage(): Int {
        return (currentOffset / limit) + 1
    }

    /**
     * Get current offset
     */
    fun getOffset(): Int {
        return currentOffset
    }

    /**
     * Check if can go to previous page
     */
    fun canGoPrevious(): Boolean {
        return currentOffset >= limit
    }

    /**
     * Check if can go to next page
     */
    fun canGoNext(): Boolean {
        return currentOffset + limit < totalItems
    }

    /**
     * Go to previous page
     * @return new offset or null if can't go previous
     */
    fun goPrevious(): Int? {
        return if (canGoPrevious()) {
            currentOffset -= limit
            updateButtonStates()
            currentOffset
        } else {
            null
        }
    }

    /**
     * Go to next page
     * @return new offset or null if can't go next
     */
    fun goNext(): Int? {
        return if (canGoNext()) {
            currentOffset += limit
            updateButtonStates()
            currentOffset
        } else {
            null
        }
    }

    /**
     * Reset pagination to first page
     */
    fun reset() {
        currentOffset = 0
        totalPages = 0
        totalItems = 0L
        paginationContainer.visibility = View.GONE
    }

    /**
     * Update button enabled/disabled states
     */
    private fun updateButtonStates() {
        btnPrevious.isEnabled = canGoPrevious()
        btnPrevious.alpha = if (canGoPrevious()) 1.0f else 0.5f

        btnNext.isEnabled = canGoNext()
        btnNext.alpha = if (canGoNext()) 1.0f else 0.5f
    }

    /**
     * Set up click listeners
     */
    fun setupListeners(onPreviousClick: () -> Unit, onNextClick: () -> Unit) {
        btnPrevious.setOnClickListener {
            if (canGoPrevious()) {
                onPreviousClick()
            }
        }

        btnNext.setOnClickListener {
            if (canGoNext()) {
                onNextClick()
            }
        }
    }
}
