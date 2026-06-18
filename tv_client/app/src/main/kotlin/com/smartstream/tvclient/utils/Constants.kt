package com.smartstream.tvclient.utils

object Constants {
    // Server Configuration (Default values)
    const val DEFAULT_SERVER_HOST = "10.0.2.2"  // Android emulator host alias
    const val DEFAULT_SERVER_PORT = "18888"

    // SharedPreferences
    const val PREFS_NAME = "smart_stream_prefs"
    const val KEY_JWT_TOKEN = "jwt_token"
    const val KEY_CLIENT_ID = "client_id"
    const val KEY_IS_LOGGED_IN = "is_logged_in"
    const val KEY_SERVER_HOST = "server_host"
    const val KEY_SERVER_PORT = "server_port"
    const val KEY_MEDIA_POSITION_PREFIX = "media_position_"

    // Network
    const val CONNECT_TIMEOUT = 30L // seconds
    const val READ_TIMEOUT = 30L    // seconds
    const val WRITE_TIMEOUT = 30L   // seconds

    // Pagination
    const val DEFAULT_PAGE_LIMIT = 12
    const val DEFAULT_PAGE_OFFSET = 0

    // Validation
    const val MIN_PORT = 1
    const val MAX_PORT = 65535
}
