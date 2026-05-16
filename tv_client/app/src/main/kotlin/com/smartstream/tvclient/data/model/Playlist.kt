package com.smartstream.tvclient.data.model

import com.google.gson.annotations.SerializedName

/**
 * Playlist model
 * According to API.md specification
 */
data class Playlist(
    @SerializedName("id")
    val id: String,

    @SerializedName("created_at")
    val createdAt: String,

    @SerializedName("name")
    val name: String,

    @SerializedName("type")
    val type: String,

    @SerializedName("poster")
    val poster: String?,

    @SerializedName("parent_id")
    val parentId: String? = null,

    @SerializedName("has_children")
    val hasChildren: Boolean = false
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
     * Get playlist type display name
     */
    fun getTypeDisplayName(): String {
        return when (type) {
            "franchise" -> "Franchise"
            "series" -> "Series"
            "favourites" -> "Favourites"
            else -> type.replaceFirstChar { it.uppercase() }
        }
    }
}
