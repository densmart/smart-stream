package com.smartstream.tvclient.data.api

import com.smartstream.tvclient.data.model.*
import retrofit2.Response
import retrofit2.http.*

/**
 * ClientAPI service interface
 * According to API.md specification
 * Base URL: http://localhost:17778 (default)
 */
interface ClientApiService {

    /**
     * Authentication
     * POST /auth/sign-in/
     */
    @POST("auth/sign-in/")
    suspend fun signIn(
        @Body request: ClientSignInRequest
    ): Response<ApiResponse<TokenResponse>>

    /**
     * Get playlists
     * GET /playlists/
     */
    @GET("playlists/")
    suspend fun getPlaylists(
        @Query("limit") limit: Int? = null,
        @Query("offset") offset: Int? = null,
        @Query("id") id: String? = null,
        @Query("name") name: String? = null,
        @Query("type") type: String? = null,
        @Query("order") order: String? = null
    ): Response<ApiResponse<List<Playlist>>>

    /**
     * Get playlist children (sub-playlists)
     * GET /playlists/:id/children/
     */
    @GET("playlists/{id}/children/")
    suspend fun getPlaylistChildren(
        @Path("id") playlistId: String,
        @Query("limit") limit: Int? = null,
        @Query("offset") offset: Int? = null,
        @Query("id") id: String? = null,
        @Query("name") name: String? = null,
        @Query("type") type: String? = null,
        @Query("order") order: String? = null
    ): Response<ApiResponse<List<Playlist>>>

    /**
     * Get media in specific playlist
     * GET /playlists/:id/media/
     */
    @GET("playlists/{id}/media/")
    suspend fun getPlaylistMedia(
        @Path("id") playlistId: String,
        @Query("limit") limit: Int? = null,
        @Query("offset") offset: Int? = null
    ): Response<ApiResponse<List<Media>>>

    /**
     * Get unassigned media (not in any playlist)
     * GET /media/
     */
    @GET("media/")
    suspend fun getUnassignedMedia(
        @Query("limit") limit: Int? = null,
        @Query("offset") offset: Int? = null,
        @Query("id") id: String? = null,
        @Query("name") name: String? = null,
        @Query("format") format: String? = null,
        @Query("order") order: String? = null
    ): Response<ApiResponse<List<Media>>>

    /**
     * Stream media
     * GET /media/:id/stream/
     * Note: This endpoint is used by ExoPlayer directly via URL
     * The stream URL format: {BASE_URL}media/{id}/stream/
     */
}
