package com.smartstream.tvclient.data.model

import com.google.gson.annotations.SerializedName

/**
 * Media item model
 * According to API.md specification
 */
data class Media(
    @SerializedName("id")
    val id: String,

    @SerializedName("created_at")
    val createdAt: String,

    @SerializedName("playlist_id")
    val playlistId: String?,

    @SerializedName("name")
    val name: String,

    @SerializedName("poster")
    val poster: String?,

    @SerializedName("format")
    val format: String,

    @SerializedName("path")
    val path: String,

    @SerializedName("duration")
    val duration: Int?,

    @SerializedName("size")
    val size: Long?,

    @SerializedName("order")
    val order: Int
) {
    /**
     * Get full poster URL
     * Poster field can be:
     * - Full path: "/static/posters/uuid.jpg"
     * - Just filename: "uuid.jpg"
     * - Absolute URL: "http://..."
     */
    fun getPosterUrl(baseUrl: String): String? {
        return poster?.let {
            when {
                it.startsWith("http") -> it
                it.startsWith("/static/posters/") -> "$baseUrl${it.removePrefix("/")}"
                else -> "${baseUrl}static/posters/$it"
            }
        }
    }

    /**
     * Get stream URL for this media
     * GET /media/:id/stream/
     */
    fun getStreamUrl(baseUrl: String): String {
        return "${baseUrl}media/$id/stream/"
    }

    /**
     * Format duration as HH:MM:SS or MM:SS
     */
    fun getFormattedDuration(): String? {
        return duration?.let {
            val hours = it / 3600
            val minutes = (it % 3600) / 60
            val seconds = it % 60

            if (hours > 0) {
                String.format("%02d:%02d:%02d", hours, minutes, seconds)
            } else {
                String.format("%02d:%02d", minutes, seconds)
            }
        }
    }

    /**
     * Format file size as KB/MB/GB
     */
    fun getFormattedSize(): String? {
        return size?.let {
            when {
                it < 1024 -> "$it B"
                it < 1024 * 1024 -> String.format("%.2f KB", it / 1024.0)
                it < 1024 * 1024 * 1024 -> String.format("%.2f MB", it / (1024.0 * 1024.0))
                else -> String.format("%.2f GB", it / (1024.0 * 1024.0 * 1024.0))
            }
        }
    }
}
