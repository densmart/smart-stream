package com.smartstream.tvclient.data.repository

import com.smartstream.tvclient.data.api.RetrofitClient
import com.smartstream.tvclient.data.model.ApiResponse
import com.smartstream.tvclient.data.model.Playlist
import com.smartstream.tvclient.utils.Constants

/**
 * Repository for playlist-related operations
 */
class PlaylistRepository {

    private val apiService = RetrofitClient.getApiService()

    /**
     * Get all playlists with ApiResponse wrapper
     */
    suspend fun getPlaylistsRaw(
        limit: Int = Constants.DEFAULT_PAGE_LIMIT,
        offset: Int = Constants.DEFAULT_PAGE_OFFSET
    ): Result<ApiResponse<List<Playlist>>> {
        return try {
            val response = apiService.getPlaylists(
                limit = limit,
                offset = offset
            )
            if (response.isSuccessful && response.body() != null) {
                Result.success(response.body()!!)
            } else {
                Result.failure(Exception("Failed to load playlists: ${response.code()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Get all playlists as simple list
     */
    suspend fun getPlaylists(
        limit: Int = Constants.DEFAULT_PAGE_LIMIT,
        offset: Int = Constants.DEFAULT_PAGE_OFFSET
    ): Result<List<Playlist>> {
        return try {
            val response = getPlaylistsRaw(limit, offset)

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
     * Get child playlists with ApiResponse wrapper
     */
    suspend fun getPlaylistChildrenRaw(
        playlistId: String,
        limit: Int = Constants.DEFAULT_PAGE_LIMIT,
        offset: Int = Constants.DEFAULT_PAGE_OFFSET
    ): Result<ApiResponse<List<Playlist>>> {
        return try {
            val response = apiService.getPlaylistChildren(
                playlistId = playlistId,
                limit = limit,
                offset = offset
            )
            if (response.isSuccessful && response.body() != null) {
                Result.success(response.body()!!)
            } else {
                Result.failure(Exception("Failed to load playlist children: ${response.code()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Get child playlists as simple list
     */
    suspend fun getPlaylistChildren(
        playlistId: String,
        limit: Int = Constants.DEFAULT_PAGE_LIMIT,
        offset: Int = Constants.DEFAULT_PAGE_OFFSET
    ): Result<List<Playlist>> {
        return try {
            val response = getPlaylistChildrenRaw(playlistId, limit, offset)

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
}
