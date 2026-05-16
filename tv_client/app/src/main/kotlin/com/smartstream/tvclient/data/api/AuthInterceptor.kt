package com.smartstream.tvclient.data.api

import com.smartstream.tvclient.utils.SharedPrefsManager
import okhttp3.Interceptor
import okhttp3.Response

/**
 * JWT authentication interceptor
 * Adds Authorization header with JWT token to all requests
 */
class AuthInterceptor : Interceptor {

    override fun intercept(chain: Interceptor.Chain): Response {
        val originalRequest = chain.request()

        // Get JWT token from SharedPreferences
        val token = SharedPrefsManager.getToken()

        // If token exists, add it to Authorization header
        val newRequest = if (token != null) {
            originalRequest.newBuilder()
                .header("Authorization", "Bearer $token")
                .build()
        } else {
            originalRequest
        }

        return chain.proceed(newRequest)
    }
}
