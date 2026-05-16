package com.smartstream.tvclient.data.repository

import com.smartstream.tvclient.data.api.RetrofitClient
import com.smartstream.tvclient.data.model.ApiResponse
import com.smartstream.tvclient.data.model.ClientSignInRequest
import com.smartstream.tvclient.data.model.TokenResponse
import com.smartstream.tvclient.utils.SharedPrefsManager

/**
 * Repository for authentication operations
 */
class AuthRepository {

    private val apiService = RetrofitClient.getApiService()

    /**
     * Sign in and get JWT token
     */
    suspend fun signIn(login: String, password: String): Result<String> {
        return try {
            val request = ClientSignInRequest(login, password)
            val response = apiService.signIn(request)

            if (response.isSuccessful && response.body() != null) {
                val apiResponse = response.body()!!
                if (apiResponse.isSuccess() && apiResponse.result != null) {
                    val token = apiResponse.result.accessToken
                    // Save token to SharedPreferences
                    SharedPrefsManager.saveToken(token)
                    Result.success(token)
                } else {
                    Result.failure(Exception(apiResponse.error))
                }
            } else {
                Result.failure(Exception("Authentication failed: ${response.code()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Logout - clear auth data
     */
    fun logout() {
        SharedPrefsManager.clearAuthData()
    }
}
