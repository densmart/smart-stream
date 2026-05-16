package com.smartstream.tvclient.data.model

import com.google.gson.annotations.SerializedName

/**
 * Unified API response wrapper
 * According to API.md specification
 */
data class ApiResponse<T>(
    @SerializedName("error")
    val error: String,

    @SerializedName("pagination")
    val pagination: Pagination?,

    @SerializedName("result")
    val result: T?
) {
    /**
     * Check if request was successful (no error)
     */
    fun isSuccess(): Boolean = error.isEmpty()

    /**
     * Check if request failed
     */
    fun isError(): Boolean = error.isNotEmpty()
}

/**
 * Pagination information
 */
data class Pagination(
    @SerializedName("total")
    val total: Long,

    @SerializedName("pages")
    val pages: Int
)
