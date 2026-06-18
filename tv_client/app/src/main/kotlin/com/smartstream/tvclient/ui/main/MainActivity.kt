package com.smartstream.tvclient.ui.main

import android.content.Intent
import android.os.Bundle
import androidx.activity.OnBackPressedCallback
import androidx.fragment.app.FragmentActivity
import com.smartstream.tvclient.R
import com.smartstream.tvclient.ui.auth.AuthActivity
import com.smartstream.tvclient.utils.SharedPrefsManager

/**
 * Main launcher activity for TV.
 * Checks authentication status and redirects accordingly.
 */
class MainActivity : FragmentActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        setupBackPressHandler()
        checkAuthenticationStatus()
    }

    private fun setupBackPressHandler() {
        onBackPressedDispatcher.addCallback(this, object : OnBackPressedCallback(true) {
            override fun handleOnBackPressed() {
                // Check if search is active in the fragment
                val fragment = supportFragmentManager.findFragmentById(R.id.main_browse_fragment) as? MainBrowseFragmentNew
                if (fragment != null && fragment.isSearchActive) {
                    // Close search instead of exiting activity
                    fragment.closeSearch()
                } else {
                    // Default back behavior - exit app
                    finish()
                }
            }
        })
    }

    override fun onResume() {
        super.onResume()
        // Check auth status on resume in case token was cleared
        // (e.g., by UnauthorizedInterceptor on 401 response)
        if (!SharedPrefsManager.isLoggedIn()) {
            redirectToAuth()
        }
    }

    private fun checkAuthenticationStatus() {
        if (!SharedPrefsManager.isLoggedIn()) {
            redirectToAuth()
        } else {
            // User is logged in, show main browse fragment
            loadMainBrowseFragment()
        }
    }

    private fun redirectToAuth() {
        // User not logged in, go to auth screen
        startActivity(Intent(this, AuthActivity::class.java))
        finish()
    }

    private fun loadMainBrowseFragment() {
        val fragment = MainBrowseFragmentNew()
        supportFragmentManager.beginTransaction()
            .replace(R.id.main_browse_fragment, fragment)
            .commit()
    }
}
