package com.smartstream.tvclient.ui.auth

import android.content.Intent
import android.os.Bundle
import android.view.View
import android.widget.Button
import android.widget.EditText
import android.widget.ImageButton
import android.widget.TextView
import androidx.fragment.app.FragmentActivity
import androidx.lifecycle.lifecycleScope
import com.smartstream.tvclient.R
import com.smartstream.tvclient.data.repository.AuthRepository
import com.smartstream.tvclient.ui.main.MainActivity
import kotlinx.coroutines.launch

/**
 * Authentication activity for TV.
 * Allows user to sign in with login and password.
 */
class AuthActivity : FragmentActivity() {

    private lateinit var etLogin: EditText
    private lateinit var etPassword: EditText
    private lateinit var btnSignIn: Button
    private lateinit var btnSettings: ImageButton
    private lateinit var tvError: TextView

    private val authRepository = AuthRepository()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_auth)

        initViews()
        setupListeners()
    }

    private fun initViews() {
        etLogin = findViewById(R.id.et_login)
        etPassword = findViewById(R.id.et_password)
        btnSignIn = findViewById(R.id.btn_sign_in)
        btnSettings = findViewById(R.id.btn_settings)
        tvError = findViewById(R.id.tv_error)

        // Request focus on login field
        etLogin.requestFocus()
    }

    private fun setupListeners() {
        btnSignIn.setOnClickListener {
            performSignIn()
        }

        btnSettings.setOnClickListener {
            openServerSettings()
        }
    }

    private fun openServerSettings() {
        ServerConfigDialog.show(this)
    }

    private fun performSignIn() {
        val login = etLogin.text.toString().trim()
        val password = etPassword.text.toString().trim()

        // Validate input
        if (login.isEmpty() || password.isEmpty()) {
            showError(getString(R.string.auth_error_empty_fields))
            return
        }

        // Hide error and disable button during sign-in
        hideError()
        setLoadingState(true)

        // Perform sign-in
        lifecycleScope.launch {
            try {
                val result = authRepository.signIn(login, password)

                result.onSuccess {
                    // Sign-in successful, navigate to main screen
                    navigateToMainScreen()
                }.onFailure { error ->
                    // Sign-in failed, show error
                    setLoadingState(false)
                    val errorMessage = when {
                        error.message?.contains("Network", ignoreCase = true) == true ->
                            getString(R.string.auth_error_network)
                        else ->
                            error.message ?: getString(R.string.auth_error)
                    }
                    showError(errorMessage)
                }
            } catch (e: Exception) {
                setLoadingState(false)
                showError(getString(R.string.auth_error_network))
            }
        }
    }

    private fun setLoadingState(loading: Boolean) {
        btnSignIn.isEnabled = !loading
        etLogin.isEnabled = !loading
        etPassword.isEnabled = !loading

        btnSignIn.text = if (loading) {
            getString(R.string.auth_signing_in)
        } else {
            getString(R.string.auth_sign_in_button)
        }
    }

    private fun showError(message: String) {
        tvError.text = message
        tvError.visibility = View.VISIBLE
    }

    private fun hideError() {
        tvError.visibility = View.GONE
    }

    private fun navigateToMainScreen() {
        val intent = Intent(this, MainActivity::class.java)
        intent.flags = Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TASK
        startActivity(intent)
        finish()
    }
}
