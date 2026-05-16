package com.smartstream.tvclient.ui.settings

import android.app.Dialog
import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Button
import android.widget.EditText
import android.widget.Toast
import androidx.fragment.app.DialogFragment
import com.smartstream.tvclient.R
import com.smartstream.tvclient.utils.SharedPrefsManager

/**
 * Dialog for configuring server settings
 */
class ServerSettingsDialogFragment : DialogFragment() {

    private lateinit var editHost: EditText
    private lateinit var editPort: EditText
    private lateinit var btnSave: Button
    private lateinit var btnCancel: Button

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        return inflater.inflate(R.layout.dialog_server_settings, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        initViews(view)
        loadCurrentSettings()
        setupButtons()
    }

    override fun onStart() {
        super.onStart()
        // Make dialog wider
        dialog?.window?.setLayout(
            ViewGroup.LayoutParams.WRAP_CONTENT,
            ViewGroup.LayoutParams.WRAP_CONTENT
        )
    }

    private fun initViews(view: View) {
        editHost = view.findViewById(R.id.edit_host)
        editPort = view.findViewById(R.id.edit_port)
        btnSave = view.findViewById(R.id.btn_save)
        btnCancel = view.findViewById(R.id.btn_cancel)
    }

    private fun loadCurrentSettings() {
        val host = SharedPrefsManager.getServerHost()
        val port = SharedPrefsManager.getServerPort()

        editHost.setText(host)
        editPort.setText(port)
    }

    private fun setupButtons() {
        btnSave.setOnClickListener {
            saveSettings()
        }

        btnCancel.setOnClickListener {
            dismiss()
        }
    }

    private fun saveSettings() {
        val host = editHost.text.toString().trim()
        val portStr = editPort.text.toString().trim()

        // Validate
        if (host.isEmpty()) {
            showError(getString(R.string.settings_invalid_host))
            return
        }

        val port = portStr.toIntOrNull()
        if (port == null || port < 1 || port > 65535) {
            showError(getString(R.string.settings_invalid_port))
            return
        }

        // Save
        SharedPrefsManager.saveServerHost(host)
        SharedPrefsManager.saveServerPort(portStr)

        val baseUrl = "http://$host:$portStr/"
        Toast.makeText(
            requireContext(),
            getString(R.string.settings_saved, baseUrl),
            Toast.LENGTH_LONG
        ).show()

        dismiss()
    }

    private fun showError(message: String) {
        Toast.makeText(requireContext(), message, Toast.LENGTH_SHORT).show()
    }
}
