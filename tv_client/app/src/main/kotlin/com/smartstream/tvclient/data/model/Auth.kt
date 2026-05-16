package com.smartstream.tvclient.data.model

import com.google.gson.annotations.SerializedName

/**
 * Client sign-in request
 * POST /auth/sign-in/
 */
data class ClientSignInRequest(
    @SerializedName("login")
    val login: String,

    @SerializedName("password")
    val password: String
)

/**
 * Authentication response with JWT token
 */
data class TokenResponse(
    @SerializedName("access_token")
    val accessToken: String
)
