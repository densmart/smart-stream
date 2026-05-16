package com.smartstream.tvclient.utils

import android.content.Context
import android.content.SharedPreferences
import com.smartstream.tvclient.SmartStreamApp

object SharedPrefsManager {

    private val prefs: SharedPreferences by lazy {
        SmartStreamApp.instance.getSharedPreferences(
            Constants.PREFS_NAME,
            Context.MODE_PRIVATE
        )
    }

    /**
     * Save JWT token
     */
    fun saveToken(token: String) {
        prefs.edit().apply {
            putString(Constants.KEY_JWT_TOKEN, token)
            putBoolean(Constants.KEY_IS_LOGGED_IN, true)
            apply()
        }
    }

    /**
     * Get JWT token
     */
    fun getToken(): String? {
        return prefs.getString(Constants.KEY_JWT_TOKEN, null)
    }

    /**
     * Check if user is logged in
     */
    fun isLoggedIn(): Boolean {
        return prefs.getBoolean(Constants.KEY_IS_LOGGED_IN, false) && getToken() != null
    }

    /**
     * Clear all auth data (logout)
     */
    fun clearAuthData() {
        prefs.edit().apply {
            remove(Constants.KEY_JWT_TOKEN)
            remove(Constants.KEY_CLIENT_ID)
            putBoolean(Constants.KEY_IS_LOGGED_IN, false)
            apply()
        }
    }

    /**
     * Save client ID (optional)
     */
    fun saveClientId(clientId: String) {
        prefs.edit().putString(Constants.KEY_CLIENT_ID, clientId).apply()
    }

    /**
     * Get client ID
     */
    fun getClientId(): String? {
        return prefs.getString(Constants.KEY_CLIENT_ID, null)
    }

    // ==================== Server Configuration ====================

    /**
     * Save server host (IP or domain)
     */
    fun saveServerHost(host: String) {
        prefs.edit().putString(Constants.KEY_SERVER_HOST, host).apply()
    }

    /**
     * Get server host
     * @return Saved host or default value
     */
    fun getServerHost(): String {
        return prefs.getString(Constants.KEY_SERVER_HOST, Constants.DEFAULT_SERVER_HOST)
            ?: Constants.DEFAULT_SERVER_HOST
    }

    /**
     * Save server port
     */
    fun saveServerPort(port: String) {
        prefs.edit().putString(Constants.KEY_SERVER_PORT, port).apply()
    }

    /**
     * Get server port
     * @return Saved port or default value
     */
    fun getServerPort(): String {
        return prefs.getString(Constants.KEY_SERVER_PORT, Constants.DEFAULT_SERVER_PORT)
            ?: Constants.DEFAULT_SERVER_PORT
    }

    /**
     * Get complete base URL (host + port)
     * @return Full URL like "http://192.168.0.17:18888/"
     */
    fun getBaseUrl(): String {
        val host = getServerHost()
        val port = getServerPort()
        return "http://$host:$port/"
    }

    /**
     * Reset server configuration to default values
     */
    fun resetServerConfig() {
        prefs.edit().apply {
            remove(Constants.KEY_SERVER_HOST)
            remove(Constants.KEY_SERVER_PORT)
            apply()
        }
    }

    // ==================== Media Playback Position ====================

    /**
     * Save playback position for a media file
     * @param mediaId Media ID
     * @param positionMs Position in milliseconds
     */
    fun saveMediaPosition(mediaId: String, positionMs: Long) {
        prefs.edit().putLong("${Constants.KEY_MEDIA_POSITION_PREFIX}$mediaId", positionMs).apply()
    }

    /**
     * Get saved playback position for a media file
     * @param mediaId Media ID
     * @return Position in milliseconds, or 0 if not found
     */
    fun getMediaPosition(mediaId: String): Long {
        return prefs.getLong("${Constants.KEY_MEDIA_POSITION_PREFIX}$mediaId", 0L)
    }

    /**
     * Clear saved playback position for a media file
     * @param mediaId Media ID
     */
    fun clearMediaPosition(mediaId: String) {
        prefs.edit().remove("${Constants.KEY_MEDIA_POSITION_PREFIX}$mediaId").apply()
    }
}
