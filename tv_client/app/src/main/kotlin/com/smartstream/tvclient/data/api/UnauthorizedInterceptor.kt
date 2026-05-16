package com.smartstream.tvclient.data.api

import android.content.Intent
import com.smartstream.tvclient.SmartStreamApp
import com.smartstream.tvclient.ui.auth.AuthActivity
import com.smartstream.tvclient.utils.SharedPrefsManager
import okhttp3.Interceptor
import okhttp3.Response

/**
 * Interceptor to handle 401 Unauthorized responses
 * Clears auth data and redirects to login screen when token expires
 */
class UnauthorizedInterceptor : Interceptor {

    override fun intercept(chain: Interceptor.Chain): Response {
        val request = chain.request()
        val response = chain.proceed(request)

        // Check if response is 401 Unauthorized
        if (response.code == 401) {
            // Token expired or invalid, clear auth data
            SharedPrefsManager.clearAuthData()

            // Redirect to auth screen
            redirectToAuthScreen()
        }

        return response
    }

    private fun redirectToAuthScreen() {
        val context = SmartStreamApp.instance
        val intent = Intent(context, AuthActivity::class.java).apply {
            flags = Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TASK
        }
        context.startActivity(intent)
    }
}
