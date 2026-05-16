package com.smartstream.tvclient.data.api

import com.smartstream.tvclient.utils.Constants
import com.smartstream.tvclient.utils.SharedPrefsManager
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import java.util.concurrent.TimeUnit

/**
 * Retrofit client singleton
 * Provides configured instance of ClientApiService
 *
 * Note: BASE_URL is dynamically loaded from SharedPreferences
 */
object RetrofitClient {

    private var apiService: ClientApiService? = null
    private var currentBaseUrl: String? = null

    /**
     * Get ClientApiService instance
     * Automatically recreates if base URL has changed
     */
    fun getApiService(): ClientApiService {
        val baseUrl = SharedPrefsManager.getBaseUrl()

        // Recreate service if URL changed or service not initialized
        if (apiService == null || currentBaseUrl != baseUrl) {
            currentBaseUrl = baseUrl
            apiService = createApiService(baseUrl)
        }

        return apiService!!
    }

    /**
     * Reset API service (e.g., when base URL changes)
     * Call this after changing server settings
     */
    fun reset() {
        apiService = null
        currentBaseUrl = null
    }

    /**
     * Create configured ApiService instance with given base URL
     */
    private fun createApiService(baseUrl: String): ClientApiService {
        val okHttpClient = OkHttpClient.Builder()
            .addInterceptor(createLoggingInterceptor())
            .addInterceptor(AuthInterceptor())
            .addInterceptor(UnauthorizedInterceptor())
            .connectTimeout(Constants.CONNECT_TIMEOUT, TimeUnit.SECONDS)
            .readTimeout(Constants.READ_TIMEOUT, TimeUnit.SECONDS)
            .writeTimeout(Constants.WRITE_TIMEOUT, TimeUnit.SECONDS)
            .build()

        val retrofit = Retrofit.Builder()
            .baseUrl(baseUrl)
            .client(okHttpClient)
            .addConverterFactory(GsonConverterFactory.create())
            .build()

        return retrofit.create(ClientApiService::class.java)
    }

    /**
     * Create HTTP logging interceptor for debugging
     */
    private fun createLoggingInterceptor(): HttpLoggingInterceptor {
        val interceptor = HttpLoggingInterceptor()
        // Set BODY for development, NONE for production
        interceptor.level = HttpLoggingInterceptor.Level.BODY
        return interceptor
    }
}
