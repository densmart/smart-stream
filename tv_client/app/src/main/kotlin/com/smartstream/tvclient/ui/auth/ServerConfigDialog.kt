package com.smartstream.tvclient.ui.auth

import android.app.Dialog
import android.os.Bundle
import android.view.View
import android.view.Window
import android.widget.Button
import android.widget.EditText
import android.widget.TextView
import android.widget.Toast
import androidx.fragment.app.DialogFragment
import androidx.fragment.app.FragmentActivity
import com.smartstream.tvclient.R
import com.smartstream.tvclient.utils.Constants
import com.smartstream.tvclient.utils.SharedPrefsManager
import java.net.InetAddress

/**
 * Dialog for configuring server IP and port settings.
 * Validates input and saves to SharedPreferences.
 */
class ServerConfigDialog : DialogFragment() {

    private lateinit var etHost: EditText
    private lateinit var etPort: EditText
    private lateinit var tvError: TextView
    private lateinit var btnSave: Button
    private lateinit var btnCancel: Button
    private lateinit var btnReset: Button

    override fun onCreateDialog(savedInstanceState: Bundle?): Dialog {
        val dialog = Dialog(requireContext())
        dialog.requestWindowFeature(Window.FEATURE_NO_TITLE)
        dialog.setContentView(R.layout.dialog_server_config)
        dialog.window?.setBackgroundDrawableResource(android.R.color.transparent)

        initViews(dialog)
        loadCurrentSettings()
        setupListeners()

        return dialog
    }

    private fun initViews(dialog: Dialog) {
        etHost = dialog.findViewById(R.id.et_host)
        etPort = dialog.findViewById(R.id.et_port)
        tvError = dialog.findViewById(R.id.tv_error)
        btnSave = dialog.findViewById(R.id.btn_save)
        btnCancel = dialog.findViewById(R.id.btn_cancel)
        btnReset = dialog.findViewById(R.id.btn_reset)

        // Request focus on first field
        etHost.requestFocus()
    }

    private fun loadCurrentSettings() {
        etHost.setText(SharedPrefsManager.getServerHost())
        etPort.setText(SharedPrefsManager.getServerPort())
    }

    private fun setupListeners() {
        btnSave.setOnClickListener {
            saveSettings()
        }

        btnCancel.setOnClickListener {
            dismiss()
        }

        btnReset.setOnClickListener {
            resetToDefault()
        }
    }

    private fun saveSettings() {
        val host = etHost.text.toString().trim()
        val port = etPort.text.toString().trim()

        // Validate host
        if (host.isEmpty()) {
            showError(getString(R.string.settings_invalid_host))
            return
        }

        if (!isValidHost(host)) {
            showError(getString(R.string.settings_invalid_host))
            return
        }

        // Validate port
        if (port.isEmpty()) {
            showError(getString(R.string.settings_invalid_port))
            return
        }

        val portNumber = port.toIntOrNull()
        if (portNumber == null || portNumber < Constants.MIN_PORT || portNumber > Constants.MAX_PORT) {
            showError(getString(R.string.settings_invalid_port))
            return
        }

        // Save settings
        SharedPrefsManager.saveServerHost(host)
        SharedPrefsManager.saveServerPort(port)

        // Show success message
        val newUrl = SharedPrefsManager.getBaseUrl()
        Toast.makeText(
            requireContext(),
            getString(R.string.settings_saved, newUrl),
            Toast.LENGTH_LONG
        ).show()

        dismiss()
    }

    private fun resetToDefault() {
        SharedPrefsManager.resetServerConfig()
        loadCurrentSettings()
        hideError()

        Toast.makeText(
            requireContext(),
            getString(R.string.settings_saved, SharedPrefsManager.getBaseUrl()),
            Toast.LENGTH_SHORT
        ).show()
    }

    private fun isValidHost(host: String): Boolean {
        // Check if it's "localhost"
        if (host.equals("localhost", ignoreCase = true)) {
            return true
        }

        // Try to validate as IP address
        return try {
            val parts = host.split(".")
            if (parts.size != 4) {
                return false
            }

            parts.all { part ->
                val num = part.toIntOrNull()
                num != null && num in 0..255
            }
        } catch (e: Exception) {
            false
        }
    }

    private fun showError(message: String) {
        tvError.text = message
        tvError.visibility = View.VISIBLE
    }

    private fun hideError() {
        tvError.visibility = View.GONE
    }

    companion object {
        const val TAG = "ServerConfigDialog"

        /**
         * Show the server configuration dialog
         */
        fun show(activity: FragmentActivity) {
            val dialog = ServerConfigDialog()
            dialog.show(activity.supportFragmentManager, TAG)
        }
    }
}
