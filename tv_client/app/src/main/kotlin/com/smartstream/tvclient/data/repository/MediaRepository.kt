package com.smartstream.tvclient.data.repository

import com.smartstream.tvclient.data.api.RetrofitClient
import com.smartstream.tvclient.data.model.ApiResponse
import com.smartstream.tvclient.data.model.Media
import com.smartstream.tvclient.utils.Constants

/**
 * Repository for media-related operations
 */
class MediaRepository {

    private val apiService = RetrofitClient.getApiService()

    /**
     * Get unassigned media (not in any playlist)
     */
    suspend fun getUnassignedMedia(
        limit: Int = Constants.DEFAULT_PAGE_LIMIT,
        offset: Int = Constants.DEFAULT_PAGE_OFFSET
    ): Result<ApiResponse<List<Media>>> {
        return try {
            val response = apiService.getUnassignedMedia(
                limit = limit,
                offset = offset
            )
            if (response.isSuccessful && response.body() != null) {
                Result.success(response.body()!!)
            } else {
                Result.failure(Exception("Failed to load media: ${response.code()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Get media from specific playlist
     */
    suspend fun getPlaylistMedia(
        playlistId: String,
        limit: Int = Constants.DEFAULT_PAGE_LIMIT,
        offset: Int = Constants.DEFAULT_PAGE_OFFSET
    ): Result<ApiResponse<List<Media>>> {
        return try {
            val response = apiService.getPlaylistMedia(
                playlistId = playlistId,
                limit = limit,
                offset = offset
            )
            if (response.isSuccessful && response.body() != null) {
                Result.success(response.body()!!)
            } else {
                Result.failure(Exception("Failed to load playlist media: ${response.code()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Get media list with ApiResponse (includes pagination)
     */
    suspend fun getMediaListWithPagination(
        onlyUnassigned: Boolean = false,
        limit: Int = Constants.DEFAULT_PAGE_LIMIT,
        offset: Int = Constants.DEFAULT_PAGE_OFFSET
    ): Result<ApiResponse<List<Media>>> {
        return if (onlyUnassigned) {
            getUnassignedMedia(limit, offset)
        } else {
            getUnassignedMedia(limit, offset) // For now, same as unassigned
        }
    }

    /**
     * Get media list with optional filter for unassigned media (backward compatibility)
     */
    suspend fun getMediaList(
        onlyUnassigned: Boolean = false,
        limit: Int = Constants.DEFAULT_PAGE_LIMIT,
        offset: Int = Constants.DEFAULT_PAGE_OFFSET
    ): Result<List<Media>> {
        return try {
            val response = getMediaListWithPagination(onlyUnassigned, limit, offset)

            response.fold(
                onSuccess = { apiResponse ->
                    if (apiResponse.isSuccess() && apiResponse.result != null) {
                        Result.success(apiResponse.result)
                    } else {
                        Result.failure(Exception(apiResponse.error ?: "Unknown error"))
                    }
                },
                onFailure = { error ->
                    Result.failure(error)
                }
            )
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Search media by name
     */
    suspend fun searchMedia(
        query: String,
        limit: Int = Constants.DEFAULT_PAGE_LIMIT,
        offset: Int = Constants.DEFAULT_PAGE_OFFSET
    ): Result<ApiResponse<List<Media>>> {
        return try {
            val response = apiService.getUnassignedMedia(
                name = query,
                limit = limit,
                offset = offset
            )
            if (response.isSuccessful && response.body() != null) {
                Result.success(response.body()!!)
            } else {
                Result.failure(Exception("Failed to search media: ${response.code()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
}
